package modules

import (
	"net/http"

	"github.com/gorilla/mux"
)

type AppModule interface {
	Name() string
	Path() string
	Setup(r *mux.Router) *mux.Router
}

type AppModuleHandler struct {
	Path    string
	Label   string
	Handler http.HandlerFunc
	Methods []string
}
