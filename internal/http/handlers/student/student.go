package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	// "github.com/AnshSinghSonkhia/golang-students-api/internal/http/handlers/student"
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

// GetByID(storage storage.Storage) returns a handler function that retrieves a student by ID.

func GetByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id") // get the ID from the URL path parameters

		slog.Info("Retrieving student by ID: ", slog.String("id", id)) // log the ID being retrieved

		intTd, err := strconv.ParseInt(id, 10, 64) // convert the ID from string to int64
		if err != nil {

			slog.Error("Error converting ID to int64", slog.String("id", id), slog.Any("error", err)) // log the error if there is an issue converting the ID

			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err)) // if there is an error converting the ID, respond with a 400 Bad Request status code
			return                                                                   // return early to avoid further processing
		}

		student, err := storage.GetStudentByID(intTd) // call the GetStudentByID method on the storage interface to retrieve the student by ID

		if err != nil {

			slog.Error("Error retrieving student by ID", slog.String("id", id), slog.Any("error", err)) // log the error if there is an issue retrieving the student

			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err)) // if there is any other error, respond with a 500 Internal Server Error status code

			return // return early to avoid further processing
		}

		response.WriteJSON(w, http.StatusOK, student) // if the student is found, respond with a 200 OK status code and the student data
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Retrieving list of students") // log the action of retrieving the list of students

		students, err := storage.GetStudents() // call the GetStudents method on the storage interface to retrieve the list of students
		if err != nil {
			slog.Error("Error retrieving list of students", slog.Any("error", err)) // log the error if there is an issue retrieving the list

			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err)) // if there is an error, respond with a 500 Internal Server Error status code

			return // return early to avoid further processing
		}

		response.WriteJSON(w, http.StatusOK, students) // if the list is retrieved successfully, respond with a 200 OK status code and the list of students
	}
}
