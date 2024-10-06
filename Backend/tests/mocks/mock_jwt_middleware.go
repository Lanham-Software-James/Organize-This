package mocks

import (
	context "context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func MockJWTMiddleware(userName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := jwt.MapClaims{
				"username": userName,
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_claims", claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
