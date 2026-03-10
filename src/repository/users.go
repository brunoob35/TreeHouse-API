package repository

import (
	"database/sql"
	"fmt"

	"github.com/brunoob35/TreeHouse-API/src/models"
)

type Users struct {
	db *sql.DB
}

func UsersNewRepo(db *sql.DB) *Users {
	return &Users{db}
}

func (u Users) Create(user models.User) (uint64, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	addressQuery := `
		INSERT INTO treehousedb.enderecos (
			rua,
			numero,
			bairro,
			cidade,
			estado,
			pais,
			complemento
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	addressStmt, err := tx.Prepare(addressQuery)
	if err != nil {
		return 0, err
	}
	defer addressStmt.Close()

	addressResult, err := addressStmt.Exec(
		user.Endereco.Rua,
		user.Endereco.Numero,
		user.Endereco.Bairro,
		user.Endereco.Cidade,
		user.Endereco.Estado,
		user.Endereco.Pais,
		nullIfEmpty(user.Endereco.Complemento),
	)
	if err != nil {
		return 0, err
	}

	addressID, err := addressResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	userQuery := `
		INSERT INTO treehousedb.usuarios (
			id_endereco,
			senha,
			nome,
			email,
			cpf,
			rg,
			telefone,
			ativo,
			nascimento
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	userStmt, err := tx.Prepare(userQuery)
	if err != nil {
		return 0, err
	}
	defer userStmt.Close()

	userResult, err := userStmt.Exec(
		addressID,
		user.Senha,
		user.Nome,
		user.Email,
		nullIfEmpty(user.CPF),
		nullIfEmpty(user.RG),
		nullIfEmpty(user.Telefone),
		user.Ativo,
		user.Nascimento,
	)
	if err != nil {
		return 0, err
	}

	userID, err := userResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	if len(user.Permissoes) > 0 {
		permQuery := `
			INSERT INTO treehousedb.usuarios_permissoes (
				id_usuario,
				id_permissao
			) VALUES (?, ?)`

		permStmt, err := tx.Prepare(permQuery)
		if err != nil {
			return 0, err
		}
		defer permStmt.Close()

		for _, permissao := range user.Permissoes {
			if permissao.ID == 0 {
				return 0, fmt.Errorf("permissão inválida")
			}

			if _, err = permStmt.Exec(userID, permissao.ID); err != nil {
				return 0, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(userID), nil
}

func (u Users) FetchAllUsers(nome string) ([]models.User, error) {
	query := `
		SELECT
			u.id,
			u.id_endereco,
			u.nome,
			u.email,
			u.cpf,
			u.rg,
			u.telefone,
			u.ativo,
			u.nascimento,
			u.created_at,
			u.updated_at
		FROM treehousedb.usuarios u`

	var rows *sql.Rows
	var err error

	if nome != "" {
		query += ` WHERE LOWER(u.nome) LIKE ? ORDER BY u.nome ASC`
		rows, err = u.db.Query(query, "%"+nome+"%")
	} else {
		query += ` ORDER BY u.nome ASC`
		rows, err = u.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&user.ID,
			&user.IDEndereco,
			&user.Nome,
			&user.Email,
			&user.CPF,
			&user.RG,
			&user.Telefone,
			&user.Ativo,
			&user.Nascimento,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u Users) FetchByID(id uint64) (models.User, error) {
	query := `
		SELECT
			u.id,
			u.id_endereco,
			u.nome,
			u.email,
			u.cpf,
			u.rg,
			u.telefone,
			u.ativo,
			u.nascimento,
			u.created_at,
			u.updated_at,
			e.id,
			e.rua,
			e.numero,
			e.bairro,
			e.cidade,
			e.estado,
			e.pais,
			e.complemento,
			e.created_at,
			e.updated_at
		FROM treehousedb.usuarios u
		INNER JOIN treehousedb.enderecos e ON e.id = u.id_endereco
		WHERE u.id = ?`

	var user models.User

	err := u.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.IDEndereco,
		&user.Nome,
		&user.Email,
		&user.CPF,
		&user.RG,
		&user.Telefone,
		&user.Ativo,
		&user.Nascimento,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Endereco.ID,
		&user.Endereco.Rua,
		&user.Endereco.Numero,
		&user.Endereco.Bairro,
		&user.Endereco.Cidade,
		&user.Endereco.Estado,
		&user.Endereco.Pais,
		&user.Endereco.Complemento,
		&user.Endereco.CreatedAt,
		&user.Endereco.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	permissoes, err := u.fetchPermissionsByUserID(id)
	if err != nil {
		return models.User{}, err
	}
	user.Permissoes = permissoes

	return user, nil
}

func (u Users) FetchByEmail(email string) (models.User, error) {
	query := `
		SELECT
			id,
			senha
		FROM treehousedb.usuarios
		WHERE email = ?`

	var user models.User

	err := u.db.QueryRow(query, email).Scan(&user.ID, &user.Senha)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u Users) fetchPermissionsByUserID(userID uint64) ([]models.Permission, error) {
	query := `
		SELECT
			p.id,
			p.nome
		FROM treehousedb.usuarios_permissoes up
		INNER JOIN treehousedb.permissoes p ON p.id = up.id_permissao
		WHERE up.id_usuario = ?
		ORDER BY p.id`

	rows, err := u.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissoes []models.Permission

	for rows.Next() {
		var p models.Permission
		if err = rows.Scan(&p.ID, &p.Nome); err != nil {
			return nil, err
		}
		permissoes = append(permissoes, p)
	}

	return permissoes, nil
}

func nullIfEmpty(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}
