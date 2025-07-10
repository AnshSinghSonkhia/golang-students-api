package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/AnshSinghSonkhia/golang-students-api/internal/storage"
	"github.com/AnshSinghSonkhia/golang-students-api/internal/types"
	"github.com/AnshSinghSonkhia/golang-students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// This file contains the handler for the root endpoint of the Golang Students API.

// New(storage storage.Storage) means -> injecting the storage dependency into the handler function. This allows the handler to access the storage layer for database operations.

func New(storage storage.Storage) http.HandlerFunc {
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

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age) // call the CreateStudent method on the storage interface to create a new student

		slog.Info("Student created successfully", slog.Int64("id", lastId), slog.String("name", student.Name), slog.String("email", student.Email), slog.Int("age", student.Age))

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err)) // if there is an error creating the student, respond with a 500 Internal Server Error status code

			return // return early to avoid further processing
		}

		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastId}) // return the last inserted ID in the response
	}
}
