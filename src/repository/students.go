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

// Create Creates a new student in DB
func (s Students) Create(student models.Students) (uint64, error) {
	query := `INSERT INTO treehousedb.alunos
				(nome)
				VALUES (?)`

	statement, err := s.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(student.Name)
	if err != nil {
		return 0, err
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertedID), nil
}

// FetchByID Fetches a student from the DB
func (s Students) FetchByID(studentID uint64) (models.Students, error) {
	query := `SELECT * FROM treehousedb.alunos WHERE id_aluno = ?`

	line, err := s.db.Query(query, studentID)
	if err != nil {
		return models.Students{}, err
	}

	defer line.Close()

	var student models.Students

	if line.Next() {
		if err = line.Scan(
			&student.ID,
			&student.Name,
			&student.Active,
		); err != nil {
			return models.Students{}, err
		}
	}

	return student, nil
}

// FetchAll Fetches all students either active or not
func (s Students) FetchAll() ([]models.Students, error) {
	query := `SELECT * FROM treehousedb.alunos`

	lines, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer lines.Close()

	var students []models.Students

	for lines.Next() {
		var student models.Students

		if err = lines.Scan(
			&student.ID,
			&student.Name,
			&student.Active,
		); err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

// FetchAllActive fetches all active students ignoring inactive ones
func (s Students) FetchAllActive() ([]models.Students, error) {
	query := `SELECT * FROM treehousedb.alunos WHERE ativo = 1`

	lines, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer lines.Close()

	var students []models.Students

	for lines.Next() {
		var student models.Students

		if err = lines.Scan(
			&student.ID,
			&student.Name,
			&student.Active,
		); err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s Students) Update(student models.Students) error {
	return nil
}
