package router

import (
	"net/http"
	"validation-api/internal/db"
	"validation-api/internal/models"
)

func (api *Router) Register(w http.ResponseWriter, r *http.Request) {
	var user models.LoginUser
	if ok := api.GetBody(w, r, &user); !ok {
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
