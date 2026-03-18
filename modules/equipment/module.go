package equipment

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

type moduleEquipment struct {
	useCase domain.EquipmentUseCase
	name    string
	path    string
}

// NewEquipmentModule cria o módulo de equipamentos.
func NewEquipmentModule(useCase domain.EquipmentUseCase) modules.AppModule {
	return &moduleEquipment{
		useCase: useCase,
		name:    "Equipment module",
		path:    "/equipamentos",
	}
}

func (m *moduleEquipment) Name() string { return m.name }
func (m *moduleEquipment) Path() string { return m.path }

func (m *moduleEquipment) Setup(r *mux.Router) *mux.Router {
	handlers := []modules.AppModuleHandler{
		{
			Handler: m.list,
			Path:    "/list",
			Label:   "Lista todos os equipamentos",
			Methods: []string{http.MethodGet},
		},
		{
			Handler: m.create,
			Path:    "/create",
			Label:   "Cadastra novo equipamento",
			Methods: []string{http.MethodPost},
		},
		{
			Handler: m.getByID,
			Path:    "/get/{id}",
			Label:   "Busca equipamento por ID",
			Methods: []string{http.MethodGet},
		},
		{
			Handler: m.update,
			Path:    "/update/{id}",
			Label:   "Atualiza equipamento por ID",
			Methods: []string{http.MethodPut},
		},
		{
			Handler: m.delete,
			Path:    "/delete/{id}",
			Label:   "Remove equipamento por ID",
			Methods: []string{http.MethodDelete},
		},
	}

	for _, h := range handlers {
		r.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}
	return r
}

func (m *moduleEquipment) list(w http.ResponseWriter, r *http.Request) {
	list, err := m.useCase.ListEquipment(r.Context())
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (m *moduleEquipment) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "id inválido"})
		return
	}

	eq, err := m.useCase.GetEquipmentByID(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	if eq == nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "equipamento não encontrado"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(eq)
}

type equipmentRequest struct {
	Nome             string `json:"nome"`
	IDSubCategoria   int64  `json:"id_sub_categoria"`
	IDUnidadeEstoque int64  `json:"id_unidade_estoque"`
}

func (m *moduleEquipment) create(w http.ResponseWriter, r *http.Request) {
	var req equipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: "corpo da requisição inválido"})
		return
	}
	if req.Nome == "" || req.IDSubCategoria <= 0 || req.IDUnidadeEstoque <= 0 {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "nome, id_sub_categoria e id_unidade_estoque são obrigatórios"})
		return
	}

	eq, err := m.useCase.Create(r.Context(), req.Nome, req.IDSubCategoria, req.IDUnidadeEstoque)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(eq)
}

func (m *moduleEquipment) update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "id inválido"})
		return
	}

	var req equipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: "corpo da requisição inválido"})
		return
	}
	if req.Nome == "" || req.IDSubCategoria <= 0 || req.IDUnidadeEstoque <= 0 {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "nome, id_sub_categoria e id_unidade_estoque são obrigatórios"})
		return
	}

	eq, err := m.useCase.Update(r.Context(), id, req.Nome, req.IDSubCategoria, req.IDUnidadeEstoque)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(eq)
}

func (m *moduleEquipment) delete(w http.ResponseWriter, r *http.Request) {
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

