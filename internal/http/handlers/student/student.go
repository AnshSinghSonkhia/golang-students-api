package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/AnshSinghSonkhia/golang-students-api/internal/types"
	"github.com/AnshSinghSonkhia/golang-students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// This file contains the handler for the root endpoint of the Golang Students API.

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a student")
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student) // decode the request body into a Student struct

		// if there is an error decoding the request body, check if it is an EOF error
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))

			return // return early to avoid further processing
		}

		// if there is an error decoding the request body, respond with a 400 Bad Request status code
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err)) // if there

			return
		}

		// Request Validataion

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors) // type assert the error to a ValidationErrors type

			// if there are validation errors, respond with a 400 Bad Request status code and the validation errors
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))

			return
		}

		// respond with a success message and a 201 Created status code
		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
