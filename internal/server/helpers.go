package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"template-api/internal/util/validation"
)

// GetBody scans into struct and validates JSON body
func (s *Server) GetBody(w http.ResponseWriter, r *http.Request, model validation.BaseModel) bool {
	err := json.NewDecoder(r.Body).Decode(model)
	defer r.Body.Close()

	if err != nil {
		var synErr *json.SyntaxError
		var unmarshalErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &synErr):
			s.ValidationError(w, fmt.Errorf("request body contains badly-formed JSON (at position %d)", synErr.Offset))
			return false
		case errors.Is(err, io.EOF):
			s.ValidationError(w, errors.New("request body must not be empty"))
			return false
		case errors.As(err, &unmarshalErr):
			s.ValidationError(w, fmt.Errorf("request body contains an invalid value for the %q field (at position %d)", unmarshalErr.Field, unmarshalErr.Offset))
			return false
		default:
			s.InternalServerError(w, err)
			return false
		}
	}

	if err = model.Validate(); err != nil {
		s.ValidationError(w, err)
		return false
	}

	return true
}
