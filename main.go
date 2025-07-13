package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"validation-api/internal/server"
	"validation-api/internal/server/middleware"
	"validation-api/internal/server/router"
	"validation-api/internal/util/helper"
	"validation-api/internal/util/services"
)

// Routes with api prefix
func registerAPIRoutes(router *router.Router) {
	router.Route("POST /", router.TestValidation)

	router.Route("POST /login", router.Login)
	router.Route("POST /register", router.Register)
}

func setup() *http.Server {
	serverHelper := &helper.ServerHelper{
		MainLogger: helper.NewLogger("server"),
	}

	s := &server.Server{
		Auth: services.NewAuthService(
			os.Getenv("APP_SECRET"),
			86_400_000, // 1000 days
		),
		ServerHelper: serverHelper,
	}

	apiRouter := router.NewRouter(s, "/api", helper.NewLogger("api"))
	registerAPIRoutes(apiRouter)

	skipAuthRoutes := map[string]bool{
		"/api/login":    true,
		"/api/register": true,
	}

	middlewareLayer := middleware.Stack(
		middleware.NewCORSMiddleware(),
		middleware.NewAuthMiddleware(s.Auth, serverHelper, func(r *http.Request) bool { return skipAuthRoutes[r.URL.Path] }),
		middleware.NewLoggingMiddleware(serverHelper),
	)

	handler := middlewareLayer(apiRouter)

	return server.NewHTTPServer(handler)
}

func main() {
	server := setup()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go helper.GracefulShutdown(server, done)

	log.Printf("Server started on http://localhost:%s", os.Getenv("PORT"))
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
