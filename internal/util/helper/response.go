package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Detail string `json:"detail"`
}

// OK writes 200 response with json data
func (helper *ServerHelper) OK(w http.ResponseWriter, model any) {
	helper.HTTPResponse(w, model, http.StatusOK)
}

// BadRequest writes 400 response with detail
func (helper *ServerHelper) BadRequest(w http.ResponseWriter, detail string) {
	helper.HTTPResponse(w, Error{Detail: detail}, http.StatusBadRequest)
}

// Unauthorized writes 401 response with detail
func (helper *ServerHelper) Unauthorized(w http.ResponseWriter, detail string) {
	helper.HTTPResponse(w, Error{Detail: detail}, http.StatusUnauthorized)
}

// Forbidden writes 403 response with detail
func (helper *ServerHelper) Forbidden(w http.ResponseWriter, detail string) {
	helper.HTTPResponse(w, Error{Detail: detail}, http.StatusForbidden)
}

// NotFound writes 404 response with detail
func (helper *ServerHelper) NotFound(w http.ResponseWriter, detail string) {
	helper.HTTPResponse(w, Error{Detail: detail}, http.StatusNotFound)
}

// Conflict writes 409 response with detail
func (helper *ServerHelper) Conflict(w http.ResponseWriter, detail string) {
	helper.HTTPResponse(w, Error{Detail: detail}, http.StatusConflict)
}

// ValidationError writes 422 response with detail
func (helper *ServerHelper) ValidationError(w http.ResponseWriter, err error) {
	helper.HTTPResponse(w, Error{Detail: err.Error()}, http.StatusUnprocessableEntity)
}

// InternalServerError writes 500 response and logs an error
func (helper *ServerHelper) InternalServerError(w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %v", err)
	helper.HTTPResponse(w, Error{Detail: "Internal ServerHelper Error"}, http.StatusInternalServerError)
}

// HTTPResponse writes response with model and status code
func (helper *ServerHelper) HTTPResponse(w http.ResponseWriter, model any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(model)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode)
	}

	w.Write(response)
}
