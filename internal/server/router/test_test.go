package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createTestDatabase(t *testing.T) {
	pool := createTestDatabase(t)

	assert.NotNil(t, pool)
	assert.Contains(t, pool.Config().ConnConfig.Database, "test_createtestdatabase")
}
