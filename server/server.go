package server

import "net/http"

// New returns new instance of Server
func New(listen string) *http.Server {
	return &http.Server{
		Addr: listen,
	}
}
