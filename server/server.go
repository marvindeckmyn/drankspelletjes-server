package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// Handler is called when a new request came in.
type Handler func(ResponseWriter, *Request)

// Middleware can be added at the beginning of a request chain.
type Middleware func(ResponseWriter, *Request) bool

// Server is the base object required to set up an HTTP/ws API.
type Server struct {
	rtr    *chi.Mux
	mWares []Middleware
}

func (s *Server) GetRouter() (*chi.Mux, error) {
	if s == nil {
		return nil, &ErrNil{}
	}

	return s.rtr, nil
}

// New creates a new router instance.
func New() *Server {
	return &Server{
		rtr:    chi.NewRouter(),
		mWares: []Middleware{},
	}
}

// ListenAndServe starts the HTTP/ws server after which clients can connect.
func (s *Server) ListenAndServe(port uint16) error {
	if s == nil {
		return &ErrNil{}
	}

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.rtr)
	if err != nil {
		return &ErrListen{err}
	}

	return nil
}

// AddMiddleware adds middleware to be executed on every route.
func (s *Server) AddMiddleware(mWare Middleware) error {
	if s == nil {
		return &ErrNil{}
	}

	s.mWares = append(s.mWares, mWare)
	return nil
}

// runMiddlewares runs all the registered middlewares until one of the middlewares stops the
// execution chain. When the one of the middlewares signals the execution chain to stop this
// function will return false.
func runMiddlewares(mWares []Middleware, rw ResponseWriter, r Request) bool {
	if mWares == nil {
		return true
	}

	for _, mWare := range mWares {
		if mWare != nil && !mWare(rw, &r) {
			return false
		}
	}

	return true
}

// httpRouterHandle returns a Handle function which is compatible with httprouter.
func (s *Server) httpRouterHandle(callback Handler, mWares []Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r)
		rw := newRespWriter(w)

		if !runMiddlewares(s.mWares, rw, req) || !runMiddlewares(mWares, rw, req) {
			//TODO: check if response headers are written
			return
		}

		callback(rw, &req)
	}
}

// Get routes the http GET calls for the DPT router
func (s *Server) Get(url string, callback Handler, mWares ...Middleware) {
	s.rtr.Get(url, s.httpRouterHandle(callback, mWares))
}

// Put routes the http PUT calls for the DPT router
func (s *Server) Put(url string, callback Handler, mWares ...Middleware) {
	s.rtr.Put(url, s.httpRouterHandle(callback, mWares))
}

// Post routes the http POST calls for the DPT router
func (s *Server) Post(url string, callback Handler, mWares ...Middleware) {
	s.rtr.Post(url, s.httpRouterHandle(callback, mWares))
}

// Delete routes the http DELETE calls for the DPT router
func (s *Server) Delete(url string, callback Handler, mWares ...Middleware) {
	s.rtr.Delete(url, s.httpRouterHandle(callback, mWares))
}

// ServeFiles is used to send a file as response
func ServeFile(w http.ResponseWriter, r *http.Request, dir string, file string) {
	_, err := os.Stat(dir + file)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, dir+"/index.html")
		return
	}

	http.ServeFile(w, r, dir+file)
}

// ServeFiles will serve files prefixed with a certain URL from a certain directory.
func (s *Server) ServeFiles(url string, dir string, mWares ...Middleware) {
	if url[len(url)-1] != '/' {
		url += "/{filepath}"
	} else {
		url += "{filepath}"
	}

	s.rtr.Get(url, func(w http.ResponseWriter, r *http.Request) {
		req := newRequest(r)
		rw := newRespWriter(w)

		if !runMiddlewares(s.mWares, rw, req) || !runMiddlewares(mWares, rw, req) {
			return
		}

		// get url param filepath
		ServeFile(w, r, dir, chi.URLParam(r, "filepath"))
	})
}
