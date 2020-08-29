package requestlogger

import (
	"context"
	"net/http"
	"os"

	"github.com/felixge/httpsnoop"
	"github.com/sirupsen/logrus"
	"github.com/splice/catalog-interview/libs/golang/requestid"
	datadog "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Since `context.WithValue` can have collisions if we use simple types, it's
// better to create a private struct that will only be used for this and always
// use the same variable. This way we avoid collisions.
type requestKey struct{}

var reqLoggerKey = &requestKey{}

// FromContext returns the logger contained in the given context. A second
// return value indicates if a logger was found in the context. If no logger is
// found, a new logger is returned.
func FromContext(ctx context.Context) (*logrus.Entry, bool) {
	v := ctx.Value(reqLoggerKey)
	if s, ok := v.(*logrus.Entry); ok {
		return s, true
	}

	lg := logrus.New()
	lg.SetFormatter(&logrus.JSONFormatter{})
	return lg.WithContext(ctx), false
}

// ContextWithLogger returns a copy of the given context with the logger
func ContextWithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, reqLoggerKey, logger)
}

// Middleware returns an http.Handler that will add a logger to the request's
// context. If the request has a request ID, that will be used as the prefix.
//
// If you don't want to log some paths, you can add them as `skipPaths`. Make sure
// they're prefixed with `/`.
func Middleware(h http.Handler, serviceName string, skipPaths ...string) http.HandlerFunc {
	var shouldSkip = func(uri string) bool {
		for _, path := range skipPaths {
			if path == uri {
				return true
			}
		}

		return false
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, _ := datadog.SpanFromContext(r.Context())

		logger, _ := FromContext(r.Context())
		logger = logger.WithContext(r.Context()).WithFields(logrus.Fields{
			"service":     serviceName,
			"env":         os.Getenv("DATADOG_ENV"),
			"req_id":      requestid.FromContext(r.Context()),
			"dd.trace_id": span.Context().TraceID(),
			"dd.span_id":  span.Context().SpanID(),
		})

		metrics := httpsnoop.CaptureMetrics(h, w, r.WithContext(ContextWithLogger(r.Context(), logger)))
		if !shouldSkip(r.URL.Path) {
			logger = logger.WithFields(logrus.Fields{
				"request_method":  r.Method,
				"path":            r.URL.String(),
				"status":          metrics.Code,
				"body_bytes_sent": metrics.Written,
				"request_time":    metrics.Duration.Seconds(),
				"remote_address":  RemoteAddr(r),
				"referer":         r.Referer(),
				"http_user_agent": r.UserAgent(),
				"url":             LoggableURL(r),
			})

			if metrics.Code < http.StatusInternalServerError {
				logger.Printf("request completed")
			} else {
				logger.Error("request completed")
			}
		}
	})
}
