package middleware

import (
	"context"
	"net/http"
	"strings"

	"lince/entities"
	"lince/httputil"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware verifica o token JWT no header Authorization e injeta o user_id no contexto.
func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				httputil.WriteError(w, entities.ErrorStruct{Code: 3, Message: "token ausente"})
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				httputil.WriteError(w, entities.ErrorStruct{Code: 3, Message: "token inválido"})
				return
			}

			token, err := jwt.ParseWithClaims(parts[1], &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				httputil.WriteError(w, entities.ErrorStruct{Code: 3, Message: "token inválido ou expirado"})
				return
			}

			claims, ok := token.Claims.(*jwt.RegisteredClaims)
			if !ok || claims.Subject == "" {
				httputil.WriteError(w, entities.ErrorStruct{Code: 3, Message: "token inválido"})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
