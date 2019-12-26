package features

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/labstack/echo/v4"

	echo2 "github.com/pratikmallya/clubhouse-webhook/pkg/echo"
	"github.com/pratikmallya/clubhouse-webhook/pkg/signature"
)

const (
	testSecretClubhouse = "the_cake_is_a_lie"
	validRequestBody    = "bobloblaw"
)

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

var (
	req    *http.Request
	server *echo.Echo
	rec    *httptest.ResponseRecorder
	resp   *http.Response
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
		Paths:  []string{"../features"},
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^request does not have clubhouse header$`, func() error {
		return nil
	})

	s.Step(`^request is rejected with status code (\d+)$`, func(code int) error {
		resp = rec.Result()
		if resp.StatusCode != code {
			return fmt.Errorf("expected %d, got %d", code, resp.StatusCode)
		}
		return nil
	})

	s.Step(`^request is accepted with status code (\d+)$`, func(code int) error {
		resp = rec.Result()
		if resp.StatusCode != code {
			return fmt.Errorf("expected %d, got %d", code, resp.StatusCode)
		}
		return nil
	})

	s.Step(`^request has a garbage clubhouse header$`, func() error {
		req.Header.Set(signature.HeaderClubHouseSignature, "jar-jar-binks")
		return nil
	})

	s.Step(`^request has a valid clubhouse header$`, func() error {
		setValidHeader(req)
		return nil
	})

	s.Step(`^request has empty body$`, func() error {
		return nil
	})

	s.Step(`^request has a valid body$`, func() error {
		req.Body = ioutil.NopCloser(bytes.NewReader([]byte(validRequestBody)))
		return nil
	})

	s.Step(`^request signature does not match request$`, func() error {
		setValidHeader(req)
		return nil
	})

	s.Step(`^request signature does match request$`, func() error {
		mac := hmac.New(sha256.New, []byte(testSecretClubhouse))
		mac.Write([]byte(validRequestBody))
		dst := make([]byte, hex.EncodedLen(len(mac.Sum(nil))))
		hex.Encode(dst, mac.Sum(nil))
		req.Header.Set(signature.HeaderClubHouseSignature, string(dst))
		req.Body = ioutil.NopCloser(bytes.NewReader([]byte(validRequestBody)))
		return nil
	})

	s.Step(`^request is made$`, func() error {
		server.ServeHTTP(rec, req)
		return nil
	})

	s.BeforeScenario(func(interface{}) {
		server = testserver() // create a new echo server for every scenario
		rec = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "http://fake", nil)
	})

}

func setValidHeader(req *http.Request) {
	// $ echo "jar-jar-binks" | xxd -p -u
	req.Header.Set(signature.HeaderClubHouseSignature, "6A61722D6A61722D62696E6B730A")
}

func testserver() *echo.Echo {
	e := echo.New()
	e.Use(echo2.HeaderVerification(echo2.NewConfig(testSecretClubhouse), nil))
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	return e
}
