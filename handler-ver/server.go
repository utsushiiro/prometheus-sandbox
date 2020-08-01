package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestCounterVec *prometheus.CounterVec

func init() {
	requestCounterVec = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sample_app",
			Name:      "http_requests_total",
			Help:      "How many HTTP requests processed, partitioned by status code.",
		},
		[]string{"code"},
	)

	requestCounterVec.WithLabelValues("200").Add(0)
	requestCounterVec.WithLabelValues("404").Add(0)
	requestCounterVec.WithLabelValues("503").Add(0)
}

func createPrometheusHandler() echo.HandlerFunc {
	h := promhttp.Handler()
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func okHandler(c echo.Context) error {
	requestCounterVec.WithLabelValues("200").Inc()
	return c.String(http.StatusOK, http.StatusText(http.StatusOK))
}

func notFoundHandler(c echo.Context) error {
	requestCounterVec.WithLabelValues("404").Inc()
	return c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func serviceUnavailableHandler(c echo.Context) error {
	requestCounterVec.WithLabelValues("503").Inc()
	return c.String(http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
}

func main() {
	e := echo.New()

	e.GET("/metrics", createPrometheusHandler())

	e.GET("/", func(c echo.Context) error {
		requestCounterVec.WithLabelValues("200").Inc()
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/200", okHandler)
	e.GET("/404", notFoundHandler)
	e.GET("/503", serviceUnavailableHandler)

	rand.Seed(time.Now().UnixNano())
	e.GET("/random", func(c echo.Context) error {
		v := rand.Intn(100)
		switch {
		case 0 <= v && v < 60:
			return okHandler(c)
		case 60 <= v && v < 80:
			return notFoundHandler(c)
		case 80 <= v && v < 100:
			return serviceUnavailableHandler(c)
		default:
			panic("system error")
		}
	})

	e.Logger.Fatal(e.Start(":8080"))
}
