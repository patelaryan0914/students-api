package storage

import "github.com/patelaryan0914/students-api/internal/types"
type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64)(types.Student,error)
}

