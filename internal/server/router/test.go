package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"validation-api/internal/db"
	"validation-api/internal/server"
	"validation-api/internal/util/helper"
	"validation-api/internal/util/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789_"
const maxLength = 63
const randLength = 10

type setup struct {
	router *Router
}

func generateRandomPostgresDBName(testName string) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	sb.Grow(maxLength)

	if len(testName) > maxLength-randLength {
		testName = testName[:maxLength-randLength]
	}

	sb.WriteString(testName)
	for i := 1; i < randLength; i++ {
		sb.WriteByte(charset[seededRand.Intn(len(charset))])
	}

	return strings.ToLower(sb.String())
}

func urlWithDbName(t *testing.T, url, db string) string {
	if url == "" {
		t.Fatal("empty database url")
	}

	urlTrimmed := strings.TrimPrefix(url, "postgres://")
	split := strings.Split(urlTrimmed, "/")
	if len(split) != 2 {
		t.Fatalf("invalid database url: %v", url)
	}
	split[1] = db
	return "postgres://" + strings.Join(split, "/")
}

func createTestDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	// example: postgres://postgres:password@localhost:5432/postgres
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("can not get database config: %v", err)
	}

	// connect to main database through template1
	connStr := urlWithDbName(t, config.ConnString(), "template1")
	poolMain, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("can not connect to main database: %v", err)
	}

	// create test database
	testName := strings.ReplaceAll(t.Name(), "/", "_")
	testDbName := generateRandomPostgresDBName(testName)
	if _, err = poolMain.Exec(ctx, fmt.Sprintf(
		"CREATE DATABASE %v WITH TEMPLATE %v OWNER %v",
		testDbName,
		config.ConnConfig.Database,
		config.ConnConfig.User),
	); err != nil {
		t.Fatalf("can not create test database: %v", err)
	}

	// connect to test db and return connection
	connStr = urlWithDbName(t, config.ConnString(), testDbName)
	poolTest, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("can not connect to test database: %v", err)
	}

	// register cleanup to close connections and delete database after test
	t.Cleanup(func() {
		poolTest.Close()
		_, err = poolMain.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", testDbName))
		if err != nil {
			t.Fatalf("can not delete test database: %v", err)
		}
		poolMain.Close()
	})

	return poolTest
}

func Setup(t *testing.T) (*setup, *httptest.ResponseRecorder) {
	t.Helper()

	poolTest := createTestDatabase(t)
	log.Printf("Test is using test database: %s", poolTest.Config().ConnConfig.Database)
	server := &server.Server{
		Store: db.New(poolTest),
		ServerHelper: &helper.ServerHelper{
			MainLogger: helper.NewLogger("server"),
		},
		Auth: services.NewAuthService(
			os.Getenv("APP_SECRET"),
			86_400_000, // 1000 days
		),
	}

	return &setup{
		router: NewRouter(server, "/test", helper.NewLogger("test")),
	}, httptest.NewRecorder()
}

func (s *setup) POST(url, body string) *http.Request {
	req := httptest.NewRequest(
		http.MethodPost,
		s.router.prefix+url,
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func scanReader[T any](t *testing.T, reader io.Reader) T {
	var model T
	err := json.NewDecoder(reader).Decode(&model)
	if err != nil {
		t.Fatal(err)
	}
	return model
}
