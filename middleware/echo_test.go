package middleware_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/octomiro/easycache"
	"github.com/octomiro/easycache/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEchoMiddleware(t *testing.T) {
	cacheMiss := false

	c := easycache.NewCache(easycache.CacheConfig{})
	e := echo.New()
	e.HideBanner = true

	e.GET("/home/:id", func(c echo.Context) error {
		cacheMiss = true
		return c.JSON(http.StatusOK, map[string]string{
			"info": "done",
		})
	}, middleware.EchoCacheMiddleware(&c))

	ts := httptest.NewServer(e)
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/home/123?query=321", ts.URL))
	assert.NoError(t, err)
	assert.Equal(t, true, cacheMiss)

	content, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	respStr := string(content)

	require.Equal(t, "{\"info\":\"done\"}\n", respStr)

	cacheMiss = false
	resp, err = http.Get(fmt.Sprintf("%s/home/123?query=321", ts.URL))
	assert.NoError(t, err)
	assert.Equal(t, false, cacheMiss)

	content, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, respStr, string(content))
}

func BenchmarkEchoMiddleware(b *testing.B) {
	c := easycache.NewCache(easycache.CacheConfig{})
	e := echo.New()
	e.HideBanner = true

	e.GET("/home/:id", func(c echo.Context) error {
		time.Sleep(100 * time.Millisecond)
		return c.JSON(http.StatusOK, map[string]string{
			"id":    c.Param("id"),
			"query": c.QueryParam("query"),
			"info":  "done",
		})
	}, middleware.EchoCacheMiddleware(&c))

	ts := httptest.NewServer(e)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 200; i++ {
			_, err := http.Get(fmt.Sprintf("%s/home/123?query=321", ts.URL))
			require.NoError(b, err)
		}
	}
}

func BenchmarkEchoWithoutMiddleware(b *testing.B) {
	e := echo.New()
	e.HideBanner = true

	e.GET("/home/:id", func(c echo.Context) error {
		time.Sleep(100 * time.Millisecond)
		return c.JSON(http.StatusOK, map[string]string{
			"id":    c.Param("id"),
			"query": c.QueryParam("query"),
			"info":  "done",
		})
	})

	ts := httptest.NewServer(e)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 200; i++ {
			_, err := http.Get(fmt.Sprintf("%s/home/123?query=321", ts.URL))
			require.NoError(b, err)
		}
	}
}
