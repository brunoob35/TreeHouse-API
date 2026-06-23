package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/security"
)

type CustomersRepository struct {
	db          *sql.DB
	auditUserID *uint64
}

func NewCustomersRepository(db *sql.DB) *CustomersRepository {
	return &CustomersRepository{db: db}
}

func (r *CustomersRepository) WithAuditUser(userID uint64) *CustomersRepository {
	r.auditUserID = &userID
	return r
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}

	return value
}

func encryptCustomerPII(value string) (interface{}, interface{}, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil, nil
	}

	encrypted, err := security.EncryptPII(trimmed)
	if err != nil {
		return nil, nil, err
	}

	hash, err := security.HashPII(trimmed)
	if err != nil {
		return nil, nil, err
	}

	return encrypted, hash, nil
}

func (r *CustomersRepository) Insert(customer models.Customer) (uint64, error) {
	cpfProtected, cpfHash, err := encryptCustomerPII(customer.CPF)
	if err != nil {
		return 0, err
	}

	rgProtected, rgHash, err := encryptCustomerPII(customer.RG)
	if err != nil {
		return 0, err
	}

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

	queryCustomer := `
		INSERT INTO treehousedb.clientes (
			nome,
			cpf,
			email,
			telefone,
			rg,
			cpf_protegido,
			rg_protegido,
			cpf_hash,
			rg_hash,
			nascimento,
			ativo,
			lgpd_aceito,
			lgpd_aceito_em,
			lgpd_finalidade
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	customerResult, err := tx.Exec(
		queryCustomer,
		customer.Nome,
		nil,
		nullableString(customer.Email),
		customer.Telefone,
		nil,
		cpfProtected,
		rgProtected,
		cpfHash,
		rgHash,
		customer.Nascimento,
		customer.Ativo,
		customer.LGPDAceito,
		customer.LGPDAceitoEm,
		nullableString(customer.LGPDFinalidade),
	)
	if err != nil {
		return 0, err
	}

	customerID, err := customerResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	if len(customer.Enderecos) > 0 {
		queryAddress := `
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

		queryAddressLink := `
			INSERT INTO treehousedb.enderecos_clientes (
				id_cliente,
				id_endereco
			) VALUES (?, ?)
		`

		for _, address := range customer.Enderecos {
			addressResult, addressErr := tx.Exec(
				queryAddress,
				address.CEP,
				address.Rua,
				address.Numero,
				address.Bairro,
				address.Cidade,
				address.Estado,
				address.Pais,
				address.Complemento,
			)
			if addressErr != nil {
				err = addressErr
				return 0, err
			}

			addressID, addressErr := addressResult.LastInsertId()
			if addressErr != nil {
				err = addressErr
				return 0, err
			}

			if _, addressErr = tx.Exec(queryAddressLink, customerID, addressID); addressErr != nil {
				err = addressErr
				return 0, err
			}
		}
	}

	if len(customer.Students) > 0 {
		queryStudent := `
			INSERT INTO treehousedb.alunos (
				nome,
				livro,
				alfabetizacao,
				nascimento,
				ativo
			) VALUES (?, ?, ?, ?, ?)
		`

		queryLink := `
			INSERT INTO treehousedb.clientes_alunos (
				id_cliente,
				id_aluno
			) VALUES (?, ?)
		`

		for _, student := range customer.Students {
			studentID := int64(student.ID)

			if studentID == 0 {
				studentResult, studentErr := tx.Exec(
					queryStudent,
					student.Nome,
					student.Livro,
					student.Alfabetizacao,
					student.Nascimento,
					true,
				)
				if studentErr != nil {
					err = studentErr
					return 0, err
				}

				studentID, studentErr = studentResult.LastInsertId()
				if studentErr != nil {
					err = studentErr
					return 0, err
				}
			}

			if _, studentErr := tx.Exec(queryLink, customerID, studentID); studentErr != nil {
				err = studentErr
				return 0, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(customerID), nil
}

func (r *CustomersRepository) FetchAll(search string) ([]models.Customer, error) {
	query := `
		SELECT
			c.id,
			c.nome,
			COALESCE(c.cpf, ''),
			COALESCE(c.email, ''),
			COALESCE(c.telefone, ''),
			COALESCE(c.rg, ''),
			c.cpf_protegido,
			c.rg_protegido,
			c.nascimento,
			c.ativo,
			c.lgpd_aceito,
			c.lgpd_aceito_em,
			COALESCE(c.lgpd_finalidade, ''),
			COUNT(DISTINCT ca.id_aluno) AS students_count,
			0 AS contracts_count,
			c.created_at,
			c.updated_at
		FROM treehousedb.clientes c
		LEFT JOIN treehousedb.clientes_alunos ca
			ON ca.id_cliente = c.id
	`

	var args []interface{}
	if search != "" {
		cpfSearchHash := ""
		digitsOnlySearch := strings.NewReplacer(".", "", "-", "", " ", "", "(", "", ")", "", "+", "").Replace(search)
		if digitsOnlySearch != "" {
			cpfSearchHash, _ = security.HashPII(digitsOnlySearch)
		}
		query += `
			WHERE
				LOWER(c.nome) LIKE ?
				OR LOWER(COALESCE(c.email, '')) LIKE ?
				OR REPLACE(COALESCE(c.cpf, ''), '.', '') LIKE REPLACE(?, '.', '')
				OR (? <> '' AND c.cpf_hash = ?)
		`
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%", cpfSearchHash, cpfSearchHash)
	}

	query += `
		GROUP BY
			c.id,
			c.nome,
			c.cpf,
			c.email,
			c.telefone,
			c.rg,
			c.cpf_protegido,
			c.rg_protegido,
			c.nascimento,
			c.ativo,
			c.lgpd_aceito,
			c.lgpd_aceito_em,
			c.lgpd_finalidade,
			c.created_at,
			c.updated_at
		ORDER BY c.nome
	`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		var cpfLegacy string
		var rgLegacy string
		var cpfProtected sql.NullString
		var rgProtected sql.NullString
		var nascimento sql.NullTime
		var lgpdAcceptedAt sql.NullTime
		var lgpdPurpose string

		err = rows.Scan(
			&customer.ID,
			&customer.Nome,
			&cpfLegacy,
			&customer.Email,
			&customer.Telefone,
			&rgLegacy,
			&cpfProtected,
			&rgProtected,
			&nascimento,
			&customer.Ativo,
			&customer.LGPDAceito,
			&lgpdAcceptedAt,
			&lgpdPurpose,
			&customer.StudentsCount,
			&customer.ContractsCount,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if nascimento.Valid {
			customer.Nascimento = &nascimento.Time
		}
		if cpfProtected.Valid {
			customer.CPF, err = security.DecryptPII(cpfProtected.String)
			if err != nil {
				return nil, err
			}
		} else {
			customer.CPF = cpfLegacy
		}
		if rgProtected.Valid {
			customer.RG, err = security.DecryptPII(rgProtected.String)
			if err != nil {
				return nil, err
			}
		} else {
			customer.RG = rgLegacy
		}
		if lgpdAcceptedAt.Valid {
			customer.LGPDAceitoEm = &lgpdAcceptedAt.Time
		}
		customer.LGPDFinalidade = lgpdPurpose

		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomersRepository) FetchByID(id uint64) (models.Customer, error) {
	query := `
		SELECT
			c.id,
			c.nome,
			COALESCE(c.cpf, ''),
			COALESCE(c.email, ''),
			COALESCE(c.telefone, ''),
			COALESCE(c.rg, ''),
			c.cpf_protegido,
			c.rg_protegido,
			c.nascimento,
			c.ativo,
			c.lgpd_aceito,
			c.lgpd_aceito_em,
			COALESCE(c.lgpd_finalidade, ''),
			COUNT(DISTINCT ca.id_aluno) AS students_count,
			0 AS contracts_count,
			c.created_at,
			c.updated_at
		FROM treehousedb.clientes c
		LEFT JOIN treehousedb.clientes_alunos ca
			ON ca.id_cliente = c.id
		WHERE c.id = ?
		GROUP BY
			c.id,
			c.nome,
			c.cpf,
			c.email,
			c.telefone,
			c.rg,
			c.cpf_protegido,
			c.rg_protegido,
			c.nascimento,
			c.ativo,
			c.lgpd_aceito,
			c.lgpd_aceito_em,
			c.lgpd_finalidade,
			c.created_at,
			c.updated_at
	`

	var customer models.Customer
	var cpfLegacy string
	var rgLegacy string
	var cpfProtected sql.NullString
	var rgProtected sql.NullString
	var nascimento sql.NullTime
	var lgpdAcceptedAt sql.NullTime
	var lgpdPurpose string
	err := r.db.QueryRow(query, id).Scan(
		&customer.ID,
		&customer.Nome,
		&cpfLegacy,
		&customer.Email,
		&customer.Telefone,
		&rgLegacy,
		&cpfProtected,
		&rgProtected,
		&nascimento,
		&customer.Ativo,
		&customer.LGPDAceito,
		&lgpdAcceptedAt,
		&lgpdPurpose,
		&customer.StudentsCount,
		&customer.ContractsCount,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Customer{}, sql.ErrNoRows
		}
		return models.Customer{}, err
	}

	if nascimento.Valid {
		customer.Nascimento = &nascimento.Time
	}
	if cpfProtected.Valid {
		customer.CPF, err = security.DecryptPII(cpfProtected.String)
		if err != nil {
			return models.Customer{}, err
		}
	} else {
		customer.CPF = cpfLegacy
	}
	if rgProtected.Valid {
		customer.RG, err = security.DecryptPII(rgProtected.String)
		if err != nil {
			return models.Customer{}, err
		}
	} else {
		customer.RG = rgLegacy
	}
	if lgpdAcceptedAt.Valid {
		customer.LGPDAceitoEm = &lgpdAcceptedAt.Time
	}
	customer.LGPDFinalidade = lgpdPurpose

	return customer, nil
}

func (r *CustomersRepository) Update(id uint64, customer models.Customer) error {
	cpfProtected, cpfHash, err := encryptCustomerPII(customer.CPF)
	if err != nil {
		return err
	}

	rgProtected, rgHash, err := encryptCustomerPII(customer.RG)
	if err != nil {
		return err
	}

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
		UPDATE treehousedb.clientes
		SET
			nome = ?,
			cpf = ?,
			email = ?,
			telefone = ?,
			rg = ?,
			cpf_protegido = ?,
			rg_protegido = ?,
			cpf_hash = ?,
			rg_hash = ?,
			nascimento = ?,
			ativo = ?,
			lgpd_aceito = ?,
			lgpd_aceito_em = ?,
			lgpd_finalidade = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := tx.Exec(
		query,
		customer.Nome,
		nil,
		nullableString(customer.Email),
		customer.Telefone,
		nil,
		cpfProtected,
		rgProtected,
		cpfHash,
		rgHash,
		customer.Nascimento,
		customer.Ativo,
		customer.LGPDAceito,
		customer.LGPDAceitoEm,
		nullableString(customer.LGPDFinalidade),
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
		return fmt.Errorf("nenhum cliente encontrado com id %d", id)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CustomersRepository) ReplaceAddresses(customerID uint64, addresses []models.Address) error {
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

	deleteLinksQuery := `
		DELETE FROM treehousedb.enderecos_clientes
		WHERE id_cliente = ?
	`

	if _, err = tx.Exec(deleteLinksQuery, customerID); err != nil {
		return err
	}

	if len(addresses) > 0 {
		insertAddressQuery := `
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

		insertLinkQuery := `
			INSERT INTO treehousedb.enderecos_clientes (
				id_cliente,
				id_endereco
			) VALUES (?, ?)
		`

		for _, address := range addresses {
			result, insertErr := tx.Exec(
				insertAddressQuery,
				address.CEP,
				address.Rua,
				address.Numero,
				address.Bairro,
				address.Cidade,
				address.Estado,
				address.Pais,
				address.Complemento,
			)
			if insertErr != nil {
				err = insertErr
				return err
			}

			addressID, insertErr := result.LastInsertId()
			if insertErr != nil {
				err = insertErr
				return err
			}

			if _, insertErr = tx.Exec(insertLinkQuery, customerID, addressID); insertErr != nil {
				err = insertErr
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CustomersRepository) ReplaceStudents(customerID uint64, students []models.Student) error {
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

	deleteLinksQuery := `
		DELETE FROM treehousedb.clientes_alunos
		WHERE id_cliente = ?
	`

	if _, err = tx.Exec(deleteLinksQuery, customerID); err != nil {
		return err
	}

	if len(students) > 0 {
		insertStudentQuery := `
			INSERT INTO treehousedb.alunos (
				nome,
				livro,
				alfabetizacao,
				nascimento,
				ativo
			) VALUES (?, ?, ?, ?, ?)
		`

		insertLinkQuery := `
			INSERT INTO treehousedb.clientes_alunos (
				id_cliente,
				id_aluno
			) VALUES (?, ?)
		`

		for _, student := range students {
			studentID := int64(student.ID)

			if studentID == 0 {
				result, insertErr := tx.Exec(
					insertStudentQuery,
					student.Nome,
					student.Livro,
					student.Alfabetizacao,
					student.Nascimento,
					true,
				)
				if insertErr != nil {
					err = insertErr
					return err
				}

				studentID, insertErr = result.LastInsertId()
				if insertErr != nil {
					err = insertErr
					return err
				}
			}

			if _, insertErr := tx.Exec(insertLinkQuery, customerID, studentID); insertErr != nil {
				err = insertErr
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CustomersRepository) SoftDelete(id uint64) error {
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
		UPDATE treehousedb.clientes
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
		return fmt.Errorf("nenhum cliente encontrado com id %d", id)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CustomersRepository) FetchStudents(customerID uint64) ([]models.Student, error) {
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
		INNER JOIN treehousedb.clientes_alunos ca
			ON ca.id_aluno = a.id
		WHERE ca.id_cliente = ?
		ORDER BY a.nome
	`

	rows, err := r.db.Query(query, customerID)
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

func (r *CustomersRepository) FetchAddresses(customerID uint64) ([]models.Address, error) {
	query := `
		SELECT
			e.id,
			COALESCE(e.cep, ''),
			e.rua,
			e.numero,
			e.bairro,
			e.cidade,
			e.estado,
			COALESCE(e.pais, 'Brasil'),
			COALESCE(e.complemento, ''),
			e.created_at,
			e.updated_at
		FROM treehousedb.enderecos e
		INNER JOIN treehousedb.enderecos_clientes ec
			ON ec.id_endereco = e.id
		WHERE ec.id_cliente = ?
		ORDER BY e.id ASC
	`

	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.Address
	for rows.Next() {
		var address models.Address
		var rua sql.NullString
		var numero sql.NullString
		var bairro sql.NullString
		var cidade sql.NullString
		var estado sql.NullString
		err = rows.Scan(
			&address.ID,
			&address.CEP,
			&rua,
			&numero,
			&bairro,
			&cidade,
			&estado,
			&address.Pais,
			&address.Complemento,
			&address.CreatedAt,
			&address.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if rua.Valid {
			address.Rua = rua.String
		}
		if numero.Valid {
			address.Numero = numero.String
		}
		if bairro.Valid {
			address.Bairro = bairro.String
		}
		if cidade.Valid {
			address.Cidade = cidade.String
		}
		if estado.Valid {
			address.Estado = estado.String
		}

		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}
