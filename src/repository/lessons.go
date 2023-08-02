package repository

import (
	"database/sql"
	"github.com/brunoob35/TreeHouse-API/src/models"
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

	result, err := statement.Exec(aula.Date)
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

func (l Lessons) FetchByTeacherID(teacherID uint32) ([]models.Lessons, error) {
	query := `SELECT a.*
				FROM treehousedb.aulas AS a
				JOIN treehousedb.alunos_aulas AS aa ON a.id_aula = aa.id_aula
				WHERE aa.id_professor = ?;`

	lines, err := l.db.Query(query, teacherID)
	if err != nil {
		return nil, err
	}

	defer lines.Close()

	var lessons []models.Lessons

	for lines.Next() {
		var lesson models.Lessons

		if err = lines.Scan(
			&lesson.ID,
			&lesson.ClassID,
			&lesson.Status,
			&lesson.Date,
			&lesson.End,
		); err != nil {
			return nil, err
		}

		lessons = append(lessons, lesson)
	}

	return lessons, nil
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

// GetStudentLesson Fetches all students Ids who participate in given lesson
func (l Lessons) GetStudentLesson(lesson models.Lessons) ([]models.Students, error) {
	query := `SELECT * FROM treehousedb.alunos_aulas where id_aula = ?`

	lines, err := l.db.Query(query, lesson.ID)

	if err != nil {
		return nil, err
	}

	defer lines.Close()

	var students []models.Students

	for lines.Next() {
		var student models.Students

		if err = lines.Scan(
			&student.ID,
		); err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

// SetStudentLesson Saves a student ID to a lesson
func (l Lessons) SetStudentLesson(lesson models.Lessons, student models.Students) (uint64, error) {
	query := `INSERT INTO treehousedb.alunos_turmas
				(id_turma, id_aluno)
				VALUES (?, ?)`

	statement, err := l.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(lesson.ID, student.ID)
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

func (l Lessons) FetchByClassID(id int) (interface{}, error) {

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
