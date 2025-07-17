package router

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
	"validation-api/internal/db"
	"validation-api/internal/server"
	"validation-api/internal/util/helper"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789_"
const maxLength = 63
const randLength = 10

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

type setup struct {
	router *Router
	db     *pgxpool.Pool
}

func urlWithDbName(url, db string) string {
	urlTrimmed := strings.TrimPrefix(url, "postgres://")
	split := strings.Split(urlTrimmed, "/")
	split[1] = db
	return "postgres://" + strings.Join(split, "/")
}

func Setup(t *testing.T) *setup {
	t.Helper()
	ctx := context.Background()

	// example: postgres://postgres:password@localhost:5432/postgres
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	// connect to main database through template1
	connStr := urlWithDbName(config.ConnString(), "template1")
	poolMain, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}

	// create test database
	testDbName := generateRandomPostgresDBName(t.Name())
	if _, err = poolMain.Exec(ctx, fmt.Sprintf(
		"CREATE DATABASE %v WITH TEMPLATE %v OWNER %v",
		testDbName,
		config.ConnConfig.Database,
		config.ConnConfig.User),
	); err != nil {
		t.Fatal(err)
	}

	// connect to test db and return connection
	connStr = urlWithDbName(config.ConnString(), testDbName)
	poolTest, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}

	server := &server.Server{
		Store: db.New(poolTest),
		ServerHelper: &helper.ServerHelper{
			MainLogger: helper.NewLogger("server"),
		},
	}

	// register cleanup to delete database after test
	t.Cleanup(func() {
		_, err = poolMain.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", testDbName))
		if err != nil {
			t.Fatal(err)
		}
	})

	return &setup{
		router: NewRouter(server, "/test", helper.NewLogger("test")),
		db:     poolTest,
	}
}

func TestSetup(t *testing.T) {
	setup := Setup(t)

	assert.NotNil(t, setup)
	assert.NotNil(t, setup.router)
	assert.NotNil(t, setup.db)
	assert.Contains(t, setup.db.Config().ConnConfig.Database, "testsetup")
}
