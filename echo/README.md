# echo middleware

Middleware to use in a server using the [echo framework].

Example usage:
```go
package main

import (
	"github.com/labstack/echo/v4"

	echoMiddleWare "github.com/pratikmallya/clubhouse-webhook/echo"
)

func main() {
	e := echo.New()
	e.Use(echoMiddleWare.HeaderVerificationMiddleware(echoMiddleWare.NewConfig(testSecretClubhouse), nil))
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	return e
}
```

[echo framework]: https://github.com/labstack/echo