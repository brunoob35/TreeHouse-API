package repository

import (
	"database/sql"
	"github.com/brunoob35/TreeHouse-API/src/models"
)

// Users receives the DB connection and handles it
// it is representing the repository
type Students struct {
	db *sql.DB
}

// StudentsNewRepo receives the DB connection and injects the dependency of it into the repo thought the struct
func StudentsNewRepo(db *sql.DB) *Students {
	return &Students{db}
}

func (s Students) Create() (models.Classes, error) {
	return nil, nil
}

func (s Students) FetchByID(id int) (interface{}, error) {
	return nil, nil
}

func (s Students) FetchAll() (interface{}, error) {
	return nil, nil
}

func (s Students) Update(student models.Students) error {
	return nil
}
