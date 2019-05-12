package server

import (
	"encoding/base64"
	"encoding/json"
	"imup/uploader"
	"io/ioutil"
	"net/http"
	"strings"

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
		Error string `json:"error"` // explained error
	}

	// JSONImage {"image":"<base64encoded==>"}
	JSONImage struct {
		Image string `json:"image"` // base64 encoded content
	}
)

func success(w http.ResponseWriter, r *http.Request, UUID uuid.UUID) {
	render.JSON(w, r, &Success{UUID: UUID})
}

func failed(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		render.JSON(w, r, &Failed{Error: "something wrong"})
		return
	}
	render.JSON(w, r, &Failed{Error: err.Error()})
}

func newController(uploader uploader.Uploader) http.Handler {
	api := chi.NewMux()
	controller := &controller{
		Router:   api,
		uploader: uploader,
	}
	api.Post("/upload/form", controller.uploadFORM)
	api.Post("/upload/json", controller.uploadJSON)
	return controller
}

func (c *controller) store(w http.ResponseWriter, r *http.Request, src io.Reader) {
	uuid, err := c.uploader.Store(src)
	if err != nil {
		failed(w, r, err)
		return
	}
	success(w, r, uuid)
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

	c.store(w, r, src)
}

func (c *controller) uploadJSON(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err == nil {
		image := new(JSONImage)
		err = json.Unmarshal(buf, image)
		if err == nil {
			src := base64.NewDecoder(base64.StdEncoding, strings.NewReader(image.Image))
			c.store(w, r, src)
			return
		}
	}
	failed(w, r, err)
}
