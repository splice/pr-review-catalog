package requestid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Since `context.WithValue` can have collisions if we use simple types, it's
// better to create a private struct that will only be used for this and always
// use the same variable. This way we avoid collisions.
type requestKey struct{}

// HeaderKey is the HTTP header we'll use to set the Request ID.
const HeaderKey = "X-Request-Id"

var reqIDKey = &requestKey{}

// FromContext returns the request ID contained in the given context.
func FromContext(ctx context.Context) string {
	v := ctx.Value(reqIDKey)
	s, _ := v.(string)
	return s
}

// AddToContext returns a copy of the given context with the request ID.
func AddToContext(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, reqIDKey, reqID)
}

// Middleware returns an http.Handler that will add a request ID to a
// HTTP request's context and will also set the `X-Request-Id` header.
func Middleware(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqUUID string

		// Only generate a new request ID if it's not present in the Header
		if reqUUID = r.Header.Get(HeaderKey); reqUUID == "" {
			reqUUID = uuid.New().String()
		}

		w.Header().Set(HeaderKey, reqUUID)
		h.ServeHTTP(w, r.WithContext(AddToContext(r.Context(), reqUUID)))
	})
}
