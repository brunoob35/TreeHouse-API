package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/models"
)

type ClassesRepository struct {
	db *sql.DB
}

const lessonStatusAulaDada uint64 = 2

func NewClassesRepository(db *sql.DB) *ClassesRepository {
	return &ClassesRepository{db}
}

type classRecurrencePayload struct {
	Weekdays    []string `json:"weekdays"`
	LessonCount int      `json:"lesson_count"`
	StartDate   string   `json:"start_date"`
	StartTime   string   `json:"start_time"`
}

const classSelectFields = `
		t.id,
		t.id_professor,
		t.id_endereco,
		t.nome,
		t.descricao_recorrencia,
		t.recorrencia_json,
		e.id,
		e.cep,
		e.rua,
		e.numero,
		e.bairro,
		e.cidade,
		e.estado,
		e.pais,
		e.complemento,
		t.created_at,
		t.updated_at,
		t.deleted_at
`

const classAddressJoin = `
		LEFT JOIN treehousedb.enderecos e ON e.id = t.id_endereco
`

func (r ClassesRepository) Create(class models.Class) (uint64, int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	addressID, err := upsertClassAddressTx(tx, nil, class.Endereco)
	if err != nil {
		return 0, 0, err
	}

	query := `
		INSERT INTO treehousedb.turmas (
			id_professor,
			id_endereco,
			nome,
			descricao_recorrencia,
			recorrencia_json
		) VALUES (?, ?, ?, ?, ?)
	`

	var teacherID interface{}
	if class.TeacherID != nil {
		teacherID = *class.TeacherID
	} else {
		teacherID = nil
	}

	result, err := tx.Exec(
		query,
		teacherID,
		nullableUint64(addressID),
		class.Name,
		nullIfEmpty(class.RecurrenceDesc),
		nullIfEmpty(class.RecurrenceJSON),
	)
	if err != nil {
		return 0, 0, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	generatedLessonsCount, err := createLessonsFromRecurrence(tx, uint64(lastInsertedID), class)
	if err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, 0, err
	}

	return uint64(lastInsertedID), generatedLessonsCount, nil
}

func createLessonsFromRecurrence(tx *sql.Tx, classID uint64, class models.Class) (int, error) {
	recurrenceRaw := strings.TrimSpace(class.RecurrenceJSON)
	if recurrenceRaw == "" {
		return 0, nil
	}

	var recurrence classRecurrencePayload
	if err := json.Unmarshal([]byte(recurrenceRaw), &recurrence); err != nil {
		return 0, fmt.Errorf("invalid recurrence_json: %w", err)
	}

	lessonDates, err := buildLessonDates(recurrence)
	if err != nil {
		return 0, err
	}

	if len(lessonDates) == 0 {
		return 0, nil
	}

	query := `
		INSERT INTO treehousedb.aulas (
			id_status,
			id_professor,
			id_turma,
			data_aula
		) VALUES (?, ?, ?, ?)
	`

	var teacherID interface{}
	if class.TeacherID != nil {
		teacherID = *class.TeacherID
	} else {
		teacherID = nil
	}

	for _, lessonDate := range lessonDates {
		_, err := tx.Exec(query, 1, teacherID, classID, lessonDate)
		if err != nil {
			return 0, err
		}
	}

	return len(lessonDates), nil
}

func buildLessonDates(recurrence classRecurrencePayload) ([]time.Time, error) {
	if strings.TrimSpace(recurrence.StartDate) == "" {
		return nil, errors.New("recurrence start_date is required")
	}
	if strings.TrimSpace(recurrence.StartTime) == "" {
		return nil, errors.New("recurrence start_time is required")
	}
	if recurrence.LessonCount <= 0 {
		return nil, errors.New("recurrence lesson_count must be greater than zero")
	}
	if len(recurrence.Weekdays) == 0 {
		return nil, errors.New("recurrence weekdays is required")
	}

	startDate, err := time.Parse("2006-01-02", recurrence.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid recurrence start_date: %w", err)
	}

	startTime, err := time.Parse("15:04", recurrence.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid recurrence start_time: %w", err)
	}

	allowedWeekdays := map[time.Weekday]struct{}{}
	for _, day := range recurrence.Weekdays {
		weekday, err := parseWeekday(day)
		if err != nil {
			return nil, err
		}
		allowedWeekdays[weekday] = struct{}{}
	}

	baseDate := time.Date(
		startDate.Year(),
		startDate.Month(),
		startDate.Day(),
		startTime.Hour(),
		startTime.Minute(),
		0,
		0,
		time.UTC,
	)

	var lessonDates []time.Time
	for offsetDays := 0; len(lessonDates) < recurrence.LessonCount && offsetDays < recurrence.LessonCount*14; offsetDays++ {
		candidate := baseDate.AddDate(0, 0, offsetDays)
		if candidate.Before(baseDate) {
			continue
		}

		if _, ok := allowedWeekdays[candidate.Weekday()]; !ok {
			continue
		}

		lessonDates = append(lessonDates, candidate)
	}

	if len(lessonDates) != recurrence.LessonCount {
		return nil, errors.New("could not generate the requested number of lessons from recurrence")
	}

	return lessonDates, nil
}

func parseWeekday(value string) (time.Weekday, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "domingo", "dom", "sunday":
		return time.Sunday, nil
	case "segunda", "seg", "monday":
		return time.Monday, nil
	case "terca", "terça", "ter", "tuesday":
		return time.Tuesday, nil
	case "quarta", "qua", "wednesday":
		return time.Wednesday, nil
	case "quinta", "qui", "thursday":
		return time.Thursday, nil
	case "sexta", "sex", "friday":
		return time.Friday, nil
	case "sabado", "sábado", "sab", "saturday":
		return time.Saturday, nil
	default:
		return time.Sunday, fmt.Errorf("invalid recurrence weekday: %s", value)
	}
}

func (r ClassesRepository) FetchByID(classID uint64) (models.Class, error) {
	query := `
		SELECT ` + classSelectFields + `
		FROM treehousedb.turmas t
		` + classAddressJoin + `
		WHERE t.id = ?
		LIMIT 1
	`

	return scanClassRow(r.db.QueryRow(query, classID))
}

func (r ClassesRepository) FetchAllActive() ([]models.Class, error) {
	query := `
	SELECT
			` + classSelectFields + `,
			(SELECT COUNT(*) FROM treehousedb.alunos_turmas at WHERE at.id_turma = t.id) AS student_count,
			(SELECT COUNT(*) FROM treehousedb.aulas a WHERE a.id_turma = t.id) AS lessons_total,
			(SELECT COUNT(*) FROM treehousedb.aulas a WHERE a.id_turma = t.id AND a.id_status = 2) AS lessons_completed
		FROM treehousedb.turmas t
		LEFT JOIN treehousedb.enderecos e ON e.id = t.id_endereco
		WHERE t.deleted_at IS NULL
		ORDER BY t.nome ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class

	for rows.Next() {
		class, scanErr := scanClassRowWithStats(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func (r ClassesRepository) FetchAllActiveByTeacherID(teacherID uint64) ([]models.Class, error) {
	query := `
	SELECT
			` + classSelectFields + `,
			(SELECT COUNT(*) FROM treehousedb.alunos_turmas at WHERE at.id_turma = t.id) AS student_count,
			(SELECT COUNT(*) FROM treehousedb.aulas a WHERE a.id_turma = t.id) AS lessons_total,
			(SELECT COUNT(*) FROM treehousedb.aulas a WHERE a.id_turma = t.id AND a.id_status = 2) AS lessons_completed
		FROM treehousedb.turmas t
		LEFT JOIN treehousedb.enderecos e ON e.id = t.id_endereco
		WHERE t.deleted_at IS NULL
		  AND t.id_professor = ?
		ORDER BY t.nome ASC
	`

	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class

	for rows.Next() {
		class, scanErr := scanClassRowWithStats(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func (r ClassesRepository) FetchAll() ([]models.Class, error) {
	query := `
	SELECT
			` + classSelectFields + `,
			(SELECT COUNT(*) FROM treehousedb.alunos_turmas at WHERE at.id_turma = t.id) AS student_count,
			(SELECT COUNT(*) FROM treehousedb.aulas a WHERE a.id_turma = t.id) AS lessons_total,
			(SELECT COUNT(*) FROM treehousedb.aulas a WHERE a.id_turma = t.id AND a.id_status = 2) AS lessons_completed
		FROM treehousedb.turmas t
		LEFT JOIN treehousedb.enderecos e ON e.id = t.id_endereco
		ORDER BY t.nome ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class

	for rows.Next() {
		class, scanErr := scanClassRowWithStats(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		classes = append(classes, class)
	}

	return classes, nil
}

func (r ClassesRepository) Update(classID uint64, class models.Class) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	currentClass, err := r.fetchByIDTx(tx, classID)
	if err != nil {
		return err
	}

	addressID, err := upsertClassAddressTx(tx, currentClass.IDEndereco, class.Endereco)
	if err != nil {
		return err
	}

	query := `
		UPDATE treehousedb.turmas
		SET
			id_professor = ?,
			id_endereco = ?,
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

	_, err = tx.Exec(
		query,
		teacherID,
		nullableUint64(addressID),
		class.Name,
		nullIfEmpty(class.RecurrenceDesc),
		nullIfEmpty(class.RecurrenceJSON),
		classID,
	)
	if err != nil {
		return err
	}

	recurrenceChanged := normalizedNullableString(currentClass.RecurrenceJSON) != normalizedNullableString(class.RecurrenceJSON)
	teacherChanged := equalTeacherID(currentClass.TeacherID, class.TeacherID) == false

	if recurrenceChanged {
		if err = deleteOpenLessonsByClassTx(tx, classID); err != nil {
			return err
		}

		if _, err = createLessonsFromRecurrence(tx, classID, class); err != nil {
			return err
		}
	}

	if teacherChanged {
		if err = updateOpenLessonsTeacherByClassTx(tx, classID, class.TeacherID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r ClassesRepository) SoftDelete(classID uint64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		UPDATE treehousedb.turmas
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND deleted_at IS NULL
	`

	_, err = tx.Exec(query, classID)
	if err != nil {
		return err
	}

	if err = cancelOpenLessonsByClassTx(tx, classID); err != nil {
		return err
	}

	return tx.Commit()
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

func (r ClassesRepository) CreatePrivateClassFromStudent(studentID uint64, class models.Class) (uint64, int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var studentName string
	queryStudent := `SELECT nome FROM treehousedb.alunos WHERE id = ? LIMIT 1`

	if err = tx.QueryRow(queryStudent, studentID).Scan(&studentName); err != nil {
		return 0, 0, err
	}

	className := strings.TrimSpace(class.Name)
	if className == "" {
		className = fmt.Sprintf("Turma %s", strings.TrimSpace(studentName))
	}

	queryClass := `
		INSERT INTO treehousedb.turmas (
			id_professor,
			id_endereco,
			nome,
			descricao_recorrencia,
			recorrencia_json
		) VALUES (?, ?, ?, ?, ?)
	`

	var teacher interface{}
	if class.TeacherID != nil {
		teacher = *class.TeacherID
	} else {
		teacher = nil
	}

	addressID, err := upsertClassAddressTx(tx, nil, class.Endereco)
	if err != nil {
		return 0, 0, err
	}

	result, err := tx.Exec(
		queryClass,
		teacher,
		nullableUint64(addressID),
		className,
		nullIfEmpty(class.RecurrenceDesc),
		nullIfEmpty(class.RecurrenceJSON),
	)
	if err != nil {
		return 0, 0, err
	}

	classID, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	class.Name = className
	generatedLessonsCount, err := createLessonsFromRecurrence(tx, uint64(classID), class)
	if err != nil {
		return 0, 0, err
	}

	queryRelation := `
		INSERT INTO treehousedb.alunos_turmas (id_aluno, id_turma)
		VALUES (?, ?)
	`

	if _, err = tx.Exec(queryRelation, studentID, classID); err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, 0, err
	}

	return uint64(classID), generatedLessonsCount, nil
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

	professorIDCopy := professorID
	if err = updateOpenLessonsTeacherByClassTx(tx, classID, &professorIDCopy); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r ClassesRepository) fetchByIDTx(tx *sql.Tx, classID uint64) (models.Class, error) {
	query := `
		SELECT ` + classSelectFields + `
		FROM treehousedb.turmas t
		` + classAddressJoin + `
		WHERE t.id = ?
		LIMIT 1
	`

	return scanClassRow(tx.QueryRow(query, classID))
}

func deleteOpenLessonsByClassTx(tx *sql.Tx, classID uint64) error {
	queryDeleteRelations := `
		DELETE aa
		FROM treehousedb.alunos_aulas aa
		INNER JOIN treehousedb.aulas a ON a.id = aa.id_aula
		WHERE a.id_turma = ?
		  AND a.id_status <> ?
	`

	if _, err := tx.Exec(queryDeleteRelations, classID, lessonStatusAulaDada); err != nil {
		return err
	}

	queryDeleteLessons := `
		DELETE FROM treehousedb.aulas
		WHERE id_turma = ?
		  AND id_status <> ?
	`

	_, err := tx.Exec(queryDeleteLessons, classID, lessonStatusAulaDada)
	return err
}

func cancelOpenLessonsByClassTx(tx *sql.Tx, classID uint64) error {
	query := `
		UPDATE treehousedb.aulas
		SET
			id_status = 3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_turma = ?
		  AND id_status <> ?
	`

	_, err := tx.Exec(query, classID, lessonStatusAulaDada)
	return err
}

func updateOpenLessonsTeacherByClassTx(tx *sql.Tx, classID uint64, teacherID *uint64) error {
	query := `
		UPDATE treehousedb.aulas
		SET
			id_professor = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_turma = ?
		  AND id_status <> ?
	`

	var teacher interface{}
	if teacherID != nil {
		teacher = *teacherID
	}

	_, err := tx.Exec(query, teacher, classID, lessonStatusAulaDada)
	return err
}

func equalTeacherID(left *uint64, right *uint64) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return *left == *right
}

func normalizedNullableString(value string) string {
	return strings.TrimSpace(value)
}

type classRowScanner interface {
	Scan(dest ...interface{}) error
}

func scanClassRow(scanner classRowScanner) (models.Class, error) {
	var class models.Class
	var teacherID sql.NullInt64
	var addressID sql.NullInt64
	var recurrenceDesc sql.NullString
	var recurrenceJSON sql.NullString
	var deletedAt sql.NullTime
	var address models.Address
	var scannedAddressID sql.NullInt64
	var cep sql.NullString
	var rua sql.NullString
	var numero sql.NullString
	var bairro sql.NullString
	var cidade sql.NullString
	var estado sql.NullString
	var pais sql.NullString
	var complemento sql.NullString

	err := scanner.Scan(
		&class.ID,
		&teacherID,
		&addressID,
		&class.Name,
		&recurrenceDesc,
		&recurrenceJSON,
		&scannedAddressID,
		&cep,
		&rua,
		&numero,
		&bairro,
		&cidade,
		&estado,
		&pais,
		&complemento,
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

	if addressID.Valid {
		aid := uint64(addressID.Int64)
		class.IDEndereco = &aid
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

	if scannedAddressID.Valid {
		address.ID = uint64(scannedAddressID.Int64)
		address.CEP = cep.String
		address.Rua = rua.String
		address.Numero = numero.String
		address.Bairro = bairro.String
		address.Cidade = cidade.String
		address.Estado = estado.String
		address.Pais = pais.String
		address.Complemento = complemento.String
		class.Endereco = &address
	}

	return class, nil
}

func scanClassRowWithStats(scanner classRowScanner) (models.Class, error) {
	var class models.Class
	var teacherID sql.NullInt64
	var addressID sql.NullInt64
	var recurrenceDesc sql.NullString
	var recurrenceJSON sql.NullString
	var deletedAt sql.NullTime
	var address models.Address
	var scannedAddressID sql.NullInt64
	var cep sql.NullString
	var rua sql.NullString
	var numero sql.NullString
	var bairro sql.NullString
	var cidade sql.NullString
	var estado sql.NullString
	var pais sql.NullString
	var complemento sql.NullString

	err := scanner.Scan(
		&class.ID,
		&teacherID,
		&addressID,
		&class.Name,
		&recurrenceDesc,
		&recurrenceJSON,
		&scannedAddressID,
		&cep,
		&rua,
		&numero,
		&bairro,
		&cidade,
		&estado,
		&pais,
		&complemento,
		&class.CreatedAt,
		&class.UpdatedAt,
		&deletedAt,
		&class.StudentCount,
		&class.LessonsTotal,
		&class.LessonsCompleted,
	)
	if err != nil {
		return models.Class{}, err
	}

	if teacherID.Valid {
		tid := uint64(teacherID.Int64)
		class.TeacherID = &tid
	}

	if addressID.Valid {
		aid := uint64(addressID.Int64)
		class.IDEndereco = &aid
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

	if scannedAddressID.Valid {
		address.ID = uint64(scannedAddressID.Int64)
		address.CEP = cep.String
		address.Rua = rua.String
		address.Numero = numero.String
		address.Bairro = bairro.String
		address.Cidade = cidade.String
		address.Estado = estado.String
		address.Pais = pais.String
		address.Complemento = complemento.String
		class.Endereco = &address
	}

	return class, nil
}

func upsertClassAddressTx(tx *sql.Tx, currentAddressID *uint64, address *models.Address) (*uint64, error) {
	if address == nil {
		return nil, nil
	}

	if currentAddressID != nil {
		query := `
			UPDATE treehousedb.enderecos
			SET
				cep = ?,
				rua = ?,
				numero = ?,
				bairro = ?,
				cidade = ?,
				estado = ?,
				pais = ?,
				complemento = ?,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`

		_, err := tx.Exec(
			query,
			nullIfEmpty(address.CEP),
			address.Rua,
			address.Numero,
			address.Bairro,
			address.Cidade,
			address.Estado,
			address.Pais,
			nullIfEmpty(address.Complemento),
			*currentAddressID,
		)
		if err != nil {
			return nil, err
		}

		return currentAddressID, nil
	}

	query := `
		INSERT INTO treehousedb.enderecos (
			cep,
			rua,
			numero,
			bairro,
			cidade,
			estado,
			pais,
			complemento
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		query,
		nullIfEmpty(address.CEP),
		address.Rua,
		address.Numero,
		address.Bairro,
		address.Cidade,
		address.Estado,
		address.Pais,
		nullIfEmpty(address.Complemento),
	)
	if err != nil {
		return nil, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	addressID := uint64(lastInsertedID)
	return &addressID, nil
}

func nullableUint64(value *uint64) interface{} {
	if value == nil {
		return nil
	}

	return *value
}
