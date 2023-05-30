package repository

import (
	"database/sql"
	"github.com/brunoob35/TreeHouse-API/src/model"
	"log"
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
	log.Println("Chegou no repo")
	query := `INSERT INTO usuarios
				(nome_usuario,
				 email_ususario,
				 senha,
				 id_acesso,
				 cpf,
				 rg,
				 celular,
				 data_nascimento)
				VALUES (?,?,?,?,?,?,?,?)`

	statement, err := u.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Nome, user.Email, user.Senha, user.IDAcesso, user.CPF, user.RG, user.Celular, user.DataNasc)
	if err != nil {
		return 0, err
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertedID), nil
}
