# Padrão de Criação de Rotas – Prompt para Aplicar em Outro Projeto

Use este documento como prompt no Cursor para replicar o sistema de rotas em outro projeto. Substitua `PROJECT_NAME` pelo nome do seu módulo Go (ex: `lince`, `meuapp`).

---

## FORMATO DAS URLs

```
/api                          ← prefixo global da API (opcional)
  └── /auth                   ← módulo de autenticação (rotas públicas)
        └── /login            ← ação do endpoint
  └── /categorias             ← prefixo do módulo
        └── /list             ← ação: listar
        └── /create           ← ação: criar
        └── /get/{id}         ← ação: buscar por ID
  └── /subcategorias
        └── /list
        └── /create
        └── /get/{id}
```

**URLs finais (com /api):**
- `POST /api/auth/login` (público, sem token)
- `GET /api/categorias/list` (protegido)
- `POST /api/categorias/create` (protegido)
- `GET /api/categorias/get/{id}` (protegido)

---

## FLUXO DE MONTAGEM DAS URLs

| Ordem | Onde | O quê | Resultado |
|-------|------|-------|-----------|
| 1 | setup.go | `r.PathPrefix("/api")` | Todas as rotas começam com /api |
| 2 | authenticationModule.Setup | Adiciona `/auth` + `/login` | `/api/auth/login` |
| 3 | authenticationModule.Setup | Retorna `protected` (com middleware) | Router para rotas protegidas |
| 4 | setup.go | `routerBase.PathPrefix(am.Path())` → `/categorias` | Prefixo do módulo |
| 5 | module.go (categoria) | `h.Path` → `/list`, `/create`, `/get/{id}` | Ação do endpoint |

---

## ARQUIVOS NECESSÁRIOS

### 1. `main.go`

Não precisa alterar. Apenas cria o router e chama SetupModules.

```go
package main

import (
	"log"
	"net/http"

	"PROJECT_NAME/entities"
	"PROJECT_NAME/setup"

	"github.com/gorilla/mux"
)

func main() {
	cfg := entities.Config{}
	r := mux.NewRouter()

	setup.SetupModules(r, cfg)

	srv := &http.Server{Addr: ":8080", Handler: r}
	log.Printf("HTTP server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
```

---

### 2. `modules/modules.go`

Interface base. **Arquivo já deve existir.** Contém `AppModule` e `AppModuleHandler`.

```go
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
	Path    string   // Ex: "/list", "/create", "/get/{id}"
	Label   string
	Handler http.HandlerFunc
	Methods []string // Ex: []string{http.MethodGet}
}
```

---

### 3. `setup/setup.go`

**ALTERAR:** Adicione o prefixo `/api` e o registro dos módulos conforme abaixo.

```go
package setup

import (
	"log"

	"PROJECT_NAME/datastore"
	"PROJECT_NAME/entities"
	"PROJECT_NAME/migrations"
	"PROJECT_NAME/modules"
	"PROJECT_NAME/modules/category"
	"PROJECT_NAME/modules/subcategory"
	"PROJECT_NAME/modules/user"

	"github.com/gorilla/mux"
)

func SetupModules(r *mux.Router, cfg entities.Config) {
	log.Println("Setup modules")

	// ... (database, migrations, repositories, use cases) ...

	// ### ENDPOINT MODULES ###
	categoryModule := category.NewCategoryModule(categoryUseCase)
	subcategoryModule := subcategory.NewSubcategoryModule(subcategoryUseCase)
	authenticationModule := user.NewAuthenticationModule(cfg, userUseCase)

	appModules := []modules.AppModule{
		categoryModule,
		subcategoryModule,
	}

	// PASSO 1: Prefixo global /api (todas as rotas ficam sob /api)
	apiRouter := r.PathPrefix("/api").Subrouter()

	// PASSO 2: Auth setup - adiciona /auth/login e retorna router protegido (com middleware)
	routerBase := authenticationModule.Setup(apiRouter)

	// PASSO 3: Cada módulo protegido - am.Path() = /categorias, /subcategorias
	for _, am := range appModules {
		moduleSubRouter := routerBase.PathPrefix(am.Path()).Subrouter()
		_ = am.Setup(moduleSubRouter)
	}
}
```

**Importante:** Se não quiser o prefixo `/api`, use `apiRouter := r` e passe `r` diretamente para `authenticationModule.Setup(r)`.

---

### 4. `modules/user/module.go` – Módulo de Autenticação (especial)

Este módulo é **diferente**: adiciona rotas públicas (login) e retorna o router base com middleware para as rotas protegidas.

```go
package user

import (
	"encoding/json"
	"net/http"

	"PROJECT_NAME/domain"
	"PROJECT_NAME/entities"
	"PROJECT_NAME/httputil"
	"PROJECT_NAME/middleware"
	"PROJECT_NAME/modules"

	"github.com/gorilla/mux"
)

type moduleAuthentication struct {
	useCase   domain.UserUseCase
	cfg       entities.Config
	name      string
	path      string   // "/auth"
	jwtSecret string
}

func NewAuthenticationModule(cfg entities.Config, useCase domain.UserUseCase) modules.AppModule {
	jwtSecret := cfg.JWTSecret
	if jwtSecret == "" {
		jwtSecret = "default-secret-key"
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

// Setup: 1) Adiciona rotas públicas em /auth; 2) Retorna router com middleware para rotas protegidas
func (m *moduleAuthentication) Setup(r *mux.Router) *mux.Router {
	handlers := []modules.AppModuleHandler{
		{
			Handler: m.login,
			Path:    "/login",       // URL final: /api/auth/login
			Label:   "Login",
			Methods: []string{http.MethodPost},
		},
	}

	// Rotas públicas: r já tem prefixo /api, então authSub = /api/auth
	authSub := r.PathPrefix(m.Path()).Subrouter()
	for _, h := range handlers {
		authSub.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}

	// Router base para módulos protegidos (com middleware de token)
	protected := r.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware([]byte(m.jwtSecret)))
	return protected
}

// ... (implementação do handler m.login)
```

**Resumo:** O `r` recebido já tem o prefixo `/api`. O `authSub` adiciona `/auth`, então `/login` vira `/api/auth/login`. O `protected` é o mesmo router com middleware aplicado, usado para os outros módulos.

---

### 5. `modules/category/module.go` – Módulo protegido (exemplo)

O `r` recebido no `Setup` já é o subrouter com prefixo `/api` + `/categorias`. Os handlers usam `h.Path` como sufixo.

```go
package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"PROJECT_NAME/domain"
	"PROJECT_NAME/entities"
	"PROJECT_NAME/httputil"
	"PROJECT_NAME/modules"

	"github.com/gorilla/mux"
)

type moduleCategory struct {
	useCase domain.CategoryUseCase
	name    string
	path    string   // "/categorias"
}

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
			Path:    "/list",        // URL final: /api/categorias/list
			Label:   "Lista todas as categorias",
			Methods: []string{http.MethodGet},
		},
		{
			Handler: m.create,
			Path:    "/create",      // URL final: /api/categorias/create
			Label:   "Cadastra nova categoria",
			Methods: []string{http.MethodPost},
		},
		{
			Handler: m.getByID,
			Path:    "/get/{id}",    // URL final: /api/categorias/get/123
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
	// ... implementação
}

func (m *moduleCategory) create(w http.ResponseWriter, r *http.Request) {
	// ... implementação
}

func (m *moduleCategory) getByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	// ... implementação
}
```

**Regra:** O `r` em `Setup(r)` já está em `/api/categorias`. Cada `h.Path` é o sufixo. Não use `PathPrefix` novamente dentro do módulo.

---

### 6. `modules/subcategory/module.go` – Outro módulo protegido

Mesmo padrão. `path: "/subcategorias"`. Handlers com `Path: "/list"`, `"/create"`, `"/get/{id}"`.

```go
package subcategory

import (
	"encoding/json"
	"net/http"
	"strconv"

	"PROJECT_NAME/domain"
	"PROJECT_NAME/entities"
	"PROJECT_NAME/httputil"
	"PROJECT_NAME/modules"

	"github.com/gorilla/mux"
)

type moduleSubcategory struct {
	useCase domain.SubcategoryUseCase
	name    string
	path    string
}

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
		{Handler: m.list, Path: "/list", Label: "Lista subcategorias", Methods: []string{http.MethodGet}},
		{Handler: m.create, Path: "/create", Label: "Cadastra subcategoria", Methods: []string{http.MethodPost}},
		{Handler: m.getByID, Path: "/get/{id}", Label: "Busca subcategoria por ID", Methods: []string{http.MethodGet}},
	}

	for _, h := range handlers {
		r.HandleFunc(h.Path, h.Handler).Methods(h.Methods...)
	}
	return r
}

// ... handlers (list, create, getByID)
```

---

## TABELA DE REFERÊNCIA RÁPIDA

| Arquivo | Campo/Valor | Gera na URL |
|---------|-------------|-------------|
| setup.go | `r.PathPrefix("/api")` | `/api` |
| user/module.go | `path: "/auth"` | `/api/auth` |
| user/module.go | `Path: "/login"` | `/api/auth/login` |
| category/module.go | `path: "/categorias"` | `/api/categorias` |
| category/module.go | `Path: "/list"` | `/api/categorias/list` |
| category/module.go | `Path: "/create"` | `/api/categorias/create` |
| category/module.go | `Path: "/get/{id}"` | `/api/categorias/get/123` |

---

## CHECKLIST PARA APLICAR NO OUTRO PROJETO

1. **setup/setup.go**
   - [ ] Adicionar `apiRouter := r.PathPrefix("/api").Subrouter()`
   - [ ] Passar `apiRouter` para `authenticationModule.Setup(apiRouter)` em vez de `r`
   - [ ] Manter o loop: `routerBase.PathPrefix(am.Path()).Subrouter()` e `am.Setup(moduleSubRouter)`

2. **modules/user/module.go** (autenticação)
   - [ ] Manter `path: "/auth"`
   - [ ] Handlers com `Path: "/login"` (ou outra ação)
   - [ ] `Setup` retornando `protected` (router com middleware)

3. **modules/{dominio}/module.go** (cada módulo protegido)
   - [ ] `path: "/nomedomodulo"` (ex: `/categorias`)
   - [ ] Handlers com `Path: "/list"`, `"/create"`, `"/get/{id}"` (ou ações desejadas)
   - [ ] Dentro de `Setup`, usar `r.HandleFunc(h.Path, h.Handler)` – o `r` já vem com o prefixo correto

---

## VARIAÇÕES

- **Sem /api:** Em setup, use `routerBase := authenticationModule.Setup(r)` em vez de criar `apiRouter`.
- **Path vazio para list/create:** Alguns preferem `Path: ""` para GET e POST no mesmo recurso (ex: `GET /categorias` e `POST /categorias`). Nesse caso, o método HTTP diferencia. Para URLs como `/list` e `/create`, use sempre o Path explícito.
- **Novo módulo:** Crie a pasta `modules/novomodulo/`, adicione `module.go` implementando `AppModule`, e inclua em `appModules` no setup.
