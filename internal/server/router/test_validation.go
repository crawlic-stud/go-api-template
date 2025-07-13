package router

import (
	"net/http"
	"validation-api/internal/models"
)

func (api *Router) TestValidation(w http.ResponseWriter, r *http.Request) {
	var user models.User
	valid := api.GetBody(w, r, &user)
	if !valid {
		return
	}

	api.OK(w, user)
}
