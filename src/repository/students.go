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
	db          *sql.DB
	auditUserID *uint64
}

// NewStudentsRepository creates a new repository instance
// bound to the provided database connection.
func NewStudentsRepository(db *sql.DB) *StudentsRepository {
	return &StudentsRepository{db: db}
}

func (r *StudentsRepository) WithAuditUser(userID uint64) *StudentsRepository {
	r.auditUserID = &userID
	return r
}

// Insert creates a new student record in the database.
// Only base student data is inserted. The "ativo" field must
// already be defined by the caller (usually true on creation).
func (r *StudentsRepository) Insert(student models.Student) (uint64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = setAuditUserTx(tx, r.auditUserID); err != nil {
		return 0, err
	}

	query := `
		INSERT INTO treehousedb.alunos (
			nome,
			livro,
			alfabetizacao,
			nascimento,
			ativo
		) VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
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

	if err = tx.Commit(); err != nil {
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
			a.id,
			a.nome,
			a.livro,
			a.alfabetizacao,
			a.nascimento,
			COALESCE(MAX(c.nome), '') AS responsavel,
			COALESCE(MAX(c.telefone), '') AS responsavel_telefone,
			COUNT(DISTINCT at.id_turma) AS classes_count,
			a.ativo,
			a.created_at,
			a.updated_at
		FROM treehousedb.alunos a
		LEFT JOIN treehousedb.clientes_alunos ca
			ON ca.id_aluno = a.id
		LEFT JOIN treehousedb.clientes c
			ON c.id = ca.id_cliente
		LEFT JOIN treehousedb.alunos_turmas at
			ON at.id_aluno = a.id
	`

	var args []interface{}

	if nome != "" {
		query += " WHERE LOWER(a.nome) LIKE ?"
		args = append(args, "%"+nome+"%")
	}

	query += `
		GROUP BY
			a.id,
			a.nome,
			a.livro,
			a.alfabetizacao,
			a.nascimento,
			a.ativo,
			a.created_at,
			a.updated_at
		ORDER BY a.nome
	`

	rows, err := r.db.Query(query, args...)
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

		err = rows.Scan(
			&student.ID,
			&student.Nome,
			&livro,
			&alfabetizacao,
			&nascimento,
			&student.Responsavel,
			&student.ResponsavelTelefone,
			&student.ClassesCount,
			&student.Ativo,
			&student.CreatedAt,
			&student.UpdatedAt,
		)
		if err != nil {
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
	var livro sql.NullString
	var alfabetizacao sql.NullString
	var nascimento sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&student.ID,
		&student.Nome,
		&livro,
		&alfabetizacao,
		&nascimento,
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

	if livro.Valid {
		student.Livro = livro.String
	}
	if alfabetizacao.Valid {
		student.Alfabetizacao = alfabetizacao.String
	}
	if nascimento.Valid {
		student.Nascimento = &nascimento.Time
	}

	return student, nil
}

// Update updates an existing student record.
//
// All editable fields can be updated including the
// active status.
func (r *StudentsRepository) Update(id uint64, student models.Student) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = setAuditUserTx(tx, r.auditUserID); err != nil {
		return err
	}

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

	result, err := tx.Exec(
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

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// SoftDelete performs a soft delete on a student record.
//
// Instead of removing the row from the database,
// the function sets the "ativo" field to false.
func (r *StudentsRepository) SoftDelete(id uint64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = setAuditUserTx(tx, r.auditUserID); err != nil {
		return err
	}

	query := `
		UPDATE treehousedb.alunos
		SET
			ativo = FALSE,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := tx.Exec(query, id)
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

	if err = tx.Commit(); err != nil {
		return err
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
			(SELECT COUNT(*) FROM treehousedb.alunos_turmas at2 WHERE at2.id_turma = t.id) AS student_count,
			(SELECT COUNT(*) FROM treehousedb.aulas a WHERE a.id_turma = t.id) AS lessons_total,
			(SELECT COUNT(*) FROM treehousedb.aulas a WHERE a.id_turma = t.id AND a.id_status = 2) AS lessons_completed,
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
			&class.StudentCount,
			&class.LessonsTotal,
			&class.LessonsCompleted,
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

func (r *StudentsRepository) FetchByTeacherID(teacherID uint64) ([]models.ProfessorStudentSummary, error) {
	query := `
		SELECT
			a.id,
			a.nome,
			a.ativo,
			t.id,
			t.nome,
			COUNT(DISTINCT CASE WHEN aa.id_aula IS NOT NULL AND aula.id_status = 2 THEN aula.id END) AS lessons_completed,
			COUNT(DISTINCT aula.id) AS lessons_total
		FROM treehousedb.turmas t
		INNER JOIN treehousedb.alunos_turmas at
			ON at.id_turma = t.id
		INNER JOIN treehousedb.alunos a
			ON a.id = at.id_aluno
		LEFT JOIN treehousedb.aulas aula
			ON aula.id_turma = t.id
		LEFT JOIN treehousedb.alunos_aulas aa
			ON aa.id_aula = aula.id
		   AND aa.id_aluno = a.id
		WHERE t.deleted_at IS NULL
		  AND t.id_professor = ?
		GROUP BY
			a.id,
			a.nome,
			a.ativo,
			t.id,
			t.nome
		ORDER BY a.nome ASC, t.nome ASC
	`

	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.ProfessorStudentSummary

	for rows.Next() {
		var student models.ProfessorStudentSummary

		if err = rows.Scan(
			&student.ID,
			&student.Nome,
			&student.Ativo,
			&student.ClassID,
			&student.ClassName,
			&student.LessonsCompleted,
			&student.LessonsTotal,
		); err != nil {
			return nil, err
		}

		if student.LessonsTotal > 0 {
			student.FrequencyPercentage = (float64(student.LessonsCompleted) / float64(student.LessonsTotal)) * 100
		}

		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
