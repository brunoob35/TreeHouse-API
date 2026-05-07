package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/brunoob35/TreeHouse-API/src/models"
)

type ContractsRepository struct {
	db *sql.DB
}

func NewContractsRepository(db *sql.DB) *ContractsRepository {
	return &ContractsRepository{db: db}
}

func nullableStringValue(value string) interface{} {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return strings.TrimSpace(value)
}

func nullableUint64Pointer(value *uint64) interface{} {
	if value == nil || *value == 0 {
		return nil
	}
	return *value
}

func nullableTimePointer(value *time.Time) interface{} {
	if value == nil || value.IsZero() {
		return nil
	}
	return *value
}

func (r *ContractsRepository) FetchStatuses() ([]models.ContractStatus, error) {
	query := `
		SELECT id, nome_status
		FROM treehousedb.contratos_status
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []models.ContractStatus
	for rows.Next() {
		var status models.ContractStatus
		if err = rows.Scan(&status.ID, &status.Name); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, rows.Err()
}

func (r *ContractsRepository) FetchTypes() ([]models.ContractType, error) {
	query := `
		SELECT id, nome_tipo
		FROM treehousedb.contratos_tipos
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.ContractType
	for rows.Next() {
		var contractType models.ContractType
		if err = rows.Scan(&contractType.ID, &contractType.Name); err != nil {
			return nil, err
		}
		types = append(types, contractType)
	}

	return types, rows.Err()
}

func scanContract(scanner interface{ Scan(dest ...interface{}) error }, contract *models.Contract) error {
	var classID sql.NullInt64
	var representativeEmail sql.NullString
	var representativeCPF sql.NullString
	var rg sql.NullString
	var representativePhone sql.NullString
	var representativeCivilStatus sql.NullString
	var installments sql.NullInt64
	var installmentsDescription sql.NullString
	var lessonsCount sql.NullInt64
	var periodicity sql.NullString
	var lessonDuration sql.NullString
	var contractDuration sql.NullString
	var startDate sql.NullTime
	var dueDate sql.NullTime
	var firstLessonDate sql.NullTime
	var studentName sql.NullString
	var responsibleName sql.NullString
	var representativeName sql.NullString
	var contractTypeName sql.NullString
	var statusName sql.NullString

	err := scanner.Scan(
		&contract.ID,
		&contract.RepresentativeCustomerID,
		&contract.ResponsibleCustomerID,
		&contract.StudentID,
		&contract.ContractTypeID,
		&contract.StatusID,
		&classID,
		&contract.Value,
		&representativeEmail,
		&representativeCPF,
		&rg,
		&representativePhone,
		&representativeCivilStatus,
		&contract.DiscountPercentage,
		&contract.FinalValue,
		&installments,
		&installmentsDescription,
		&lessonsCount,
		&periodicity,
		&lessonDuration,
		&contractDuration,
		&startDate,
		&dueDate,
		&firstLessonDate,
		&contract.CreatedAt,
		&contract.UpdatedAt,
		&studentName,
		&responsibleName,
		&representativeName,
		&contractTypeName,
		&statusName,
	)
	if err != nil {
		return err
	}

	if classID.Valid {
		value := uint64(classID.Int64)
		contract.ClassID = &value
	}
	if representativeEmail.Valid {
		contract.RepresentativeEmail = representativeEmail.String
	}
	if representativeCPF.Valid {
		contract.RepresentativeCPF = representativeCPF.String
	}
	if rg.Valid {
		contract.RG = rg.String
	}
	if representativePhone.Valid {
		contract.RepresentativePhone = representativePhone.String
	}
	if representativeCivilStatus.Valid {
		contract.RepresentativeCivilStatus = representativeCivilStatus.String
	}
	if installments.Valid {
		value := uint64(installments.Int64)
		contract.Installments = &value
	}
	if installmentsDescription.Valid {
		contract.InstallmentsDescription = installmentsDescription.String
	}
	if lessonsCount.Valid {
		value := uint64(lessonsCount.Int64)
		contract.LessonsCount = &value
	}
	if periodicity.Valid {
		contract.Periodicity = periodicity.String
	}
	if lessonDuration.Valid {
		contract.LessonDuration = lessonDuration.String
	}
	if contractDuration.Valid {
		contract.ContractDuration = contractDuration.String
	}
	if startDate.Valid {
		contract.StartDate = &startDate.Time
	}
	if dueDate.Valid {
		contract.DueDate = &dueDate.Time
	}
	if firstLessonDate.Valid {
		contract.FirstLessonDate = &firstLessonDate.Time
	}
	if studentName.Valid {
		contract.StudentName = studentName.String
	}
	if responsibleName.Valid {
		contract.ResponsibleName = responsibleName.String
	}
	if representativeName.Valid {
		contract.RepresentativeName = representativeName.String
	}
	if contractTypeName.Valid {
		contract.ContractTypeName = contractTypeName.String
	}
	if statusName.Valid {
		contract.StatusName = statusName.String
	}

	contract.EffectiveStatusID, contract.EffectiveStatusName = models.ComputeEffectiveContractStatus(*contract, time.Now())
	return nil
}

func (r *ContractsRepository) baseContractQuery() string {
	return `
		SELECT
			c.id,
			c.id_cliente_representante,
			c.id_cliente_responsavel,
			c.id_aluno,
			c.id_tipo_contrato,
			c.id_status,
			c.id_turma,
			c.valor,
			c.email_representante,
			c.cpf_representante,
			c.rg,
			c.telefone_representante,
			c.est_civil_representante,
			c.desconto_porcentagem,
			c.valor_final,
			c.parcelas,
			c.parcelas_descricao,
			c.numero_aulas,
			c.periodicidade,
			c.tempo_aula,
			c.tempo_contrato,
			c.inicio_contrato,
			c.vencimento_contrato,
			c.primeira_aula,
			c.created_at,
			c.updated_at,
			a.nome AS student_name,
			cr.nome AS responsible_name,
			cp.nome AS representative_name,
			ct.nome_tipo,
			cs.nome_status
		FROM treehousedb.contratos c
		INNER JOIN treehousedb.alunos a
			ON a.id = c.id_aluno
		INNER JOIN treehousedb.clientes cr
			ON cr.id = c.id_cliente_responsavel
		INNER JOIN treehousedb.clientes cp
			ON cp.id = c.id_cliente_representante
		INNER JOIN treehousedb.contratos_tipos ct
			ON ct.id = c.id_tipo_contrato
		INNER JOIN treehousedb.contratos_status cs
			ON cs.id = c.id_status
	`
}

func (r *ContractsRepository) FetchAll(search string) ([]models.Contract, error) {
	query := r.baseContractQuery()
	var args []interface{}

	if search != "" {
		query += `
			WHERE
				LOWER(a.nome) LIKE ?
				OR LOWER(cr.nome) LIKE ?
				OR CAST(c.id AS CHAR) LIKE ?
		`
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query += " ORDER BY c.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contracts []models.Contract
	for rows.Next() {
		var contract models.Contract
		if err = scanContract(rows, &contract); err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}

	return contracts, rows.Err()
}

func (r *ContractsRepository) FetchByID(contractID uint64) (models.Contract, error) {
	query := r.baseContractQuery() + " WHERE c.id = ? LIMIT 1"

	var contract models.Contract
	err := scanContract(r.db.QueryRow(query, contractID), &contract)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Contract{}, sql.ErrNoRows
		}
		return models.Contract{}, err
	}

	return contract, nil
}

func insertCustomerTx(tx *sql.Tx, customer models.Customer) (uint64, error) {
	queryCustomer := `
		INSERT INTO treehousedb.clientes (
			nome,
			email,
			cpf,
			rg,
			telefone,
			ativo,
			nascimento
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		queryCustomer,
		customer.Nome,
		nullableStringValue(customer.Email),
		nullableStringValue(customer.CPF),
		nullableStringValue(customer.RG),
		nullableStringValue(customer.Telefone),
		true,
		nullableTimePointer(customer.Nascimento),
	)
	if err != nil {
		return 0, err
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, address := range customer.Enderecos {
		addressQuery := `
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

		addressResult, addressErr := tx.Exec(
			addressQuery,
			address.CEP,
			address.Rua,
			address.Numero,
			address.Bairro,
			address.Cidade,
			address.Estado,
			address.Pais,
			nullableStringValue(address.Complemento),
		)
		if addressErr != nil {
			return 0, addressErr
		}

		addressID, addressErr := addressResult.LastInsertId()
		if addressErr != nil {
			return 0, addressErr
		}

		if _, addressErr = tx.Exec(
			`INSERT INTO treehousedb.enderecos_clientes (id_cliente, id_endereco) VALUES (?, ?)`,
			customerID,
			addressID,
		); addressErr != nil {
			return 0, addressErr
		}
	}

	return uint64(customerID), nil
}

func insertStudentTx(tx *sql.Tx, student models.Student) (uint64, error) {
	result, err := tx.Exec(
		`INSERT INTO treehousedb.alunos (nome, livro, alfabetizacao, nascimento, ativo) VALUES (?, ?, ?, ?, ?)`,
		student.Nome,
		nullIfEmpty(student.Livro),
		nullIfEmpty(student.Alfabetizacao),
		nullableTimePointer(student.Nascimento),
		true,
	)
	if err != nil {
		return 0, err
	}

	studentID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(studentID), nil
}

func ensureCustomerStudentLinkTx(tx *sql.Tx, customerID, studentID uint64) error {
	_, err := tx.Exec(
		`INSERT INTO treehousedb.clientes_alunos (id_cliente, id_aluno) VALUES (?, ?)
		 ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP`,
		customerID,
		studentID,
	)
	return err
}

func (r *ContractsRepository) resolveContractPartiesTx(tx *sql.Tx, contract *models.Contract) error {
	if contract.ResponsibleCustomerID == 0 {
		if contract.ResponsibleCustomer == nil {
			return errors.New("cliente responsável é obrigatório")
		}
		customerID, err := insertCustomerTx(tx, *contract.ResponsibleCustomer)
		if err != nil {
			return err
		}
		contract.ResponsibleCustomerID = customerID
	}

	if contract.RepresentativeCustomerID == 0 {
		if contract.RepresentativeCustomer != nil {
			customerID, err := insertCustomerTx(tx, *contract.RepresentativeCustomer)
			if err != nil {
				return err
			}
			contract.RepresentativeCustomerID = customerID
		} else {
			contract.RepresentativeCustomerID = contract.ResponsibleCustomerID
		}
	}

	if contract.StudentID == 0 {
		if contract.Student == nil {
			return errors.New("aluno é obrigatório")
		}
		studentID, err := insertStudentTx(tx, *contract.Student)
		if err != nil {
			return err
		}
		contract.StudentID = studentID
	}

	return ensureCustomerStudentLinkTx(tx, contract.ResponsibleCustomerID, contract.StudentID)
}

func (r *ContractsRepository) Insert(contract models.Contract) (uint64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = r.resolveContractPartiesTx(tx, &contract); err != nil {
		return 0, err
	}

	query := `
		INSERT INTO treehousedb.contratos (
			id_cliente_representante,
			id_cliente_responsavel,
			id_aluno,
			id_tipo_contrato,
			id_status,
			id_turma,
			valor,
			email_representante,
			cpf_representante,
			rg,
			telefone_representante,
			est_civil_representante,
			desconto_porcentagem,
			valor_final,
			parcelas,
			parcelas_descricao,
			numero_aulas,
			periodicidade,
			tempo_aula,
			tempo_contrato,
			inicio_contrato,
			vencimento_contrato,
			primeira_aula
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(
		query,
		contract.RepresentativeCustomerID,
		contract.ResponsibleCustomerID,
		contract.StudentID,
		contract.ContractTypeID,
		contract.StatusID,
		nullableUint64Pointer(contract.ClassID),
		contract.Value,
		nullableStringValue(contract.RepresentativeEmail),
		nullableStringValue(contract.RepresentativeCPF),
		nullableStringValue(contract.RG),
		nullableStringValue(contract.RepresentativePhone),
		nullableStringValue(contract.RepresentativeCivilStatus),
		contract.DiscountPercentage,
		contract.FinalValue,
		nullableUint64Pointer(contract.Installments),
		nullableStringValue(contract.InstallmentsDescription),
		nullableUint64Pointer(contract.LessonsCount),
		nullableStringValue(contract.Periodicity),
		nullableStringValue(contract.LessonDuration),
		nullableStringValue(contract.ContractDuration),
		nullableTimePointer(contract.StartDate),
		nullableTimePointer(contract.DueDate),
		nullableTimePointer(contract.FirstLessonDate),
	)
	if err != nil {
		return 0, err
	}

	contractID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(contractID), nil
}

func (r *ContractsRepository) Update(contractID uint64, contract models.Contract) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = r.resolveContractPartiesTx(tx, &contract); err != nil {
		return err
	}

	query := `
		UPDATE treehousedb.contratos
		SET
			id_cliente_representante = ?,
			id_cliente_responsavel = ?,
			id_aluno = ?,
			id_tipo_contrato = ?,
			id_status = ?,
			id_turma = ?,
			valor = ?,
			email_representante = ?,
			cpf_representante = ?,
			rg = ?,
			telefone_representante = ?,
			est_civil_representante = ?,
			desconto_porcentagem = ?,
			valor_final = ?,
			parcelas = ?,
			parcelas_descricao = ?,
			numero_aulas = ?,
			periodicidade = ?,
			tempo_aula = ?,
			tempo_contrato = ?,
			inicio_contrato = ?,
			vencimento_contrato = ?,
			primeira_aula = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := tx.Exec(
		query,
		contract.RepresentativeCustomerID,
		contract.ResponsibleCustomerID,
		contract.StudentID,
		contract.ContractTypeID,
		contract.StatusID,
		nullableUint64Pointer(contract.ClassID),
		contract.Value,
		nullableStringValue(contract.RepresentativeEmail),
		nullableStringValue(contract.RepresentativeCPF),
		nullableStringValue(contract.RG),
		nullableStringValue(contract.RepresentativePhone),
		nullableStringValue(contract.RepresentativeCivilStatus),
		contract.DiscountPercentage,
		contract.FinalValue,
		nullableUint64Pointer(contract.Installments),
		nullableStringValue(contract.InstallmentsDescription),
		nullableUint64Pointer(contract.LessonsCount),
		nullableStringValue(contract.Periodicity),
		nullableStringValue(contract.LessonDuration),
		nullableStringValue(contract.ContractDuration),
		nullableTimePointer(contract.StartDate),
		nullableTimePointer(contract.DueDate),
		nullableTimePointer(contract.FirstLessonDate),
		contractID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("nenhum contrato encontrado com id %d", contractID)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ContractsRepository) CreateClassFromContract(contractID uint64, class models.Class) (uint64, int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var studentID uint64
	var studentName string
	var existingClassID sql.NullInt64

	queryContract := `
		SELECT c.id_aluno, a.nome, c.id_turma
		FROM treehousedb.contratos c
		INNER JOIN treehousedb.alunos a
			ON a.id = c.id_aluno
		WHERE c.id = ?
		LIMIT 1
	`

	if err = tx.QueryRow(queryContract, contractID).Scan(&studentID, &studentName, &existingClassID); err != nil {
		return 0, 0, err
	}

	if existingClassID.Valid {
		return 0, 0, errors.New("o contrato já possui uma turma vinculada")
	}

	class.Name = strings.TrimSpace(class.Name)
	if class.Name == "" {
		class.Name = fmt.Sprintf("Turma %s", strings.TrimSpace(studentName))
	}

	queryClass := `
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

	result, err := tx.Exec(
		queryClass,
		teacherID,
		class.Name,
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

	generatedLessonsCount, err := createLessonsFromRecurrence(tx, uint64(classID), class)
	if err != nil {
		return 0, 0, err
	}

	if _, err = tx.Exec(
		`INSERT INTO treehousedb.alunos_turmas (id_aluno, id_turma) VALUES (?, ?)`,
		studentID,
		classID,
	); err != nil {
		return 0, 0, err
	}

	if _, err = tx.Exec(
		`UPDATE treehousedb.contratos
		 SET id_turma = ?, numero_aulas = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		classID,
		generatedLessonsCount,
		contractID,
	); err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, 0, err
	}

	return uint64(classID), generatedLessonsCount, nil
}
