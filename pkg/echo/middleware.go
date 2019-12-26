package echo

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/pratikmallya/clubhouse-webhook/pkg/signature"
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

			verified, err := signature.Verify(c.Request(), config.Key)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("webhook verification failed. Error: %w", err))
			}
			if !verified {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("webhook signature did not match"))
			}
			return next(c)
		}
	}
}
