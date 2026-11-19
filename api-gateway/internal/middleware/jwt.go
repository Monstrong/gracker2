package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/monstrong/gracker2/pkg/jwtparser"
)

const userClaimsKey string = "user_claims"

func AuthMiddleware(tokenManager jwtparser.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			parts := strings.Fields(authHeader)

			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
				http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
    			return
			}
			
			tokenString := parts[1]
			claims, err := tokenManager.Parse(tokenString)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, userClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(ctx context.Context) (*jwtparser.Claims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*jwtparser.Claims)
	return claims, ok
}