# EasyCache
easycache is simple go-cache wrapper for caching API requests for golang web frameworks (echo, fiber, net/http, etc). This package
provides easy-to-plug middlewares for different libraries. The current version supports only the echo/v4 framework.

With EasyCache, you can easily define the TTL, Interval of cleanups, endpoints to ignore, and the interval of status codes
that result in caching a response.

Configuration:
```go
type CacheConfig struct {
	// TTL for each cache entry
	TimeToLive int
	// Interval of cache evictions
	CleanUpInterval int
	// Http Status code limit for which responses are cached. Defaults to 300
	CacheIfStatusCodeLessThan int
	// List of endpoints (Paths) to ignore
	IgnoreEndpoints map[string]interface{}
}
```

To create an easycache instance:
```go
c := easycache.NewCache(easycache.CacheConfig{})
```

## Middlewares
### Echo Middleware
```go
e.GET("/home/:id", func(c echo.Context) error {
        cacheMiss = true
        return c.JSON(http.StatusOK, map[string]string{
                "info": "done",
                })
        }, middleware.EchoCacheMiddleware(&c))
```
