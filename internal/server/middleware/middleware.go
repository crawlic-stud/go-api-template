package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Stack - stacks middleware functions one after another
func Stack(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, m := range middleware {
			next = m(next)
		}
		return next
	}
}
