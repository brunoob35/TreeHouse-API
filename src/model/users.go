package model

import (
	"errors"
	"github.com/badoux/checkmail"
	"strings"
	"time"
)

// User represents any user in the sistem, most user fields are in portuguese since the DB is also in portugues
type User struct {
	ID          uint64    `json:"id,omitempty"`
	Nome        string    `json:"nome_usuario,omitempty"`
	Email       string    `json:"email_usuario,omitempty"`
	Senha       string    `json:"senha,omitempty"`
	IDAcesso    uint64    `json:"id_acesso,omitempty"`
	CPF         string    `json:"cpf,omitempty"`
	RG          string    `json:"rg,omitempty"`
	Celular     string    `json:"celular,omitempty"`
	DataNasc    time.Time `json:"data_nascimento,omitempty"`
	DataCriacao time.Time `json:"data_criacao,omitempty"`
}

// Prepare Treats user info and also validate it
func (user *User) Prepare(step string) error {
	//optimize: test if validate > format is a good sequence since we need to also sanitize RG and CPF.
	if err := user.validate(step); err != nil {
		return err
	}

	if err := user.format(step); err != nil {
		return err
	}

	return nil
}

// validate Prevents the sistem to save any blank space or invalid info.
func (user *User) validate(step string) error {
	if user.Nome == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}

	if user.Email == "" {
		return errors.New("O e-mail é obrigatório e não pode estar em branco")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("O e-mail inserido é inválido")
	}
	//TODO: find a CPF/RG validator service or algorithm

	if step == "cadastro" && user.Senha == "" {
		return errors.New("A senha é obrigatória e não pode estar em branco")
	}

	return nil
}

// format Formats and trims blank spaces. Also applies hashing to the password.
func (user *User) format(step string) error {
	user.Nome = strings.TrimSpace(user.Nome)
	user.Email = strings.TrimSpace(user.Email)

	//if step == "cadastro" {
	//	hashedPAssword, erro := seguranca.Hash(user.Senha)
	//	if erro != nil {
	//		return erro
	//	}
	//
	//	user.Senha = string(hashedPAssword)
	//}

	return nil
}
