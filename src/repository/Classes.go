package repository

import (
	"database/sql"
	"github.com/brunoob35/TreeHouse-API/src/models"
)

// Classes receives the DB connection and handles it
// it is representing the repository
type Classes struct {
	db *sql.DB
}

// StudentsNewRepo receives the DB connection and injects the dependency of it into the repo thought the struct
func ClassesNewRepo(db *sql.DB) *Classes {
	return &Classes{db}
}

// Create Creates a new "Turma" in the DB
func (c Classes) Create(class models.Classes) (uint64, error) {
	query := `INSERT INTO turmas
				(nome_turma)
				VALUES (?)`

	statement, err := c.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(class.Name)
	if err != nil {
		return 0, err
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertedID), nil
}

func (c Classes) FetchByID(id int) (interface{}, error) {
	return nil, nil
}

func (c Classes) FetchAll() (interface{}, error) {
	return nil, nil
}

func (c Classes) Update(interface{}) error {
	return nil
}
