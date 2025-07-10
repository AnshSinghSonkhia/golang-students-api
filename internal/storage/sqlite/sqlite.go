package sqlite

import (
	"database/sql"

	"github.com/AnshSinghSonkhia/golang-students-api/internal/config"

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
