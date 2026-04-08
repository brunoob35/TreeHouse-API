package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/models"
)

type ClassesRepository struct {
	db *sql.DB
}

func NewClassesRepository(db *sql.DB) *ClassesRepository {
	return &ClassesRepository{db}
}

func (r ClassesRepository) Create(class models.Class) (uint64, error) {
	query := `
		INSERT INTO treehousedb.turmas (
			id_professor,
			nome,
			descricao_recorrencia,
			recorrencia_json
		) VALUES (?, ?, ?, ?)
	`

	var teacherID interface{}
	if class.TeacherID != nil {
		teacherID = *class.TeacherID
	} else {
		teacherID = nil
	}

	result, err := r.db.Exec(
		query,
		teacherID,
		class.Name,
		nullIfEmpty(class.RecurrenceDesc),
		nullIfEmpty(class.RecurrenceJSON),
	)
	if err != nil {
		return 0, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertedID), nil
}

func (r ClassesRepository) FetchByID(classID uint64) (models.Class, error) {
	query := `
		SELECT
			id,
			id_professor,
			nome,
			descricao_recorrencia,
			recorrencia_json,
			created_at,
			updated_at,
			deleted_at
		FROM treehousedb.turmas
		WHERE id = ?
		LIMIT 1
	`

	var class models.Class
	var teacherID sql.NullInt64
	var recurrenceDesc sql.NullString
	var recurrenceJSON sql.NullString
	var deletedAt sql.NullTime

	err := r.db.QueryRow(query, classID).Scan(
		&class.ID,
		&teacherID,
		&class.Name,
		&recurrenceDesc,
		&recurrenceJSON,
		&class.CreatedAt,
		&class.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return models.Class{}, err
	}

	if teacherID.Valid {
		tid := uint64(teacherID.Int64)
		class.TeacherID = &tid
	}

	if recurrenceDesc.Valid {
		class.RecurrenceDesc = recurrenceDesc.String
	}

	if recurrenceJSON.Valid {
		class.RecurrenceJSON = recurrenceJSON.String
	}

	if deletedAt.Valid {
		class.DeletedAt = &deletedAt.Time
	}

	return class, nil
}

func (r ClassesRepository) FetchAllActive() ([]models.Class, error) {
	query := `
		SELECT
			id,
			id_professor,
			nome,
			descricao_recorrencia,
			recorrencia_json,
			created_at,
			updated_at,
			deleted_at
		FROM treehousedb.turmas
		WHERE deleted_at IS NULL
		ORDER BY nome ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class

	for rows.Next() {
		var class models.Class
		var teacherID sql.NullInt64
		var recurrenceDesc sql.NullString
		var recurrenceJSON sql.NullString
		var deletedAt sql.NullTime

		if err = rows.Scan(
			&class.ID,
			&teacherID,
			&class.Name,
			&recurrenceDesc,
			&recurrenceJSON,
			&class.CreatedAt,
			&class.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}

		if teacherID.Valid {
			tid := uint64(teacherID.Int64)
			class.TeacherID = &tid
		}

		if recurrenceDesc.Valid {
			class.RecurrenceDesc = recurrenceDesc.String
		}

		if recurrenceJSON.Valid {
			class.RecurrenceJSON = recurrenceJSON.String
		}

		if deletedAt.Valid {
			class.DeletedAt = &deletedAt.Time
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func (r ClassesRepository) FetchAll() ([]models.Class, error) {
	query := `
		SELECT
			id,
			id_professor,
			nome,
			descricao_recorrencia,
			recorrencia_json,
			created_at,
			updated_at,
			deleted_at
		FROM treehousedb.turmas
		ORDER BY nome ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class

	for rows.Next() {
		var class models.Class
		var teacherID sql.NullInt64
		var recurrenceDesc sql.NullString
		var recurrenceJSON sql.NullString
		var deletedAt sql.NullTime

		if err = rows.Scan(
			&class.ID,
			&teacherID,
			&class.Name,
			&recurrenceDesc,
			&recurrenceJSON,
			&class.CreatedAt,
			&class.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}

		if teacherID.Valid {
			tid := uint64(teacherID.Int64)
			class.TeacherID = &tid
		}

		if recurrenceDesc.Valid {
			class.RecurrenceDesc = recurrenceDesc.String
		}

		if recurrenceJSON.Valid {
			class.RecurrenceJSON = recurrenceJSON.String
		}

		if deletedAt.Valid {
			class.DeletedAt = &deletedAt.Time
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func (r ClassesRepository) Update(classID uint64, class models.Class) error {
	query := `
		UPDATE treehousedb.turmas
		SET
			id_professor = ?,
			nome = ?,
			descricao_recorrencia = ?,
			recorrencia_json = ?
		WHERE id = ?
		  AND deleted_at IS NULL
	`

	var teacherID interface{}
	if class.TeacherID != nil {
		teacherID = *class.TeacherID
	} else {
		teacherID = nil
	}

	_, err := r.db.Exec(
		query,
		teacherID,
		class.Name,
		nullIfEmpty(class.RecurrenceDesc),
		nullIfEmpty(class.RecurrenceJSON),
		classID,
	)

	return err
}

func (r ClassesRepository) SoftDelete(classID uint64) error {
	query := `
		UPDATE treehousedb.turmas
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, classID)
	return err
}

func (r ClassesRepository) AddStudent(classID uint64, studentID uint64) error {
	query := `
		INSERT INTO treehousedb.alunos_turmas (id_aluno, id_turma)
		VALUES (?, ?)
	`

	_, err := r.db.Exec(query, studentID, classID)
	return err
}

func (r ClassesRepository) RemoveStudent(classID uint64, studentID uint64) error {
	query := `
		DELETE FROM treehousedb.alunos_turmas
		WHERE id_aluno = ?
		  AND id_turma = ?
	`

	_, err := r.db.Exec(query, studentID, classID)
	return err
}

func (r ClassesRepository) FetchStudents(classID uint64) ([]models.Student, error) {
	query := `
		SELECT
			a.id,
			a.nome,
			a.livro,
			a.alfabetizacao,
			a.nascimento,
			a.ativo,
			a.created_at,
			a.updated_at
		FROM treehousedb.alunos a
		INNER JOIN treehousedb.alunos_turmas at ON at.id_aluno = a.id
		WHERE at.id_turma = ?
		ORDER BY a.nome ASC
	`

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student
		var livro sql.NullString
		var alfabetizacao sql.NullString
		var nascimento sql.NullTime

		if err = rows.Scan(
			&student.ID,
			&student.Nome,
			&livro,
			&alfabetizacao,
			&nascimento,
			&student.Ativo,
			&student.CreatedAt,
			&student.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if livro.Valid {
			student.Livro = livro.String
		}

		if alfabetizacao.Valid {
			student.Alfabetizacao = alfabetizacao.String
		}

		if nascimento.Valid {
			student.Nascimento = &nascimento.Time
		}

		students = append(students, student)
	}

	return students, nil
}

func (r ClassesRepository) CreatePrivateClassFromStudent(studentID uint64, teacherID *uint64) (uint64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var studentName string
	queryStudent := `SELECT nome FROM treehousedb.alunos WHERE id = ? LIMIT 1`

	if err = tx.QueryRow(queryStudent, studentID).Scan(&studentName); err != nil {
		return 0, err
	}

	className := fmt.Sprintf("%s - Particular", strings.TrimSpace(studentName))

	queryClass := `
		INSERT INTO treehousedb.turmas (
			id_professor,
			nome
		) VALUES (?, ?)
	`

	var teacher interface{}
	if teacherID != nil {
		teacher = *teacherID
	} else {
		teacher = nil
	}

	result, err := tx.Exec(queryClass, teacher, className)
	if err != nil {
		return 0, err
	}

	classID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	queryRelation := `
		INSERT INTO treehousedb.alunos_turmas (id_aluno, id_turma)
		VALUES (?, ?)
	`

	if _, err = tx.Exec(queryRelation, studentID, classID); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(classID), nil
}

func (r *ClassesRepository) AssignProfessorToClass(classID, professorID uint64) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var professorExists bool
	queryProfessor := `
		SELECT EXISTS(
			SELECT 1
			FROM treehousedb.usuarios u
			INNER JOIN treehousedb.usuarios_permissoes up
				ON up.id_usuario = u.id
			WHERE u.id = ?
			  AND u.ativo = TRUE
			  AND up.id_permissao = ?
		)
	`

	if err = tx.QueryRow(queryProfessor, professorID, authentication.PermProfessor).Scan(&professorExists); err != nil {
		return err
	}
	if !professorExists {
		return fmt.Errorf("professor não encontrado ou inválido")
	}

	var classExists bool
	queryClass := `
		SELECT EXISTS(
			SELECT 1
			FROM treehousedb.turmas
			WHERE id = ?
		)
	`

	if err = tx.QueryRow(queryClass, classID).Scan(&classExists); err != nil {
		return err
	}
	if !classExists {
		return fmt.Errorf("turma não encontrada")
	}

	queryUpdate := `
		UPDATE treehousedb.turmas
		SET
			id_professor = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := tx.Exec(queryUpdate, professorID, classID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("nenhuma turma encontrada com id %d", classID)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
