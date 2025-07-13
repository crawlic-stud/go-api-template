package router

import (
	"fmt"
	"net/http"
	"validation-api/internal/models"
)

func (api *Router) Login(w http.ResponseWriter, r *http.Request) {
	var user models.LoginUser
	if ok := api.GetBody(w, r, &user); !ok {
		return
	}

	userDb, err := api.Store.GetUserByUsername(r.Context(), user.Username)
	if err != nil {
		api.NotFound(w, fmt.Sprintf("User '%s' not found", user.Username))
		return
	}

	if !api.Auth.CheckPasswordHash(user.Password, userDb.HashedPassword) {
		api.Unauthorized(w, "Password is incorrect")
		return
	}

	token, err := api.Auth.GenerateToken(userDb.ID.String())
	if err != nil {
		api.InternalServerError(w, err)
		return
	}

	api.OK(w, map[string]any{
		"token": token,
	})
}
