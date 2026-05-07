package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/brunoob35/TreeHouse-API/src/models"
)

type CustomersRepository struct {
	db *sql.DB
}

func NewCustomersRepository(db *sql.DB) *CustomersRepository {
	return &CustomersRepository{db: db}
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}

	return value
}

func (r *CustomersRepository) Insert(customer models.Customer) (uint64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	queryCustomer := `
		INSERT INTO treehousedb.clientes (
			nome,
			cpf,
			email,
			telefone,
			rg,
			nascimento,
			ativo
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	customerResult, err := tx.Exec(
		queryCustomer,
		customer.Nome,
		customer.CPF,
		nullableString(customer.Email),
		customer.Telefone,
		nullableString(customer.RG),
		customer.Nascimento,
		customer.Ativo,
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

			studentID, studentErr := studentResult.LastInsertId()
			if studentErr != nil {
				err = studentErr
				return 0, err
			}

			if _, studentErr = tx.Exec(queryLink, customerID, studentID); studentErr != nil {
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
			c.nascimento,
			c.ativo,
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
		query += `
			WHERE
				LOWER(c.nome) LIKE ?
				OR LOWER(COALESCE(c.email, '')) LIKE ?
				OR REPLACE(COALESCE(c.cpf, ''), '.', '') LIKE REPLACE(?, '.', '')
		`
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query += `
		GROUP BY
			c.id,
			c.nome,
			c.cpf,
			c.email,
			c.telefone,
			c.rg,
			c.nascimento,
			c.ativo,
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

		err = rows.Scan(
			&customer.ID,
			&customer.Nome,
			&customer.CPF,
			&customer.Email,
			&customer.Telefone,
			&customer.RG,
			&customer.Nascimento,
			&customer.Ativo,
			&customer.StudentsCount,
			&customer.ContractsCount,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

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
			c.nascimento,
			c.ativo,
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
			c.nascimento,
			c.ativo,
			c.created_at,
			c.updated_at
	`

	var customer models.Customer
	err := r.db.QueryRow(query, id).Scan(
		&customer.ID,
		&customer.Nome,
		&customer.CPF,
		&customer.Email,
		&customer.Telefone,
		&customer.RG,
		&customer.Nascimento,
		&customer.Ativo,
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

	return customer, nil
}

func (r *CustomersRepository) Update(id uint64, customer models.Customer) error {
	query := `
		UPDATE treehousedb.clientes
		SET
			nome = ?,
			cpf = ?,
			email = ?,
			telefone = ?,
			rg = ?,
			nascimento = ?,
			ativo = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.Exec(
		query,
		customer.Nome,
		customer.CPF,
		nullableString(customer.Email),
		customer.Telefone,
		nullableString(customer.RG),
		customer.Nascimento,
		customer.Ativo,
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
	query := `
		UPDATE treehousedb.clientes
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
		return fmt.Errorf("nenhum cliente encontrado com id %d", id)
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
		err = rows.Scan(
			&address.ID,
			&address.CEP,
			&address.Rua,
			&address.Numero,
			&address.Bairro,
			&address.Cidade,
			&address.Estado,
			&address.Pais,
			&address.Complemento,
			&address.CreatedAt,
			&address.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}
