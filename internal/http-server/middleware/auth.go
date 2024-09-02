package middleware

import (
	"UrlShorterService/internal/services/jwt"
	"fmt"
	"net/http"
	"strings"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("Authorization"), " ")[1]
		fmt.Println(token)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		status, err := jwt.IsCorrectJwtToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		if !status {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(err.Error()))
		}
		next.ServeHTTP(w, r)
	})
}
