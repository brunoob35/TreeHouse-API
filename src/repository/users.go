package repository

import (
	"database/sql"
	"github.com/brunoob35/TreeHouse-API/src/model"
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
func (u Users) Create(user model.User) (uint64, error) {
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

func (u Users) Fetch(nick string) ([]model.User, error) {
	//nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) // %nomeOuNick%
	//
	//linhas, erro := repositorio.db.Query(
	//	"select id, nome, nick, email, criadoEm from usuarios where nome LIKE ? or nick LIKE ?",
	//	nomeOuNick, nomeOuNick,
	//)
	//
	//if erro != nil {
	//	return nil, erro
	//}
	//defer linhas.Close()
	//
	//var usuarios []modelos.Usuario
	//
	//for linhas.Next() {
	//	var usuario modelos.Usuario
	//
	//	if erro = linhas.Scan(
	//		&usuario.ID,
	//		&usuario.Nome,
	//		&usuario.Nick,
	//		&usuario.Email,
	//		&usuario.CriadoEm,
	//	); erro != nil {
	//		return nil, erro
	//	}
	//
	//	usuarios = append(usuarios, usuario)
	//}
	//
	//return usuarios, nil

	return nil, nil
}

func (u Users) FetchByID(ID uint64) (model.User, error) {
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
	return model.User{}, nil
}

// FetchByEmail Fetches a user by email
func (u Users) FetchByEmail(email string) (model.User, error) {
	query := `SELECT 
    			id_usuario, 
    			senha
			FROM usuarios WHERE email_usuario = ?`

	line, err := u.db.Query(query, email)
	if err != nil {
		return model.User{}, err
	}
	defer line.Close()

	var user model.User

	if line.Next() {
		if err = line.Scan(&user.ID, &user.Senha); err != nil {
			return model.User{}, err
		}
	}

	return user, nil
}
