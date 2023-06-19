package repository

import (
	"database/sql"
	"github.com/brunoob35/TreeHouse-API/src/models"
)

// Users receives the DB connection and handles it
// it is representing the repository
type Users struct {
	db *sql.DB
}

// UsersNewRepo receives the DB connection and injects the dependency of it into the repo thought the struct
func UsersNewRepo(db *sql.DB) *Users {
	return &Users{db}
}

// Create inserts and user into the DB and returns the new user ID
func (u Users) Create(user models.User) (uint64, error) {
	query := `INSERT INTO usuarios
				(nome_usuario,
				 email_usuario,
				 senha,
				 id_acesso,
				 cpf,
				 rg,
				 celular)
				VALUES (?,?,?,?,?,?,?)`

	statement, err := u.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Nome, user.Email, user.Senha, user.IDAcesso, user.CPF, user.RG, user.Celular)
	if err != nil {
		return 0, err
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertedID), nil
}

func (u Users) FetchAllUsers(nick string) ([]models.User, error) {
	query := `SELECT * FROM ususarios`

	lines, err := u.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer lines.Close()

	var users []models.User

	for lines.Next() {
		var user models.User

		if err = lines.Scan(
			&user.ID,
			&user.Email,
			&user.Nome,
			&user.IDAcesso,
		); err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil
}

func (u Users) FetchByID(ID uint64) (models.User, error) {
	//linhas, erro := repositorio.db.Query(
	//	"select id, nome, nick, email, criadoEm from usuarios where id = ?",
	//	ID,
	//)
	//if erro != nil {
	//	return modelos.Usuario{}, erro
	//}
	//defer linhas.Close()
	//
	//var usuario modelos.Usuario
	//
	//if linhas.Next() {
	//	if erro = linhas.Scan(
	//		&usuario.ID,
	//		&usuario.Nome,
	//		&usuario.Nick,
	//		&usuario.Email,
	//		&usuario.CriadoEm,
	//	); erro != nil {
	//		return modelos.Usuario{}, erro
	//	}
	//}
	//
	//return usuario, nil
	return models.User{}, nil
}

// FetchByEmail Fetches a user by email
func (u Users) FetchByEmail(email string) (models.User, error) {
	query := `SELECT 
    			id_usuario, 
    			senha
			FROM usuarios WHERE email_usuario = ?`

	line, err := u.db.Query(query, email)
	if err != nil {
		return models.User{}, err
	}
	defer line.Close()

	var user models.User

	if line.Next() {
		if err = line.Scan(&user.ID, &user.Senha); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}
