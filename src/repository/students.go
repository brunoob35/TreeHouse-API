package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/brunoob35/TreeHouse-API/src/models"
)

// StudentsRepository is responsible for all database operations
// related to student records.
type StudentsRepository struct {
	db *sql.DB
}

// NewStudentsRepository creates a new repository instance
// bound to the provided database connection.
func NewStudentsRepository(db *sql.DB) *StudentsRepository {
	return &StudentsRepository{db: db}
}

// Insert creates a new student record in the database.
// Only base student data is inserted. The "ativo" field must
// already be defined by the caller (usually true on creation).
func (r *StudentsRepository) Insert(student models.Student) (uint64, error) {
	query := `
		INSERT INTO treehousedb.alunos (
			nome,
			livro,
			alfabetizacao,
			nascimento,
			ativo
		) VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		student.Nome,
		student.Livro,
		student.Alfabetizacao,
		student.Nascimento,
		student.Ativo,
	)
	if err != nil {
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(insertedID), nil
}

// FetchAll returns a list of students optionally filtered by name.
// If the "nome" parameter is provided, a case-insensitive
// search will be performed using a LIKE clause.
func (r *StudentsRepository) FetchAll(nome string) ([]models.Student, error) {
	query := `
		SELECT
			id,
			nome,
			livro,
			alfabetizacao,
			nascimento,
			ativo,
			created_at,
			updated_at
		FROM treehousedb.alunos
	`

	var args []interface{}

	if nome != "" {
		query += " WHERE LOWER(nome) LIKE ?"
		args = append(args, "%"+nome+"%")
	}

	query += " ORDER BY nome"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student

		err = rows.Scan(
			&student.ID,
			&student.Nome,
			&student.Livro,
			&student.Alfabetizacao,
			&student.Nascimento,
			&student.Ativo,
			&student.CreatedAt,
			&student.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// FetchByID retrieves a single student by its ID.
//
// If no record exists, sql.ErrNoRows is returned.
func (r *StudentsRepository) FetchByID(id uint64) (models.Student, error) {
	query := `
		SELECT
			id,
			nome,
			livro,
			alfabetizacao,
			nascimento,
			ativo,
			created_at,
			updated_at
		FROM treehousedb.alunos
		WHERE id = ?
	`

	var student models.Student

	err := r.db.QueryRow(query, id).Scan(
		&student.ID,
		&student.Nome,
		&student.Livro,
		&student.Alfabetizacao,
		&student.Nascimento,
		&student.Ativo,
		&student.CreatedAt,
		&student.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Student{}, sql.ErrNoRows
		}
		return models.Student{}, err
	}

	return student, nil
}

// Update updates an existing student record.
//
// All editable fields can be updated including the
// active status.
func (r *StudentsRepository) Update(id uint64, student models.Student) error {
	query := `
		UPDATE treehousedb.alunos
		SET
			nome = ?,
			livro = ?,
			alfabetizacao = ?,
			nascimento = ?,
			ativo = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.Exec(
		query,
		student.Nome,
		student.Livro,
		student.Alfabetizacao,
		student.Nascimento,
		student.Ativo,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum aluno encontrado com id %d", id)
	}

	return nil
}

// SoftDelete performs a soft delete on a student record.
//
// Instead of removing the row from the database,
// the function sets the "ativo" field to false.
func (r *StudentsRepository) SoftDelete(id uint64) error {
	query := `
		UPDATE treehousedb.alunos
		SET
			ativo = FALSE,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum aluno encontrado com id %d", id)
	}

	return nil
}

// FetchClasses returns every class associated with a student,
// including inactive ones for historical review.
func (r *StudentsRepository) FetchClasses(studentID uint64) ([]models.Class, error) {
	query := `
		SELECT
			t.id,
			t.id_professor,
			t.nome,
			t.descricao_recorrencia,
			t.recorrencia_json,
			t.created_at,
			t.updated_at,
			t.deleted_at
		FROM treehousedb.turmas t
		INNER JOIN treehousedb.alunos_turmas at
			ON at.id_turma = t.id
		WHERE at.id_aluno = ?
		ORDER BY
			CASE WHEN t.deleted_at IS NULL THEN 0 ELSE 1 END,
			t.nome ASC
	`

	rows, err := r.db.Query(query, studentID)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return classes, nil
}
