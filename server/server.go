package server

import (
	"imup/uploader"
	"net/http"

	"github.com/go-chi/chi"
)

// New returns new instance of Server
func New(listen string) *http.Server {
	mux := chi.NewMux()
	server := &http.Server{
		Addr:    listen,
		Handler: mux,
	}
	mux.Mount("/api/v1", newController(uploader.NewDirUploader()))
	return server
}
