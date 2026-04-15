package stock_movement

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

type moduleStockMovement struct {
	useCase domain.StockMovementUseCase
	name    string
	path    string
}

func NewStockMovementModule(useCase domain.StockMovementUseCase) modules.AppModule {
	return &moduleStockMovement{
		useCase: useCase,
		name:    "StockMovement module",
		path:    "/stock_movements",
	}
}

func (m *moduleStockMovement) Name() string { return m.name }
func (m *moduleStockMovement) Path() string { return m.path }

func (m *moduleStockMovement) Setup(r *mux.Router) *mux.Router {
	handlers := []modules.AppModuleHandler{
		{
			Handler: m.list,
			Path:    "/list",
			Label:   "Lista todos os movimentos de estoque",
			Methods: []string{http.MethodGet},
		},
		{
			Handler: m.create,
			Path:    "/create",
			Label:   "Cria um movimento de estoque",
			Methods: []string{http.MethodPost},
		},
		{
			Handler: m.createItemEstoque,
			Path:    "/item_estoque/create",
			Label:   "Cadastra item de estoque",
			Methods: []string{http.MethodPost},
		},
		{
			Handler: m.getByID,
			Path:    "/get/{id}",
			Label:   "Busca movimento de estoque por ID",
			Methods: []string{http.MethodGet},
		},
		{
			Handler: m.rankingCategorias,
			Path:    "/stats/ranking_categorias",
			Label:   "Ranking de categorias com maior movimentação",
			Methods: []string{http.MethodGet},
		},
	}

	for _, h := range handlers {
		r.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}
	return r
}

func (m *moduleStockMovement) list(w http.ResponseWriter, r *http.Request) {
	list, err := m.useCase.ListStockMovements(r.Context())
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

type createStockMovementRequest struct {
	TipoMovimentacao string  `json:"tipo_movimentacao"`
	IDUser           int64   `json:"id_user"`
	IDsItemEstoque   []int64 `json:"ids_item_estoque"`
}

func (m *moduleStockMovement) create(w http.ResponseWriter, r *http.Request) {
	var req createStockMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: "corpo da requisição inválido"})
		return
	}

	sm, err := m.useCase.CreateStockMovement(
		r.Context(),
		req.TipoMovimentacao,
		req.IDUser,
		req.IDsItemEstoque,
	)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sm)
}

type createItemEstoqueRequest struct {
	IDEquipamento int64  `json:"id_equipamento"`
	StatusCode    string `json:"status_code"`
	Codigo        string `json:"codigo"`
}

func (m *moduleStockMovement) createItemEstoque(w http.ResponseWriter, r *http.Request) {
	var req createItemEstoqueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: "corpo da requisição inválido"})
		return
	}

	item, err := m.useCase.CreateItemEstoque(
		r.Context(),
		req.IDEquipamento,
		req.StatusCode,
		req.Codigo,
	)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

func (m *moduleStockMovement) rankingCategorias(w http.ResponseWriter, r *http.Request) {
	ranking, err := m.useCase.RankCategoriesByMovement(r.Context())
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ranking)
}

func (m *moduleStockMovement) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "id inválido"})
		return
	}

	sm, err := m.useCase.GetStockMovementByID(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: err.Error()})
		return
	}
	if sm == nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "movimentação não encontrada"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(sm)
}
