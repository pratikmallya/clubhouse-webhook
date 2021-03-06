package echo

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/pratikmallya/clubhouse-webhook/pkg/signature"
)

// Config specifies configuration options for the middleware.
type Config struct {
	// Key is the secret used for generating HMAC digest. This is the secret provided to Cluhouse when configuring
	// webhooks.
	Key []byte
}

// NewConfig returns a new Config with the provided configuration options.
func NewConfig(key string) Config {
	return Config{
		Key: []byte(key),
	}
}

// HeaderVerification is a verification middleware that only allows requests that are verified to originate
// from Clubhouse.
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
