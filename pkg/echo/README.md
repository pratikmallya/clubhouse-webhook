# echo middleware

Middleware to use in a server using the [echo framework].

Example usage:
```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	echoMiddleWare "github.com/pratikmallya/clubhouse-webhook/echo"
)

func main() {
	e := echo.New()
	// testSecretClubhouse is the secret that you configure when setting up webhooks in clubhouse
	testSecretClubhouse := "some-secret"
	e.Use(echoMiddleWare.HeaderVerification(echoMiddleWare.NewConfig(testSecretClubhouse), nil))
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
```

[echo framework]: https://github.com/labstack/echo