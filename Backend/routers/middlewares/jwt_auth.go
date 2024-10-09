package middlewares

import (
	"context"
	"net/http"
	"organize-this/controllers"
	"organize-this/helpers"
	"strings"
	"time"
)

func JWTAuth(handler controllers.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the JWT from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				helpers.BadRequest(w, "Missing Authorization Header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				helpers.BadRequest(w, "Invalid Authorization header format")
				return
			}

			token, err := handler.TokenHelper.VerifyToken(tokenString, true)
			if err != nil && strings.Contains(err.Error(), "Failed to get JWKS") {
				helpers.BadRequest(w, err)
				return
			} else if err != nil {
				helpers.UnaunthorizedRequest(w, err)
				return
			}

			claims, err := handler.TokenHelper.ExtractClaims(token)
			if err != nil {
				helpers.UnaunthorizedRequest(w, err)
				return
			}

			// Verify token_use claim for access tokens
			if claims["token_use"] != "access" {
				http.Error(w, "Invalid token use", http.StatusUnauthorized)
				return
			}

			// Check expiration
			exp, ok := claims["exp"].(float64)
			if !ok {
				http.Error(w, "Invalid exp claim", http.StatusUnauthorized)
				return
			}

			isRefreshRequest := (r.Method == "PUT" && strings.Contains(r.URL.Path, "token"))

			if time.Now().Unix() > int64(exp) && !isRefreshRequest {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}

			// Add claims to request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_claims", claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
