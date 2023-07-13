package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/brunoob35/TreeHouse-API/src/security"
	"github.com/brunoob35/TreeHouse-API/src/utils"
	"log"
	"strings"
	"time"
)

// User represents any user in the sistem, most user fields are in portuguese since the DB is also in portugues
type User struct {
	ID       uint64    `json:"id,omitempty"`
	Name     string    `json:"nome_usuario,omitempty"`
	Email    string    `json:"email_usuario,omitempty"`
	Password string    `json:"senha,omitempty"`
	IDAccess uint64    `json:"id_acesso,omitempty"`
	CPF      string    `json:"cpf,omitempty"`
	RG       string    `json:"rg,omitempty"`
	Phone    string    `json:"celular,omitempty"`
	Birth    time.Time `json:"data_nascimento,omitempty"`
	Creation time.Time `json:"data_criacao,omitempty"`
}

// Prepare Treats user info and validates it
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
	log.Println("DEBUG: Entrou no validate")

	if user.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}

	if user.Email == "" {
		return errors.New("O e-mail é obrigatório e não pode estar em branco")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("O e-mail inserido é inválido")
	}

	//TODO: find a RG validator service or algorithm
	if err := utils.CPFValidator(user.CPF); err != nil {
		return err
	}

	if step == "cadastro" && user.Password == "" {
		return errors.New("A senha é obrigatória e não pode estar em branco")
	}

	return nil
}

// format Formats and trims blank spaces. Also applies hashing to the password.
func (user *User) format(step string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if step == "register" {
		hashedPassword, err := security.Hash(user.Password)
		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)
	}

	return nil
}
