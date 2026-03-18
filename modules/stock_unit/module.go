package stock_unit

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

type moduleStockUnit struct {
	useCase domain.StockUnitUseCase
	name    string
	path    string
}

// NewStockUnitModule cria o módulo de unidades de estoque.
func NewStockUnitModule(useCase domain.StockUnitUseCase) modules.AppModule {
	return &moduleStockUnit{
		useCase: useCase,
		name:    "StockUnit module",
		path:    "/unidades_de_estoque",
	}
}

func (m *moduleStockUnit) Name() string { return m.name }
func (m *moduleStockUnit) Path() string { return m.path }

func (m *moduleStockUnit) Setup(r *mux.Router) *mux.Router {
	handlers := []modules.AppModuleHandler{
		{
			Handler: m.list,
			Path:    "/list",
			Label:   "Lista todas as unidades de estoque",
			Methods: []string{http.MethodGet},
		},
		{
			Handler: m.create,
			Path:    "/create",
			Label:   "Cadastra nova unidade de estoque",
			Methods: []string{http.MethodPost},
		},
		{
			Handler: m.getByID,
			Path:    "/get/{id}",
			Label:   "Busca unidade de estoque por ID",
			Methods: []string{http.MethodGet},
		},
		{
			Handler: m.update,
			Path:    "/update/{id}",
			Label:   "Atualiza unidade de estoque por ID",
			Methods: []string{http.MethodPut},
		},
		{
			Handler: m.delete,
			Path:    "/delete/{id}",
			Label:   "Remove unidade de estoque por ID",
			Methods: []string{http.MethodDelete},
		},
	}

	for _, h := range handlers {
		r.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}
	return r
}

func (m *moduleStockUnit) list(w http.ResponseWriter, r *http.Request) {
	list, err := m.useCase.ListStockUnits(r.Context())
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (m *moduleStockUnit) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "id inválido"})
		return
	}

	su, err := m.useCase.GetStockUnitByID(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	if su == nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "unidade de estoque não encontrada"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(su)
}

type stockUnitRequest struct {
	Nome string `json:"nome"`
}

func (m *moduleStockUnit) create(w http.ResponseWriter, r *http.Request) {
	var req stockUnitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: "corpo da requisição inválido"})
		return
	}
	if req.Nome == "" {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "nome é obrigatório"})
		return
	}

	su, err := m.useCase.Create(r.Context(), req.Nome)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(su)
}

func (m *moduleStockUnit) update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "id inválido"})
		return
	}

	var req stockUnitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: "corpo da requisição inválido"})
		return
	}
	if req.Nome == "" {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "nome é obrigatório"})
		return
	}

	su, err := m.useCase.Update(r.Context(), id, req.Nome)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(su)
}

func (m *moduleStockUnit) delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "id inválido"})
		return
	}

	if err := m.useCase.Delete(r.Context(), id); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

