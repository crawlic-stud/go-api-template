package middleware

import (
	"context"
	"net/http"
	"strings"
	"template-api/internal/util/helper"
	"template-api/internal/util/services"
	"time"
)

func getTokenFromHeader(r *http.Request) string {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) == 2 {
		return authHeader[1]
	}
	return ""
}

func NewAuthMiddleware(service *services.AuthService, helper *helper.ServerHelper, skipper func(r *http.Request) bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if skipper(r) { // skips auth based on some condition
				next.ServeHTTP(w, r)
				return
			}

			jwtToken := getTokenFromHeader(r)
			if jwtToken == "" {
				helper.Unauthorized(w, "Malformed token")
				return

			} else {
				claims, err := service.VerifyToken(jwtToken)

				if err != nil {
					helper.Unauthorized(w, "Token is invalid")
					return

				} else if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
					helper.Unauthorized(w, "Token is expired")
					return

				} else {
					ctx := context.WithValue(r.Context(), service.AuthContextKey, claims)
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		})
	}
}
