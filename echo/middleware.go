package echo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	HeaderClubHouseSignature = "Clubhouse-Signature"
)

type Config struct {
	// Key is the secret used for generating HMAC digest. This is the secret provided to Cluhouse when configuring
	// webhooks.
	Key []byte
}

func NewConfig(key string) Config {
	return Config{
		Key: []byte(key),
	}
}

// HeaderVerification is a verification middleware that only allows requests that are verified to originate
// from Clubhouse when generated with a secret.
// From https://clubhouse.io/api/webhook/v1/#Signature:
// If you provide a secret when you create the Outgoing Webhook, it will include an HTTP header named
// Clubhouse-Signature. The value of this header is a cryptographic hash encoded in hexadecimal.
//
// The signature is computed by the HMAC-SHA-256 algorithm. The ‘message’ is the HTTP request body encoded in UTF-8.
// The ‘secret’ is the secret string you provided, also encoded in UTF-8.
func HeaderVerification(config Config, skipper middleware.Skipper) echo.MiddlewareFunc {

	if skipper == nil {
		skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			hexMessageMAC := c.Request().Header.Get(HeaderClubHouseSignature)
			if hexMessageMAC == "" {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s header not specified", HeaderClubHouseSignature))
			}

			messageMAC := make([]byte, hex.DecodedLen(len(hexMessageMAC)))
			if _, err := hex.Decode(messageMAC, []byte(hexMessageMAC)); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error when decoding header %s: %v", HeaderClubHouseSignature, err))
			}

			message, err := ioutil.ReadAll(c.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error occured when reading request body: %w", err))
			}
			if !validMAC(message, messageMAC, config.Key) {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unable to verify request origin from clubhouse"))
			}

			// Restore request body for processing down the chain
			// TODO: is this  the recommended way to restore request body
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(message))
			return next(c)
		}
	}
}

// copied verbatim from https://golang.org/pkg/crypto/hmac/
func validMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
