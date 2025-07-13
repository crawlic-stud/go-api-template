package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"validation-api/internal/db"
	"validation-api/internal/server"
	"validation-api/internal/server/middleware"
	"validation-api/internal/server/router"
	"validation-api/internal/util/helper"
	"validation-api/internal/util/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupPool() *pgxpool.Pool {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), config.ConnString())
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database!")

	return pool
}

// Routes with api prefix
func registerAPIRoutes(router *router.Router) {
	router.Route("POST /login", router.Login)
	router.Route("POST /register", router.Register)
}

func setup() (*http.Server, func()) {
	pool := setupPool()

	serverHelper := &helper.ServerHelper{
		MainLogger: helper.NewLogger("server"),
	}

	s := &server.Server{
		Store: db.New(pool),
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

	cleanup := func() {
		pool.Close()
	}

	return server.NewHTTPServer(handler), cleanup
}

func main() {
	server, cleanup := setup()
	defer cleanup()

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
