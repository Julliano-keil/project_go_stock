package subcategory

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

type moduleSubcategory struct {
	useCase domain.SubcategoryUseCase
	name    string
	path    string
}

// NewSubcategoryModule cria o módulo de subcategoria.
func NewSubcategoryModule(useCase domain.SubcategoryUseCase) modules.AppModule {
	return &moduleSubcategory{
		useCase: useCase,
		name:    "Subcategory module",
		path:    "/subcategorias",
	}
}

func (m *moduleSubcategory) Name() string { return m.name }
func (m *moduleSubcategory) Path() string { return m.path }

func (m *moduleSubcategory) Setup(r *mux.Router) *mux.Router {
	handlers := []modules.AppModuleHandler{
		{
			Handler: m.list,
			Path:    "",
			Label:   "Lista todas as subcategorias",
			Methods: []string{http.MethodGet},
		},
		{
			Handler: m.create,
			Path:    "",
			Label:   "Cadastra nova subcategoria",
			Methods: []string{http.MethodPost},
		},
		{
			Handler: m.getByID,
			Path:    "/{id}",
			Label:   "Busca subcategoria por ID",
			Methods: []string{http.MethodGet},
		},
	}

	for _, h := range handlers {
		r.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}
	return r
}

func (m *moduleSubcategory) list(w http.ResponseWriter, r *http.Request) {
	list, err := m.useCase.ListSubcategories(r.Context())
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (m *moduleSubcategory) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "id inválido"})
		return
	}
	sub, err := m.useCase.GetSubcategoryByID(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	if sub == nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "subcategoria não encontrada"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(sub)
}

type createSubcategoriaRequest struct {
	IDCategoria int64  `json:"id_categoria"`
	Nome        string `json:"nome"`
}

func (m *moduleSubcategory) create(w http.ResponseWriter, r *http.Request) {
	var req createSubcategoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: "corpo da requisição inválido"})
		return
	}
	if req.Nome == "" || req.IDCategoria <= 0 {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "nome e id_categoria são obrigatórios"})
		return
	}

	sub, err := m.useCase.Create(r.Context(), req.IDCategoria, req.Nome)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sub)
}
