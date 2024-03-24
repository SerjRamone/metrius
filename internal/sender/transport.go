package sender

import "net/http"

// internalRoundTripper is a holder function to make the process of
// creating middleware a bit easier without requiring the consumer to
// implement the RoundTripper interface.
type internalRoundTripper func(*http.Request) (*http.Response, error)

// RoundTrip is a RoundTripper interface implementation
func (rt internalRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// middleware is our middleware creation functionality.
type middleware func(http.RoundTripper) http.RoundTripper

// Chain is a handy function to wrap a base RoundTripper (optional)
// with the middlewares.
func chain(rt http.RoundTripper, middlewares ...middleware) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	for _, m := range middlewares {
		rt = m(rt)
	}

	return rt
}
