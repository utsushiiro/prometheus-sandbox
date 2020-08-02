package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/utsushiiro/prometheus-sandbox/sample-app/metrics"
)

func init() {
	// for randomHandler
	rand.Seed(time.Now().UnixNano())
}

func randomCodeHandler(c echo.Context) error {
	v := rand.Intn(100)
	switch {
	case 0 <= v && v < 60:
		return c.String(http.StatusOK, http.StatusText(http.StatusOK))
	case 60 <= v && v < 85:
		return c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	case 85 <= v && v < 100:
		return c.String(http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
	default:
		panic("system error")
	}
}

func main() {
	e := echo.New()

	e.GET("/", randomCodeHandler, metrics.Measure)
	e.GET("/metrics", metrics.CreateEndpointHandler())

	e.Logger.Fatal(e.Start(":8080"))
}
