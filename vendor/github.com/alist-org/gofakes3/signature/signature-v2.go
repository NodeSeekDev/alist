package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// ref: https://github.com/minio/minio/cmd/auth-handler.go

const (
	signV2Algorithm = "AWS"
)

// AWS S3 Signature V2 calculation rule is give here:
// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html#RESTAuthenticationStringToSign

// Whitelist resource list that will be used in query string for signature-V2 calculation.
var resourceList = []string{
	"acl",
	"lifecycle",
	"location",
	"logging",
	"notification",
	"partNumber",
	"policy",
	"requestPayment",
	"torrent",
	"uploadId",
	"uploads",
	"versionId",
	"versioning",
	"versions",
	"website",
}

// CredentialsV2 - for signature-V2 calculation
type CredentialsV2 struct {
	AccessKey string
	SecretKey string
}

type signV2Values struct {
	Credential CredentialsV2
	Signature  string
}

var log = logger.New()

// Sign - return the Authorization header value.
func (c CredentialsV2) SignV2(method string,
	encodedResource string,
	encodedQuery string,
	headers http.Header,
	expires string) string {
	canonicalHeaaders := canonicalizedAmzHeadersV2(headers)
	if len(canonicalHeaaders) > 0 {
		canonicalHeaaders += "\n"
	}

	date := headers.Get("Date")
	if expires != "" {
		date = expires
	}

	stringToSign := strings.Join([]string{
		method,
		headers.Get("Content-MD5"),
		headers.Get("Content-Type"),
		date,
		canonicalHeaaders,
	}, "\n") + canonicalizedResourceV2(encodedResource, encodedQuery)

	hm := hmac.New(sha1.New, []byte(c.SecretKey))
	log.Debugf("stringToSign: %s", stringToSign)
	hm.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(hm.Sum(nil))
	return signature
	//return fmt.Sprintf("%s %s:%s", SignV2Algorithm, c.AccessKey, signature)
}

// Return canonical headers.
func canonicalizedAmzHeadersV2(headers http.Header) string {
	var keys []string
	keyval := make(map[string]string)
	for key := range headers {
		lkey := strings.ToLower(key)
		if !strings.HasPrefix(lkey, "x-amz-") {
			continue
		}
		keys = append(keys, lkey)
		keyval[lkey] = strings.Join(headers[key], ",")
	}
	sort.Strings(keys)
	var canonicalHeaders []string
	for _, key := range keys {
		canonicalHeaders = append(canonicalHeaders, key+":"+keyval[key])
	}
	return strings.Join(canonicalHeaders, "\n")
}

// Return canonical resource string.
func canonicalizedResourceV2(encodedPath string, encodedQuery string) string {
	queries := strings.Split(encodedQuery, "&")
	keyval := make(map[string]string)
	for _, query := range queries {
		key := query
		val := ""
		index := strings.Index(query, "=")
		if index != -1 {
			key = query[:index]
			val = query[index+1:]
		}
		keyval[key] = val
	}
	var canonicalQueries []string
	for _, key := range resourceList {
		val, ok := keyval[key]
		if !ok {
			continue
		}
		if val == "" {
			canonicalQueries = append(canonicalQueries, key)
			continue
		}
		canonicalQueries = append(canonicalQueries, key+"="+val)
	}
	if len(canonicalQueries) == 0 {
		return encodedPath
	}
	// the queries will be already sorted as resourceList is sorted.
	return encodedPath + "?" + strings.Join(canonicalQueries, "&")
}

func V2SignVerify(r *http.Request) ErrorCode {
	req := *r

	// Save authorization header.
	// example: AWS NXM00iAJpLvPYC94VXXX:q20e3MP/dZOpnzBZTvo619OXXXX=
	v2Auth := req.Header.Get(headerAuth)
	expires := ""
	if v2Auth == "" {
		queryf := req.URL.Query()

		urlEncodedSignature := queryf.Get("Signature")
		signature, err := url.QueryUnescape(urlEncodedSignature)
		if err != nil {
			log.Warnf("Error decoding signature: %s", urlEncodedSignature)
			return errCredMalformed
		}

		expires = queryf.Get("Expires")

		v2Auth = fmt.Sprintf("%s %s:%s", signV2Algorithm, queryf.Get("AWSAccessKeyId"), signature)
		if v2Auth == "" {
			return errMissingCredTag
		}
	}

	// Get the AWS access key and secret key
	signV2Values, Err := ParseSignV2(v2Auth)
	if Err != ErrNone {
		return Err
	}

	credential := signV2Values.Credential
	accessKey := credential.AccessKey
	cred, _, Err := checkKeyValid(r, accessKey)
	if Err != ErrNone {
		return Err
	}

	secretKey := cred.SecretKey
	credential.SecretKey = secretKey
	receivedSignature := signV2Values.Signature
	log.Debugf("receivedSignature: %s", receivedSignature)

	encodedResource := req.URL.RawPath
	encodedQuery := req.URL.RawQuery
	if encodedResource == "" {
		splits := strings.Split(req.URL.Path, "?")
		if len(splits) > 0 {
			encodedResource = splits[0]
		}
	}

	// Create the string to sign
	expectedSignature := credential.SignV2(req.Method, encodedResource, encodedQuery, req.Header, expires)
	log.Debugf("expectedSignature: %s", expectedSignature)

	// Compare the received signature with the expected signature
	if receivedSignature != expectedSignature {
		return errSignatureDoesNotMatch
	}

	return ErrNone
}
