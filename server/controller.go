package server

import (
	"imup/uploader"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
)

type (
	controller struct {
		chi.Router
		uploader uploader.Uploader
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

func newController(uploader uploader.Uploader) http.Handler {
	api := chi.NewMux()
	controller := &controller{
		Router:   api,
		uploader: uploader,
	}
	api.Post("/upload/form", controller.uploadFORM)
	return controller
}

func (c *controller) uploadFORM(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files
	// TODO customizable set maxupload
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		failed(w, r, err)
		return
	}
	src, _, err := r.FormFile("image")
	if err != nil {
		failed(w, r, err)
		return
	}
	defer src.Close()

	uuid, err := c.uploader.Store(src)
	if err != nil {
		failed(w, r, err)
		return
	}
	success(w, r, uuid)
}
