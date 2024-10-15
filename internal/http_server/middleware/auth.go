package middleware

import (
	"UrlShorterService/internal/services/jwt"
	"context"
	"net/http"
	"strings"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("Authorization"), " ")[1]
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId, err := jwt.DecodeJWT(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
