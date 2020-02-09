package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	HeaderClubHouseSignature = "Clubhouse-Signature"
)

// Verify decides if the given request is a valid, signed request originating from Clubhouse.
//
// From https://clubhouse.io/api/webhook/v1/#Signature:
// If you provide a secret when you create the Outgoing Webhook, it will include an HTTP header named
// Clubhouse-Signature. The value of this header is a cryptographic hash encoded in hexadecimal.
//
// The signature is computed by the HMAC-SHA-256 algorithm. The ‘message’ is the HTTP request body encoded in UTF-8.
// The ‘secret’ is the secret string you provided, also encoded in UTF-8.
func Verify(req *http.Request, secret []byte) (bool, error) {

	hexMessageMAC := req.Header.Get(HeaderClubHouseSignature)
	if hexMessageMAC == "" {
		return false, fmt.Errorf("%s header not specified", HeaderClubHouseSignature)
	}

	messageMAC := make([]byte, hex.DecodedLen(len(hexMessageMAC)))
	if _, err := hex.Decode(messageMAC, []byte(hexMessageMAC)); err != nil {
		return false, fmt.Errorf("error when decoding header %s: %v", HeaderClubHouseSignature, err)
	}

	message, err := ioutil.ReadAll(req.Body)
	// Restore request body for processing down the chain
	// TODO: is this  the recommended way to restore request body?
	defer func() { req.Body = ioutil.NopCloser(bytes.NewBuffer(message)) }()
	if err != nil {
		return false, fmt.Errorf("error occured when reading request body: %w", err)
	}

	if !validMAC(message, messageMAC, secret) {
		return false, nil
	}

	return true, nil
}

// copied verbatim from https://golang.org/pkg/crypto/hmac/
func validMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
