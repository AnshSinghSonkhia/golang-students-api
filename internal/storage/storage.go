package storage

import "github.com/AnshSinghSonkhia/golang-students-api/internal/types"

type Storage interface {
	// CreateStudent creates a new student in the storage.
	CreateStudent(name string, email string, age int) (int64, error)

	// GetStudentByID retrieves a student by ID from the storage.
	GetStudentByID(id int64) (types.Student, error)

	// GetStudents retrieves all students from the storage.
	GetStudents() ([]types.Student, error)

	// UpdateStudent updates an existing student in the storage.
	UpdateStudent(id int64, name string, email string, age int) error

	// DeleteStudent deletes a student by ID from the storage.
	DeleteStudent(id int64) error
}
