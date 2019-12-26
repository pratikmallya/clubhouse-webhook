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
