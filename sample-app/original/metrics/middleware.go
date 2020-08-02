package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestCount *prometheus.CounterVec
var requestDuration *prometheus.HistogramVec
var requestSize *prometheus.SummaryVec
var responseSize *prometheus.SummaryVec

func init() {

	// 4XX, 5XXのような粒度がほしい場合どうするか
	// - PromQLで対処
	// 		- label名のマッチングの際に正規表現を利用する
	// 		- 主要なコードをすべて指定
	// - 計測で対処
	//		- code % 100の値をcodeとは別のラベルで付与する
	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_count_total",
		Help: "Counter of HTTP requests made.",
	}, []string{"code"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "sample_app",
		Name:      "http_request_duration_seconds",
		Help:      "A histogram of latencies for requests.",
		Buckets:   append([]float64{0.000001, 0.001, 0.003}, prometheus.DefBuckets...),
	}, []string{"code"})

	requestSize = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "sample_app",
		Name:      "http_request_size_bytes",
		Help:      "A summary of request sizes for requests.",
	}, []string{"code"})

	responseSize = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "sample_app",
		Name:      "http_response_size_bytes",
		Help:      "A summary of response sizes for requests.",
	}, []string{"code"})
}

// Measure : HTTPの各種メトリクスを計測するミドルウェア
func Measure(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		end := time.Now()

		// エラーが空でStatusが定まっている場合(値が0でない場合)にのみ各種値を計測する
		// このあたりの扱いはエラーハンドリングの方針に依存するため利用する際に調整が必要
		statusCode := c.Response().Status
		if err == nil && statusCode != 0 {
			requestCount.WithLabelValues(strconv.Itoa(statusCode)).Inc()
			requestDuration.WithLabelValues(strconv.Itoa(statusCode)).Observe(end.Sub(start).Seconds())
			requestSize.WithLabelValues(strconv.Itoa(statusCode)).Observe(float64(computeApproximateRequestSize(c.Request())))
			responseSize.WithLabelValues(strconv.Itoa(statusCode)).Observe(float64(c.Response().Size))
		}

		return err
	}
}

// https://github.com/prometheus/client_golang/blob/c650ae9fa1039588f417cbc526ddb8155ace7613/prometheus/promhttp/instrument_server.go#L298-L320
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s += len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
