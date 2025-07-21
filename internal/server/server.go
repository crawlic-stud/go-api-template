package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"template-api/internal/db"
	"template-api/internal/util/helper"
	"template-api/internal/util/services"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	*helper.ServerHelper

	Store *db.Queries

	Auth *services.AuthService
}

func NewHTTPServer(handler http.Handler) *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
