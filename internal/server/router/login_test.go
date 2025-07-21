package router

import (
	"context"
	"template-api/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	t.Run("it works", func(t *testing.T) {
		t.Parallel()

		setup, recorder := Setup(t)
		req := setup.POST("/login", `{"username": "testuser", "password": "password"}`)

		hashedPassword, err := setup.router.Auth.HashPassword("password")
		assert.NoError(t, err)
		setup.router.Store.CreateUser(context.Background(), db.CreateUserParams{
			Username:       "testuser",
			HashedPassword: hashedPassword,
		})

		setup.router.Login(recorder, req)

		assert.Equal(t, 200, recorder.Result().StatusCode)

		result := fromReader[map[string]string](t, recorder.Body)
		token, ok := result["token"]
		assert.True(t, ok)
		assert.NotEmpty(t, token)
	})
}
