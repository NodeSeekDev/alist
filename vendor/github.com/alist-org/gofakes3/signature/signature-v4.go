package signature

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	// TimeNow is a variable that holds the function to get the current time
	TimeNow = time.Now
)

// ref: https://github.com/minio/minio/cmd/auth-handler.go

const (
	signV4Algorithm = "AWS4-HMAC-SHA256"
	iso8601Format   = "20060102T150405Z"
	yyyymmdd        = "20060102"
	serviceS3       = "s3"
	slashSeparator  = "/"
	stype           = serviceS3

	headerAuth       = "Authorization"
	headerDate       = "Date"
	amzContentSha256 = "X-Amz-Content-Sha256"
	amzDate          = "X-Amz-Date"
	amzAlgorithm     = "X-Amz-Algorithm"
	amzCredential    = "X-Amz-Credential"
	amzSignedHeaders = "X-Amz-SignedHeaders"
	amzSignature     = "X-Amz-Signature"
	amzExpires       = "X-Amz-Expires"
)

// getCanonicalHeaders generate a list of request headers with their values
func getCanonicalHeaders(signedHeaders http.Header) string {
	var headers []string
	vals := make(http.Header)
	for k, vv := range signedHeaders {
		headers = append(headers, strings.ToLower(k))
		vals[strings.ToLower(k)] = vv
	}
	sort.Strings(headers)

	var buf bytes.Buffer
	for _, k := range headers {
		buf.WriteString(k)
		buf.WriteByte(':')
		for idx, v := range vals[k] {
			if idx > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(signV4TrimAll(v))
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

// getSignedHeaders generate a string i.e alphabetically sorted, semicolon-separated list of lowercase request header names
func getSignedHeaders(signedHeaders http.Header) string {
	var headers []string
	for k := range signedHeaders {
		headers = append(headers, strings.ToLower(k))
	}
	sort.Strings(headers)
	return strings.Join(headers, ";")
}

// compareSignatureV4 returns true if and only if both signatures
// are equal. The signatures are expected to be HEX encoded strings
// according to the AWS S3 signature V4 spec.
func compareSignatureV4(sig1, sig2 string) bool {
	// The CTC using []byte(str) works because the hex encoding
	// is unique for a sequence of bytes. See also compareSignatureV2.
	return subtle.ConstantTimeCompare([]byte(sig1), []byte(sig2)) == 1
}

// getSignature final signature in hexadecimal form.
func getSignature(signingKey []byte, stringToSign string) string {
	return hex.EncodeToString(sumHMAC(signingKey, []byte(stringToSign)))
}

// Trim leading and trailing spaces and replace sequential spaces with one space, following Trimall()
// in http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
func signV4TrimAll(input string) string {
	// Compress adjacent spaces (a space is determined by
	// unicode.IsSpace() internally here) to one space and return
	return strings.Join(strings.Fields(input), " ")
}

// getCanonicalRequest generate a canonical request of style
//
// canonicalRequest =
//
//	<HTTPMethod>\n
//	<CanonicalURI>\n
//	<CanonicalQueryString>\n
//	<CanonicalHeaders>\n
//	<SignedHeaders>\n
//	<HashedPayload>
func getCanonicalRequest(extractedSignedHeaders http.Header, payload, queryStr, urlPath, method string) string {
	rawQuery := strings.ReplaceAll(queryStr, "+", "%20")
	encodedPath := encodePath(urlPath)
	canonicalRequest := strings.Join([]string{
		method,
		encodedPath,
		rawQuery,
		getCanonicalHeaders(extractedSignedHeaders),
		getSignedHeaders(extractedSignedHeaders),
		payload,
	}, "\n")
	return canonicalRequest
}

// getStringToSign a string based on selected query values.
func getStringToSign(canonicalRequest string, t time.Time, scope string) string {
	stringToSign := signV4Algorithm + "\n" + t.Format(iso8601Format) + "\n"
	stringToSign += scope + "\n"
	canonicalRequestBytes := sha256.Sum256([]byte(canonicalRequest))
	stringToSign += hex.EncodeToString(canonicalRequestBytes[:])
	return stringToSign
}

// getSigningKey hmac seed to calculate final signature.
func getSigningKey(secretKey string, t time.Time, region string) []byte {
	date := sumHMAC([]byte("AWS4"+secretKey), []byte(t.Format(yyyymmdd)))
	regionBytes := sumHMAC(date, []byte(region))
	service := sumHMAC(regionBytes, []byte(stype))
	signingKey := sumHMAC(service, []byte("aws4_request"))
	return signingKey
}

// V4SignVerify - Verify authorization header with calculated header in accordance with
//   - http://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html
//
// returns nil if signature matches.
func V4SignVerify(r *http.Request) ErrorCode {
	// Copy request.
	req := *r
	queryf := req.URL.Query()
	isUnsignedPayload := req.Header.Get("X-Amz-Content-Sha256") == "UNSIGNED-PAYLOAD"

	// Save authorization header.
	v4Auth := req.Header.Get(headerAuth)

	// If the header is empty but the query string has the signature, then it's QueryString authentication. (https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-query-string-auth.html)
	if v4Auth == "" && queryf.Get(amzSignature) != "" {
		// QueryString authentications are always "UNSIGNED-PAYLOAD".
		isUnsignedPayload = true
		v4Auth = fmt.Sprintf("%s Credential=%s, SignedHeaders=%s, Signature=%s", queryf.Get(amzAlgorithm), queryf.Get(amzCredential), queryf.Get(amzSignedHeaders), queryf.Get(amzSignature))
		if queryf.Get(amzCredential) == "" {
			return errMissingCredTag
		}
	}

	// Parse signature version '4' header.
	signV4Values, Err := ParseSignV4(v4Auth)
	if Err != ErrNone {
		return Err
	}

	cred, _, Err := checkKeyValid(r, signV4Values.Credential.accessKey)
	if Err != ErrNone {
		return Err
	}

	// Extract all the signed headers along with its values.
	extractedSignedHeaders, ErrCode := extractSignedHeaders(signV4Values.SignedHeaders, r)
	if ErrCode != ErrNone {
		return ErrCode
	}

	// Extract date from various possible sources
	date := req.Header.Get(amzDate)
	if date == "" {
		date = req.Header.Get(headerDate)
	}
	if date == "" {
		date = queryf.Get(amzDate)
	}

	// If date is still empty after checking all sources, return an error
	if date == "" {
		return errMissingDateHeader
	}

	// Parse date header.
	t, e := time.Parse(iso8601Format, date)
	if e != nil {
		return errMalformedDate
	}

	// Check expiration
	expiresStr := queryf.Get(amzExpires)
	var expires time.Duration
	if expiresStr == "" {
		// If expires is not set, use the default of 15 minutes
		expires = 15 * time.Minute
	} else {
		expiresInt, err := strconv.ParseInt(expiresStr, 10, 64)
		if err != nil {
			return errMalformedExpires
		}
		expires = time.Duration(expiresInt) * time.Second
	}
	if TimeNow().After(t.Add(expires)) {
		return errExpiredRequest
	}

	// Query string.
	queryf.Del(amzSignature)
	rawquery := queryf.Encode()

	// Get hmac signing key.
	signingKey := getSigningKey(cred.SecretKey, signV4Values.Credential.scope.date, signV4Values.Credential.scope.region)

	var newSignature string
	if isUnsignedPayload {
		hashedPayload := "UNSIGNED-PAYLOAD"
		canonicalRequest := getCanonicalRequest(extractedSignedHeaders, hashedPayload, rawquery, req.URL.Path, req.Method)
		stringToSign := getStringToSign(canonicalRequest, t, signV4Values.Credential.getScope())
		newSignature = getSignature(signingKey, stringToSign)
	} else {
		hashedPayload := getContentSha256Cksum(r)
		canonicalRequest := getCanonicalRequest(extractedSignedHeaders, hashedPayload, rawquery, req.URL.Path, req.Method)
		stringToSign := getStringToSign(canonicalRequest, t, signV4Values.Credential.getScope())
		newSignature = getSignature(signingKey, stringToSign)
	}

	// Verify if signature match.
	if !compareSignatureV4(newSignature, signV4Values.Signature) {
		return errSignatureDoesNotMatch
	}

	// Return Error none.
	return ErrNone
}
