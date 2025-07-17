package router

import (
	"fmt"
	"net/http"
	"validation-api/internal/db"
	"validation-api/internal/models"
)

func (api *Router) Register(w http.ResponseWriter, r *http.Request) {
	var user models.LoginUser
	if ok := api.GetBody(w, r, &user); !ok {
		return
	}

	exists, err := api.Store.UsernameExists(r.Context(), user.Username)
	if err != nil {
		api.InternalServerError(w, err)
		return
	}

	if exists {
		api.Conflict(w, fmt.Sprintf("Username '%s' already exists", user.Username))
		return
	}

	hashedPassword, err := api.Auth.HashPassword(user.Password)
	if err != nil {
		api.InternalServerError(w, err)
		return
	}

	if err = api.Store.CreateUser(r.Context(), db.CreateUserParams{
		Username:       user.Username,
		HashedPassword: hashedPassword,
	}); err != nil {
		api.InternalServerError(w, err)
		return
	}

	api.OK(w, user)
}
