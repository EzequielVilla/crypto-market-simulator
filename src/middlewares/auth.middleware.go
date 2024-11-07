package middlewares

import (
	"context"
	"crypto-market-simulator/internal/lib/my_jwt"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenWithBearer := r.Header.Get("Authorization")
		words := strings.Split(tokenWithBearer, " ")
		token := words[1]
		if token == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		userId, walletId, err := my_jwt.ParseToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
		ctx := context.WithValue(r.Context(), "userId", userId)
		ctx = context.WithValue(ctx, "walletId", walletId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
