package httpserver

import (
	"net/http"
)

type (
	Server struct {
		*http.Server
	}
)

var defaultHTTPAddress = "127.0.0.1:8080"

func New(handler http.Handler, opts ...Option) *Server {
	s := &Server{
		&http.Server{
			Handler: handler,
			Addr:    defaultHTTPAddress,
		},
	}
	for _, opt := range opts {
		opt(s)
	}

	return s
}
