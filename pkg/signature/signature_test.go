package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testSecretClubhouse = "the_cake_is_a_lie"
	validRequestBody    = "bobloblaw"
)

func TestVerifyValidSignaturePass(t *testing.T) {
	req := httptest.NewRequest("GET", "http://fake", nil)
	mac := hmac.New(sha256.New, []byte(testSecretClubhouse))
	mac.Write([]byte(validRequestBody))
	dst := make([]byte, hex.EncodedLen(len(mac.Sum(nil))))
	hex.Encode(dst, mac.Sum(nil))
	req.Header.Set(HeaderClubHouseSignature, string(dst))
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(validRequestBody)))

	verified, err := Verify(req, []byte(testSecretClubhouse))
	assert.NoError(t, err)
	assert.True(t, verified)
}

func TestRequestBodyPreserved(t *testing.T) {
	req := httptest.NewRequest("GET", "http://fake", nil)
	mac := hmac.New(sha256.New, []byte(testSecretClubhouse))
	mac.Write([]byte(validRequestBody))
	dst := make([]byte, hex.EncodedLen(len(mac.Sum(nil))))
	hex.Encode(dst, mac.Sum(nil))
	req.Header.Set(HeaderClubHouseSignature, string(dst))
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(validRequestBody)))

	Verify(req, []byte(testSecretClubhouse))

	body, _ := ioutil.ReadAll(req.Body)
	assert.Equal(t, []byte(validRequestBody), body)

}