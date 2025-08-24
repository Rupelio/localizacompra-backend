package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"

var jwtSecret = []byte("sua-chave-super-secreta")

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Não autorizado: token não encontrado", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := (*claims)["sub"].(float64)
		if !ok {
			http.Error(w, "ID do usuário não encontrado no token", http.StatusUnauthorized)
			return
		}

		userRole, ok := (*claims)["role"].(string)
		if !ok {
			http.Error(w, "Cargo do utilizador não encontrado no token", http.StatusUnauthorized)
			return
		}

		userID := int64(userIDFloat)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(UserRoleKey).(string)

		if !ok || (role != "admin" && role != "store_admin") {
			http.Error(w, "Acesso negado: rota apenas para administradores", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SuperAdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(UserRoleKey).(string)

		if !ok || (role != "super_admin") {
			http.Error(w, "Acesso negado: rota apenas para super administradores", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
