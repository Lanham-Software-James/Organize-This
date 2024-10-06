package mocks

import (
	context "context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// func createMockJWT() string {
// 	claims := jwt.MapClaims{
// 		"sub": "testuser",
// 		"iss": "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_example",
// 		"exp": time.Now().Add(time.Hour).Unix(),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 	tokenString, _ := token.SignedString([]byte("your-test-secret-key"))
// 	return tokenString
// }

func MockJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := jwt.MapClaims{
			"username": "testuser",
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
