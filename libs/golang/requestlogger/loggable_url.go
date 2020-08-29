package requestlogger

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// LoggableURL returns a full URL sans-sensitive parameter values.
func LoggableURL(r *http.Request) string {
	qParams := LoggableParams(r.URL.Query())
	urlStr := fmt.Sprintf(
		"%s://%s%s?%s",
		r.URL.Scheme,
		r.URL.Host,
		r.URL.Path,
		qParams.Encode(),
	)
	return urlStr
}

// LoggableParams returns a set of Query Parameters with sensitive paramters
// stripped from the URL. This **will delete** the parameters from `params`.
// It will not make a copy of the object.
func LoggableParams(params url.Values) url.Values {
	for _, param := range []string{"password", "token", "nonce"} {
		params.Del(param)
		params.Del(strings.ToUpper(param))
	}

	for k := range params {
		if params.Get(k) == "" {
			params.Del(k)
		}
	}

	return params
}
