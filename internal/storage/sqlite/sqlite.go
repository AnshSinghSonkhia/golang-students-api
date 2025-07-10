package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/AnshSinghSonkhia/golang-students-api/internal/config"
	"github.com/AnshSinghSonkhia/golang-students-api/internal/types"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
	// use of _ means we are importing the package solely for its side effects (registering the driver) - i.e., indirect usess
)

type Sqlite struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// Create the students table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		age INTEGER NOT NULL
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	return &Sqlite{
		DB: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	// Prepare the SQL statement to insert a new student
	stmt, err := s.DB.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)") // ? are placeholders for the values to be inserted

	if err != nil {
		return 0, err // Return an error if the statement preparation fails
	}
	defer stmt.Close() // Ensure the statement is closed after use

	// Execute the statement with the provided values
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err // Return an error if the execution fails
	}

	// Get the last inserted ID
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err // Return an error if retrieving the last inserted ID fails
	}

	// Return the last inserted ID and no error
	return lastId, nil
}

func (s *Sqlite) GetStudentByID(id int64) (types.Student, error) {
	stmt, err := s.DB.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1") // Prepare the SQL statement to select a student by ID
	if err != nil {
		return types.Student{}, err // Return an empty Student struct and an error if preparation fails
	}

	defer stmt.Close() // Ensure the statement is closed after use

	var student types.Student // Create a Student struct to hold the retrieved data

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age) // Execute the query and scan the result into the Student struct

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with ID %d not found", id) // Return an empty Student struct and a not found error if no rows are returned
		}

		return types.Student{}, fmt.Errorf("query error: %w", err) // Return an empty Student struct and an error if the query fails
	}

	// Return the retrieved Student struct and no error
	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {

	// Prepare the SQL statement to select all students
	stmt, err := s.DB.Prepare("SELECT id, name, email, age FROM students")

	if err != nil {
		return nil, err // Return nil and an error if the statement preparation fails
	}

	defer stmt.Close() // Ensure the statement is closed after use

	rows, err := stmt.Query() // Execute the query to get all students
	if err != nil {
		return nil, err // Return nil and an error if the query execution fails
	}

	defer rows.Close() // Ensure the rows are closed after use

	var students []types.Student // Create a slice to hold the retrieved students

	for rows.Next() { // Iterate over the rows returned by the query
		var student types.Student // Create a Student struct to hold the data for each row

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age) // Scan the row data into the Student struct
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err) // Return nil and an error if scanning fails
		}

		students = append(students, student) // Append the Student struct to the slice of students
	}

	return students, nil // Return the slice of students and no error
}
