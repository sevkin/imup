package docs

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewSwagger returns swagger handler
// api -> /api/v1
// swagurl -> https://example.com:5566/swagger  /doc.json
func NewSwagger(api, swagurl string) (string, http.Handler) {
	route := "/swagger"
	mux := chi.NewMux()

	SwaggerInfo.BasePath = api

	url, err := url.Parse(swagurl)
	if err == nil {
		route = url.Path
		// TODO SwaggerInfo.Host
	}

	mux.Get("/*", httpSwagger.Handler(httpSwagger.URL(swagurl+"/doc.json")))
	return route, mux
}
