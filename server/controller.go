package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
)

type (
	controller struct {
		chi.Router
	}

	// Success response returns json with uploaded image uuid
	Success struct {
		UUID uuid.UUID `json:"uuid"` // uploaded image uuid
	}

	// Failed response returns json with explained error
	Failed struct {
		Message string `json:"message"` // explained error
	}
)

func success(w http.ResponseWriter, r *http.Request, UUID uuid.UUID) {
	render.JSON(w, r, &Success{UUID: UUID})
}

func failed(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		render.JSON(w, r, &Failed{Message: "something wrong"})
		return
	}
	render.JSON(w, r, &Failed{Message: err.Error()})
}

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

	// success(w, r, uuid.Must(uuid.NewV4()))
	failed(w, r, errors.New("not implemented"))
}
