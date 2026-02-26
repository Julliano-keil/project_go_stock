package user

import (
	"encoding/json"
	"net/http"

	"lince/domain"
	"lince/entities"
	"lince/httputil"
	"lince/middleware"
	"lince/modules"

	"github.com/gorilla/mux"
)

type moduleAuthentication struct {
	useCase    domain.UserUseCase
	cfg        entities.Config
	name       string
	path       string
	jwtSecret  string
}

// NewAuthenticationModule cria o módulo de autenticação (login + router base para rotas protegidas).
func NewAuthenticationModule(cfg entities.Config, useCase domain.UserUseCase) modules.AppModule {
	jwtSecret := cfg.JWTSecret
	if jwtSecret == "" {
		jwtSecret = "lince-secret-key"
	}
	return &moduleAuthentication{
		useCase:   useCase,
		cfg:       cfg,
		name:      "Authentication module",
		path:      "/auth",
		jwtSecret: jwtSecret,
	}
}

func (m *moduleAuthentication) Name() string { return m.name }
func (m *moduleAuthentication) Path() string { return m.path }

// Setup adiciona as rotas de login e retorna o router base para rotas protegidas (com middleware).
func (m *moduleAuthentication) Setup(r *mux.Router) *mux.Router {
	handlers := []modules.AppModuleHandler{
		{
			Handler: m.login,
			Path:    "/login",
			Label:   "Login do usuário",
			Methods: []string{http.MethodPost},
		},
	}

	authSub := r.PathPrefix(m.Path()).Subrouter()
	for _, h := range handlers {
		authSub.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}

	// Router base com middleware para rotas protegidas
	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware([]byte(m.jwtSecret)))
	return protected
}

type loginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type loginResponse struct {
	Usuario *entities.Usuario `json:"usuario"`
	Token   string            `json:"token"`
}

func (m *moduleAuthentication) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 1, Message: "corpo da requisição inválido"})
		return
	}
	if req.Email == "" || req.Senha == "" {
		httputil.WriteError(w, entities.ErrorStruct{Code: 4, Message: "email e senha são obrigatórios"})
		return
	}

	usr, token, err := m.useCase.Login(r.Context(), req.Email, req.Senha)
	if err != nil {
		httputil.WriteError(w, entities.ErrorStruct{Code: 3, Message: "credenciais inválidas"})
		return
	}

	resp := loginResponse{Usuario: usr, Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
