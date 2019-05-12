package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type (
	controller struct {
		chi.Router
	}
)

func newController() http.Handler {
	api := chi.NewMux()
	controller := &controller{
		Router: api,
	}
	api.Post("/upload/form", controller.uploadFORM)
	return controller
}

func (c *controller) uploadFORM(w http.ResponseWriter, r *http.Request) {
	log.Printf("uploadFORM not implemented")
}
