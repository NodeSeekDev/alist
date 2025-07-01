package signature

import (
	"strings"
)

// Parses signature version '2' header of the following form.
//
//	Authorization: AWS NXM00iAJpLvPYC94VRDd:q20e3MP/dZOpnzBZTvo619ONSpU=
func ParseSignV2(v2Auth string) (sv signV2Values, err ErrorCode) {
	if !strings.HasPrefix(v2Auth, signV2Algorithm) {
		return sv, ErrUnsupportAlgorithm
	}

	rawCred := strings.ReplaceAll(strings.TrimPrefix(v2Auth, signV2Algorithm), " ", "")

	signV2Values := signV2Values{}

	authFields := strings.Split(strings.TrimSpace(rawCred), ":")
	if len(authFields) != 2 {
		return sv, errMissingFields
	}

	credential := CredentialsV2{
		AccessKey: authFields[0],
	}
	signV2Values.Credential = credential
	signV2Values.Signature = authFields[1]

	// Return the structure here.
	return signV2Values, ErrNone
}
