package middleware_test

import (
	"context"
	"io"
	"net/http"
	"testing"

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

	go func() {
		_ = e.Start(":8001")
	}()
	defer func() {
		_ = e.Shutdown(context.Background())
	}()

	e.GET("/home/:id", func(c echo.Context) error {
		cacheMiss = true
		return c.JSON(http.StatusOK, map[string]string{
			"info": "done",
		})
	}, middleware.EchoCacheMiddleware(&c))

	resp, err := http.Get("http://localhost:8001/home/123?query=321")
	assert.NoError(t, err)
	assert.Equal(t, true, cacheMiss)

	content, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	respStr := string(content)

	require.Equal(t, "{\"info\":\"done\"}\n", respStr)

	cacheMiss = false
	resp, err = http.Get("http://localhost:8001/home/123?query=321")
	assert.NoError(t, err)
	assert.Equal(t, false, cacheMiss)

	content, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, respStr, string(content))
}
