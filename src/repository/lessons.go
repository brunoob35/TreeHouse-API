package repositories

import (
	"database/sql"
	"time"

	"github.com/brunoob35/TreeHouse-API/src/models"
)

type LessonsRepository struct {
	db *sql.DB
}

func NewLessonsRepository(db *sql.DB) *LessonsRepository {
	return &LessonsRepository{db}
}

func (r LessonsRepository) FetchStatuses() ([]models.LessonStatus, error) {
	query := `
		SELECT id, nome_status
		FROM treehousedb.aulas_status
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []models.LessonStatus
	for rows.Next() {
		var status models.LessonStatus
		if err = rows.Scan(&status.ID, &status.Name); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, rows.Err()
}

func (r LessonsRepository) Create(lesson models.Lesson) (uint64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		INSERT INTO treehousedb.aulas (
			id_status,
			id_professor,
			id_turma,
			assunto,
			vocabulario,
			saldo,
			observacoes,
			data_aula
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	var teacherID interface{}
	if lesson.TeacherID != nil {
		teacherID = *lesson.TeacherID
	} else {
		teacherID = nil
	}

	result, err := tx.Exec(
		query,
		1,
		teacherID,
		lesson.ClassID,
		nullIfEmpty(lesson.Subject),
		nullIfEmpty(lesson.Vocabulary),
		nullIfEmpty(lesson.Balance),
		nullIfEmpty(lesson.Notes),
		lesson.LessonDate,
	)
	if err != nil {
		return 0, err
	}

	lessonID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	queryClassStudents := `
		SELECT id_aluno
		FROM treehousedb.alunos_turmas
		WHERE id_turma = ?
	`

	rows, err := tx.Query(queryClassStudents, lesson.ClassID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var studentIDs []uint64
	for rows.Next() {
		var studentID uint64
		if err = rows.Scan(&studentID); err != nil {
			return 0, err
		}
		studentIDs = append(studentIDs, studentID)
	}

	queryInsertLessonStudent := `
		INSERT INTO treehousedb.alunos_aulas (
			id_aluno,
			id_aula,
			origem_registro
		) VALUES (?, ?, 'automatica')
	`

	for _, studentID := range studentIDs {
		if _, err = tx.Exec(queryInsertLessonStudent, studentID, lessonID); err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(lessonID), nil
}

func (r LessonsRepository) FetchByID(lessonID uint64) (models.Lesson, error) {
	query := `
		SELECT
			id,
			id_status,
			(SELECT nome_status FROM treehousedb.aulas_status WHERE id = a.id_status) AS status_name,
			id_professor,
			id_turma,
			assunto,
			vocabulario,
			saldo,
			observacoes,
			data_aula,
			data_aula_original,
			data_aula_solicitada,
			created_at,
			updated_at
		FROM treehousedb.aulas a
		WHERE id = ?
		LIMIT 1
	`

	var lesson models.Lesson
	var teacherID sql.NullInt64
	var statusName sql.NullString
	var subject sql.NullString
	var vocabulary sql.NullString
	var balance sql.NullString
	var notes sql.NullString
	var originalLessonDate sql.NullTime
	var requestedLessonDate sql.NullTime

	err := r.db.QueryRow(query, lessonID).Scan(
		&lesson.ID,
		&lesson.StatusID,
		&statusName,
		&teacherID,
		&lesson.ClassID,
		&subject,
		&vocabulary,
		&balance,
		&notes,
		&lesson.LessonDate,
		&originalLessonDate,
		&requestedLessonDate,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)
	if err != nil {
		return models.Lesson{}, err
	}

	if teacherID.Valid {
		tid := uint64(teacherID.Int64)
		lesson.TeacherID = &tid
	}
	if statusName.Valid {
		lesson.StatusName = statusName.String
	}

	if subject.Valid {
		lesson.Subject = subject.String
	}
	if vocabulary.Valid {
		lesson.Vocabulary = vocabulary.String
	}
	if balance.Valid {
		lesson.Balance = balance.String
	}
	if notes.Valid {
		lesson.Notes = notes.String
	}
	if originalLessonDate.Valid {
		original := originalLessonDate.Time
		lesson.OriginalLessonDate = &original
	}
	if requestedLessonDate.Valid {
		requested := requestedLessonDate.Time
		lesson.RequestedLessonDate = &requested
	}

	return lesson, nil
}

func (r LessonsRepository) FetchAll() ([]models.Lesson, error) {
	query := `
		SELECT
			id,
			id_status,
			(SELECT nome_status FROM treehousedb.aulas_status WHERE id = a.id_status) AS status_name,
			id_professor,
			id_turma,
			assunto,
			vocabulario,
			saldo,
			observacoes,
			data_aula,
			data_aula_original,
			data_aula_solicitada,
			created_at,
			updated_at
		FROM treehousedb.aulas a
		ORDER BY data_aula DESC, id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.Lesson

	for rows.Next() {
		var lesson models.Lesson
		var teacherID sql.NullInt64
		var statusName sql.NullString
		var subject sql.NullString
		var vocabulary sql.NullString
		var balance sql.NullString
		var notes sql.NullString
		var originalLessonDate sql.NullTime
		var requestedLessonDate sql.NullTime

		if err = rows.Scan(
			&lesson.ID,
			&lesson.StatusID,
			&statusName,
			&teacherID,
			&lesson.ClassID,
			&subject,
			&vocabulary,
			&balance,
			&notes,
			&lesson.LessonDate,
			&originalLessonDate,
			&requestedLessonDate,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if teacherID.Valid {
			tid := uint64(teacherID.Int64)
			lesson.TeacherID = &tid
		}
		if statusName.Valid {
			lesson.StatusName = statusName.String
		}
		if subject.Valid {
			lesson.Subject = subject.String
		}
		if vocabulary.Valid {
			lesson.Vocabulary = vocabulary.String
		}
		if balance.Valid {
			lesson.Balance = balance.String
		}
		if notes.Valid {
			lesson.Notes = notes.String
		}
		if originalLessonDate.Valid {
			original := originalLessonDate.Time
			lesson.OriginalLessonDate = &original
		}
		if requestedLessonDate.Valid {
			requested := requestedLessonDate.Time
			lesson.RequestedLessonDate = &requested
		}

		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (r LessonsRepository) FetchByClass(classID uint64) ([]models.Lesson, error) {
	query := `
		SELECT
			id,
			id_status,
			(SELECT nome_status FROM treehousedb.aulas_status WHERE id = a.id_status) AS status_name,
			id_professor,
			id_turma,
			assunto,
			vocabulario,
			saldo,
			observacoes,
			data_aula,
			data_aula_original,
			data_aula_solicitada,
			created_at,
			updated_at
		FROM treehousedb.aulas a
		WHERE id_turma = ?
		ORDER BY data_aula ASC, id ASC
	`

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.Lesson

	for rows.Next() {
		var lesson models.Lesson
		var teacherID sql.NullInt64
		var statusName sql.NullString
		var subject sql.NullString
		var vocabulary sql.NullString
		var balance sql.NullString
		var notes sql.NullString
		var originalLessonDate sql.NullTime
		var requestedLessonDate sql.NullTime

		if err = rows.Scan(
			&lesson.ID,
			&lesson.StatusID,
			&statusName,
			&teacherID,
			&lesson.ClassID,
			&subject,
			&vocabulary,
			&balance,
			&notes,
			&lesson.LessonDate,
			&originalLessonDate,
			&requestedLessonDate,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if teacherID.Valid {
			tid := uint64(teacherID.Int64)
			lesson.TeacherID = &tid
		}
		if statusName.Valid {
			lesson.StatusName = statusName.String
		}
		if subject.Valid {
			lesson.Subject = subject.String
		}
		if vocabulary.Valid {
			lesson.Vocabulary = vocabulary.String
		}
		if balance.Valid {
			lesson.Balance = balance.String
		}
		if notes.Valid {
			lesson.Notes = notes.String
		}
		if originalLessonDate.Valid {
			original := originalLessonDate.Time
			lesson.OriginalLessonDate = &original
		}
		if requestedLessonDate.Valid {
			requested := requestedLessonDate.Time
			lesson.RequestedLessonDate = &requested
		}

		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (r LessonsRepository) Update(lessonID uint64, lesson models.Lesson) error {
	query := `
		UPDATE treehousedb.aulas
		SET
			id_professor = ?,
			assunto = ?,
			vocabulario = ?,
			saldo = ?,
			observacoes = ?,
			data_aula = ?
		WHERE id = ?
	`

	var teacherID interface{}
	if lesson.TeacherID != nil {
		teacherID = *lesson.TeacherID
	} else {
		teacherID = nil
	}

	_, err := r.db.Exec(
		query,
		teacherID,
		nullIfEmpty(lesson.Subject),
		nullIfEmpty(lesson.Vocabulary),
		nullIfEmpty(lesson.Balance),
		nullIfEmpty(lesson.Notes),
		lesson.LessonDate,
		lessonID,
	)

	return err
}

func (r LessonsRepository) Delete(lessonID uint64) error {
	query := `
		DELETE FROM treehousedb.aulas
		WHERE id = ?
	`

	_, err := r.db.Exec(query, lessonID)
	return err
}

func (r LessonsRepository) AddStudent(lessonID uint64, studentID uint64, note string) error {
	query := `
		INSERT INTO treehousedb.alunos_aulas (
			id_aluno,
			id_aula,
			origem_registro,
			observacao_aluno_aula
		) VALUES (?, ?, 'manual', ?)
	`

	_, err := r.db.Exec(query, studentID, lessonID, nullIfEmpty(note))
	return err
}

func (r LessonsRepository) RemoveStudent(lessonID uint64, studentID uint64) error {
	query := `
		DELETE FROM treehousedb.alunos_aulas
		WHERE id_aluno = ?
		  AND id_aula = ?
	`

	_, err := r.db.Exec(query, studentID, lessonID)
	return err
}

func (r LessonsRepository) FetchStudents(lessonID uint64) ([]models.Student, error) {
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
		INNER JOIN treehousedb.alunos_aulas aa ON aa.id_aluno = a.id
		WHERE aa.id_aula = ?
		ORDER BY a.nome ASC
	`

	rows, err := r.db.Query(query, lessonID)
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

func (r LessonsRepository) UpdateStatus(lessonID uint64, statusID uint64) error {
	query := `
		UPDATE treehousedb.aulas
		SET id_status = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, statusID, lessonID)
	return err
}

func (r LessonsRepository) RequestReschedule(lessonID uint64, requestedDate time.Time) error {
	query := `
		UPDATE treehousedb.aulas
		SET
			id_status = 5,
			data_aula_original = COALESCE(data_aula_original, data_aula),
			data_aula_solicitada = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, requestedDate, lessonID)
	return err
}

func (r LessonsRepository) ApproveReschedule(lessonID uint64) error {
	query := `
		UPDATE treehousedb.aulas
		SET
			data_aula = data_aula_solicitada,
			id_status = 4,
			data_aula_original = NULL,
			data_aula_solicitada = NULL
		WHERE id = ?
		  AND data_aula_solicitada IS NOT NULL
	`

	_, err := r.db.Exec(query, lessonID)
	return err
}

func (r LessonsRepository) RejectReschedule(lessonID uint64) error {
	query := `
		UPDATE treehousedb.aulas
		SET
			id_status = 1,
			data_aula_original = NULL,
			data_aula_solicitada = NULL
		WHERE id = ?
	`

	_, err := r.db.Exec(query, lessonID)
	return err
}
