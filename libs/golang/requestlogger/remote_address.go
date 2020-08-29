package requestlogger

import "net/http"

// RemoteAddrHeaders are the well-known headers we'll look at to get the
// remote user's address of a request.
var RemoteAddrHeaders = []string{
	"Cf-Connecting-Ip", // cloudflare - Must be first
	"Client-Ip",
	"X-Real-Ip",       // nginx proxy
	"X-Forwarded-For", // common proxies
	"X-Forwarded",
	"X-Cluster-Client-Ip",
	"Forwarded-For",
	"Forwarded",
}

// RemoteAddr makes a best-effort attempt at figuring out the client
// remote address. It starts by looking at well-known http headers,
// and falls back on http.Request.RemoteAddr.
func RemoteAddr(r *http.Request) string {
	for _, h := range RemoteAddrHeaders {
		val := r.Header.Get(h)
		if val != "" {
			return val
		}
	}

	return r.RemoteAddr
}
