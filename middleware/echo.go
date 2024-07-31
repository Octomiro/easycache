package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/octomiro/easycache"
)

type MultiReponseWriter struct {
	Writer         io.Writer
	ResponseWriter http.ResponseWriter
}

func (w *MultiReponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *MultiReponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *MultiReponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func EchoCacheMiddleware(ec *easycache.EasyCache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			path := c.Path()
			if _, ok := ec.IgnoreEndpoints[path]; ok {
				return next(c)
			}

			key := c.Request().URL.String()

			val, ok := ec.Cache().Get(key)
			if ok {
				log.Printf("cache hit for %s", key)
				resp, ok := val.(easycache.Response)
				if !ok {
					log.Println("could not read cached response")
					return next(c)
				}

				for k, v := range resp.Header {
					c.Response().Header().Set(k, strings.Join(v, ","))
				}

				_, err := c.Response().Write(resp.Response)
				if err != nil {
					log.Println("could not write cached response to client")
					return next(c)
				}
				return err
			}
			log.Printf("cache miss for %s", key)

			// Capture Response
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &MultiReponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			err := next(c)
			if err != nil {
				return err
			}

			if c.Response().Status < ec.CacheIfStatusCodeLessThan {
				ec.Cache().SetDefault(key, easycache.Response{
					Header:   c.Response().Header(),
					Response: resBody.Bytes(),
				})
			}

			return err
		}
	}
}
