package server

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Request represents a HTTP request.
type Request struct {
	// R represents the default http Request.
	R *http.Request

	// Params represents the Query params.
	QueryParams map[string][]string

	// MiddlewareParams represents the parameters which are added by the middleware.
	MiddlewareParams map[string]interface{}

	// Optional user data for higher level library wrappers.
	UserData map[string]interface{}
}

//GetURLParam returns the URL parameters that match with the given name.
func (r *Request) GetURLParam(name string) string {
	return chi.URLParam(r.R, name)
}

// parseQueryParams parses the URL query parameters and adds them the the map in the Request.
func (r *Request) parseQueryParams() {
	queryMap := r.R.URL.Query()

	for key, val := range queryMap {
		r.QueryParams[key] = val
	}
}

// NewRequest creates a new and empty Request wrapped around the native Go request.
func newRequest(r *http.Request) Request {
	req := Request{
		R:                r,
		QueryParams:      map[string][]string{},
		MiddlewareParams: map[string]interface{}{},
		UserData:         map[string]interface{}{},
	}

	req.parseQueryParams()
	return req
}

// Domain returns the main domain and top level domain from the URL.
func (r *Request) TopDomain() (string, error) {
	if r == nil {
		return "", &ErrNil{}
	}

	if strings.Contains(r.R.Host, "localhost") || strings.Contains(r.R.Host, "127.0.0.1") {
		return "localhost", nil
	}

	splitted := strings.Split(r.R.Host, ".")

	if len(splitted) > 2 {
		mainDomain := splitted[len(splitted)-2]
		tld := splitted[len(splitted)-1]

		return mainDomain + "." + tld, nil
	}

	return "", &ErrInvalidURL{r.R.RequestURI}
}

// Cookie returns the value of a cookie with a given name.
func (r *Request) Cookie(name string) (*string, error) {
	if r == nil {
		return nil, &ErrNil{}
	}

	for _, c := range r.R.Cookies() {
		if c.Name == name {
			return &c.Value, nil
		}
	}

	return nil, &ErrCookieNotFound{name}
}
