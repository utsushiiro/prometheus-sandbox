package metrics

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// CreateEndpointHandler : metricsをexposeするハンドラを返す
func CreateEndpointHandler() echo.HandlerFunc {
	h := promhttp.Handler()
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// ExposeEndpointMiddleware : /metricsでmetricsをexposeするミドルウェア
func ExposeEndpointMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	h := promhttp.Handler()
	return func(c echo.Context) error {
		if c.Path() == "metrics" {
			h.ServeHTTP(c.Response(), c.Request())
			return nil
		}
		return next(c)
	}
}
