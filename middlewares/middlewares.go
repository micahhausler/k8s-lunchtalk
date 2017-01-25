package middlewares

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "Counter of requests broken out for each verb, path, and response code.",
		},
		[]string{"verb", "path", "code"},
	)
	requestLatencies = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_latencies",
			Help: "Response latency distribution in microseconds for each verb and path",
			// Use buckets ranging from 125 ms to 8 seconds.
			Buckets: prometheus.ExponentialBuckets(125000, 2.0, 7),
		},
		[]string{"verb", "path"},
	)
	requestLatenciesSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_latencies_summary",
			Help: "Response latency summary in microseconds for each verb and path.",
			// Make the sliding window of 1h.
			MaxAge: time.Hour,
		},
		[]string{"verb", "path"},
	)
)

func Register() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestLatencies)
	prometheus.MustRegister(requestLatenciesSummary)
}

func Monitor(verb, path string, httpCode int, reqStart time.Time) {
	elapsed := float64((time.Since(reqStart)) / time.Microsecond)
	requestCounter.WithLabelValues(verb, path, strconv.Itoa(httpCode)).Inc()
	requestLatencies.WithLabelValues(verb, path).Observe(elapsed)
	requestLatenciesSummary.WithLabelValues(verb, path).Observe(elapsed)
}

func init() {
	Register()
}

// Middleware is a type for decorating requests
type Middleware func(http.Handler) http.Handler

// Apply a list of middlewares to a handler
func Apply(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, adapter := range middlewares {
		h = adapter(h)
	}
	return h
}

// Logging is a middleware for adding a request log
func InstrumentRoute() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			wrappedWriter := &statusLoggingResponseWriter{w, http.StatusOK, 0}

			defer func() {
				Monitor(r.Method, r.URL.Path, wrappedWriter.status, now)
			}()
			h.ServeHTTP(wrappedWriter, r)

		})
	}
}

type statusLoggingResponseWriter struct {
	http.ResponseWriter
	status    int
	bodyBytes int
}

func (w *statusLoggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
func (w *statusLoggingResponseWriter) Write(data []byte) (int, error) {
	length, err := w.ResponseWriter.Write(data)
	w.bodyBytes += length
	return length, err
}

// Logging is a middleware for adding a request log
func Logging() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			wrappedWriter := &statusLoggingResponseWriter{w, http.StatusOK, 0}
			err := r.ParseForm()
			if err != nil {
				log.Printf("Error parsing form: %s", err.Error())
				http.Error(w, `{"error": "error parsing form"}`, http.StatusBadRequest)
				return
			}

			defer func() {
				reqLog := map[string]interface{}{
					"verb":        r.Method,
					"path":        r.URL.Path,
					"status":      strconv.Itoa(wrappedWriter.status),
					"remote_addr": r.RemoteAddr,
					"timestamp":   time.Now().Format(time.RFC3339),
					"body_bytes":  wrappedWriter.bodyBytes,
				}
				json.NewEncoder(os.Stdout).Encode(reqLog)
			}()
			h.ServeHTTP(wrappedWriter, r)

		})
	}
}
