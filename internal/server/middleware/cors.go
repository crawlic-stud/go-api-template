package middleware

import (
	"net/http"
)

type CORSOpt func(http.ResponseWriter)

var defaultCORSSettings = [4]CORSOpt{
	WithAllowOrigin("*"),
	WithAllowMethods("GET, POST, PUT, DELETE, OPTIONS, PATCH"),
	WithAllowHeaders("Accept, Authorization, Content-Type, X-CSRF-Token"),
	WithAllowCredentials("false"),
}

// WithAllowOrigin sets Allow-Origin header.
// Default is *.
func WithAllowOrigin(origins string) CORSOpt {
	return func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Origin", origins)
	}
}

// WithAllowMethods sets Allow-Methods header.
// Default is GET, POST, PUT, DELETE, OPTIONS, PATCH.
func WithAllowMethods(methods string) CORSOpt {
	return func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Methods", methods)
	}
}

// WithAllowHeaders sets Allow-Headers header.
// Default is Accept, Authorization, Content-Type, X-CSRF-Token.
func WithAllowHeaders(headers string) CORSOpt {
	return func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Headers", headers)
	}
}

// WithAllowCredentials sets Allow-Credentials header.
// Default is false.
func WithAllowCredentials(credentials string) CORSOpt {
	return func(w http.ResponseWriter) {
		w.Header().Set("Access-Control-Allow-Credentials", credentials)
	}
}

// NewCORSMiddleware constructs CORS middleware with user defined options
func NewCORSMiddleware(opts ...CORSOpt) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, defaultOpt := range defaultCORSSettings {
				defaultOpt(w)
			}

			for _, userOpt := range opts {
				userOpt(w)
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}
