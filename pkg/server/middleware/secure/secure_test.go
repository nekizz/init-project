package secure_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/nekizz/init-project/pkg/server/middleware/secure"
)

func echoHandler(mw ...echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	for _, v := range mw {
		e.Use(v)
	}
	e.GET("/hello", hwHandler)
	return e
}

func hwHandler(c echo.Context) error {
	return c.String(200, "Hello World")
}

func TestSecureHeaders(t *testing.T) {
	ts := httptest.NewServer(echoHandler(secure.Headers()))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/hello")
	if err != nil {
		t.Fatal("Did not expect http.Get to fail")
	}
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	// assert.Equal(t, "off", resp.Header.Get("X-DNS-Prefetch-Control"))
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	// assert.Equal(t, "max-age=5184000; includeSubDomains", resp.Header.Get("Strict-Transport-Security"))
	// assert.Equal(t, "noopen", resp.Header.Get("X-Download-Options"))
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
}

func TestCORS(t *testing.T) {
	ts := httptest.NewServer(echoHandler(secure.CORS(nil)))
	defer ts.Close()
	var cl http.Client
	req, _ := http.NewRequest("OPTIONS", ts.URL+"/hello", nil)
	resp, _ := cl.Do(req)
	t.Logf("%+v", resp.Header)
	assert.Equal(t, "OPTIONS, GET", resp.Header.Get("Allow"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	// assert.Equal(t, "86400", resp.Header.Get("Access-Control-Max-Age"))
	// assert.Equal(t, "Origin,Content-Type,Accept,Authorization", resp.Header.Get("Access-Control-Allow-Headers"))
	// assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
	// assert.Equal(t, "Content-Length", resp.Header.Get("Access-Control-Expose-Headers"))
	// assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}
