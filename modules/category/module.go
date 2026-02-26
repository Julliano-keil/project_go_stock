package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"lince/domain"
	"lince/entities"
	"lince/httputil"
	"lince/modules"

	"github.com/gorilla/mux"
)

type moduleCategory struct {
	useCase domain.CategoryUseCase
	name    string
	path    string
}

// NewCategoryModule cria o módulo de categoria.
func NewCategoryModule(useCase domain.CategoryUseCase) modules.AppModule {
	return &moduleCategory{
		useCase: useCase,
		name:    "Category module",
		path:    "/categorias",
	}
}

func (m *moduleCategory) Name() string { return m.name }
func (m *moduleCategory) Path() string { return m.path }

func (m *moduleCategory) Setup(r *mux.Router) *mux.Router {
	handlers := []modules.AppModuleHandler{
		{
			Handler: m.list,
			Path:    "",
			Label:   "Lista todas as categorias",
			Methods: []string{http.MethodGet},
		},
		{
			Handler: m.create,
			Path:    "",
			Label:   "Cadastra nova categoria",
			Methods: []string{http.MethodPost},
		},
		{
			Handler: m.getByID,
			Path:    "/{id}",
			Label:   "Busca categoria por ID",
			Methods: []string{http.MethodGet},
		},
	}

	for _, h := range handlers {
		r.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}
	return r
}

func (m *moduleCategory) list(w http.ResponseWriter, r *http.Request) {
	list, err := m.useCase.ListCategories(r.Context())
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (m *moduleCategory) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "id inválido"})
		return
	}
	cat, err := m.useCase.GetCategoryByID(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	if cat == nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "categoria não encontrada"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(cat)
}

type createCategoriaRequest struct {
	Nome string `json:"nome"`
}

func (m *moduleCategory) create(w http.ResponseWriter, r *http.Request) {
	var req createCategoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: "corpo da requisição inválido"})
		return
	}
	if req.Nome == "" {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "nome é obrigatório"})
		return
	}

	cat, err := m.useCase.Create(r.Context(), req.Nome)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(cat)
}
