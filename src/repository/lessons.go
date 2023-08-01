package repository

import (
	"database/sql"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"log"
)

// Classes receives the DB connection and handles it
// it is representing the repository
type Lessons struct {
	db *sql.DB
}

// StudentsNewRepo receives the DB connection and injects the dependency of it into the repo thought the struct
func LessonsNewRepo(db *sql.DB) *Lessons {
	return &Lessons{db}
}

// CreateLessons Creates a new "Aula" in the DB
func (l Lessons) Create(aula models.Lessons) (uint64, error) {
	query := `INSERT INTO treehousedb.aulas
				(datahora_aula, id_turma)
				VALUES (?, ?)`

	statement, err := l.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(aula)
	if err != nil {
		return 0, err
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertedID), nil
}

func (l Lessons) Update(interface{}) error {
	return nil
}

func (l Lessons) FetchByID(id int) (interface{}, error) {
	return nil, nil
}

func (l Lessons) FetchAll() ([]models.Classes, error) {
	query := `SELECT * FROM treehousedb.turmas`

	lines, err := l.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer lines.Close()

	var classes []models.Classes

	for lines.Next() {
		var class models.Classes

		if err = lines.Scan(
			&class.ID,
			&class.Name,
			&class.Teacher,
			&class.Active,
		); err != nil {
			return nil, err
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func (l Lessons) FetchAllActive() ([]models.Classes, error) {
	query := `SELECT * FROM treehousedb.turmas WHERE ativo = 1`

	lines, err := l.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer lines.Close()

	var classes []models.Classes

	for lines.Next() {
		var class models.Classes

		if err = lines.Scan(
			&class.ID,
			&class.Name,
			&class.Teacher,
			&class.Active,
		); err != nil {
			return nil, err
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func (l Lessons) SelectClassStudents(class models.Classes) ([]models.Students, error) {
	query := `SELECT * FROM treehousedb.alunos_turmas WHERE id_turma = ?`

	lines, err := l.db.Query(query, class.ID)
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

func (l Lessons) SetStudent(classID int, student models.Students) (uint64, error) {
	log.Println("Aluno", student)

	query := `INSERT INTO treehousedb.alunos_turmas
				(id_turma, id_aluno)
				VALUES (?, ?)`

	statement, err := l.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(classID, student.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return uint64(rowsAffected), nil
}

func (l Lessons) RemoveStudent(classID int, student models.Students) (uint64, error) {
	query := `DELETE FROM treehousedb.alunos_turmas WHERE id_turma = ? AND id_aluno = ?`

	statement, err := l.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(classID, student.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return uint64(rowsAffected), nil
}

//func (c Classes) SelectClassTeachers(class models.Classes) ([]models.Teachers, error) {
//	query := `SELECT * FROM treehousedb.professores_turmas WHERE id_turma = ?`
//
//	lines, err := c.db.Query(query, class.ID)
//
//	if err != nil {
//		return nil, err
//	}
//
//	defer lines.Close()
//
//	var teachers []models.Teachers
//
//	for lines.Next() {
//		var teacher models.Teachers
//
//		if err = lines.Scan(
//			&teacher.ID,
//			&teacher.UserID,
//			&teacher.Name,
//			&teacher.Active,
//		); err != nil {
//			return nil, err
//		}
//
//		teachers = append(teachers, teacher)
//	}
//
//	return teachers, nil
//}
