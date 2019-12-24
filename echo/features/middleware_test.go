package features

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/labstack/echo/v4"

	echoMiddleWare "github.com/pratikmallya/clubhouse-webhook/echo"
)

const (
	testSecretClubhouse = "the_cake_is_a_lie"
)

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

var (
	req *http.Request
	server *echo.Echo
	rec  *httptest.ResponseRecorder
	resp *http.Response
)

func TestMain(m *testing.M) {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: format,
		Paths:     []string{"../features"},
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^request does not have clubhouse header$`, func() error {
		req = httptest.NewRequest("GET", "http://fake", nil)
		return nil
	})
	s.Step(`^request is made$`, func() error {
		server.ServeHTTP(rec, req)
		return nil
	})
	s.Step(`^request is rejected with status code (\d+)$`, func(code int) error {
		resp = rec.Result()
		if resp.StatusCode != code {
			return fmt.Errorf("expected %d, got %d", code, resp.StatusCode)
		}
		return nil
	})

	s.BeforeScenario(func(interface{}) {
		server = testserver() // create a new echo server for every scenario
		rec = httptest.NewRecorder()
	})

}

func testserver() *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleWare.HeaderVerificationMiddleware(echoMiddleWare.NewConfig(testSecretClubhouse), nil))
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	return e
}
