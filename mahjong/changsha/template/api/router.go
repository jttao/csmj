package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	templatePath = "/template"
)

func Router(r *mux.Router) {
	r.Path(templatePath).HandlerFunc(http.HandlerFunc(handleTemplate))
}
