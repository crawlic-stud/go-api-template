package router

import (
	"context"
	"template-api/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("it works", func(t *testing.T) {
		t.Parallel()

		setup, recorder := Setup(t)
		req := setup.POST("/register", `{"username": "testuser", "password": "password"}`)

		setup.router.Register(recorder, req)

		assert.Equal(t, 200, recorder.Result().StatusCode)
		user, err := setup.router.Store.GetUserByUsername(context.Background(), "username")
		assert.NoError(t, err)
		assert.Equal(t, "username", user.Username)
	})

	t.Run("username already exists", func(t *testing.T) {
		t.Parallel()

		setup, recorder := Setup(t)
		req := setup.POST("/register", `{"username": "testuser", "password": "password"}`)

		err := setup.router.Store.CreateUser(context.Background(), db.CreateUserParams{
			Username:       "testuser",
			HashedPassword: "password",
		})
		assert.NoError(t, err)

		setup.router.Register(recorder, req)

		assert.Equal(t, 409, recorder.Result().StatusCode)
		assert.Contains(t, recorder.Body.String(), "Username 'testuser' already exists")
	})
}
