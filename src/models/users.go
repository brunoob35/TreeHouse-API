package models

import (
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/brunoob35/TreeHouse-API/src/security"
	"github.com/brunoob35/TreeHouse-API/src/utils"
)

type User struct {
	ID         uint64       `json:"id,omitempty"`
	IDEndereco *uint64      `json:"id_endereco,omitempty"`
	Nome       string       `json:"nome"`
	Email      string       `json:"email"`
	Senha      string       `json:"senha"`
	CPF        string       `json:"cpf,omitempty"`
	RG         string       `json:"rg,omitempty"`
	Telefone   string       `json:"telefone,omitempty"`
	Ativo      bool         `json:"ativo"`
	Nascimento *time.Time   `json:"nascimento,omitempty"`
	CreatedAt  time.Time    `json:"created_at,omitempty"`
	UpdatedAt  time.Time    `json:"updated_at,omitempty"`
	Endereco   *Address     `json:"endereco,omitempty"`
	Permissoes []Permission `json:"permissoes,omitempty"`
}

type Permission struct {
	ID   uint64 `json:"id"`
	Nome string `json:"nome,omitempty"`
}

// Prepare treats user info and validates it.
func (user *User) Prepare(step string) error {
	if err := user.format(step); err != nil {
		return err
	}

	if err := user.validate(step); err != nil {
		return err
	}

	return nil
}

// validate prevents the system from saving invalid info.
func (user *User) validate(step string) error {
	if user.Nome == "" {
		return errors.New("o nome é obrigatório")
	}

	if user.Email == "" {
		return errors.New("o e-mail é obrigatório")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("o e-mail inserido é inválido")
	}

	if user.CPF != "" {
		if err := utils.CPFValidator(user.CPF); err != nil {
			return err
		}
	}

	if step == "create" && user.Senha == "" {
		return errors.New("a senha é obrigatória")
	}

	return nil
}

func (user *User) format(step string) error {
	user.Nome = strings.TrimSpace(user.Nome)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.CPF = strings.TrimSpace(user.CPF)
	user.RG = strings.TrimSpace(user.RG)
	user.Telefone = strings.TrimSpace(user.Telefone)

	if user.Endereco != nil {
		user.Endereco.Rua = strings.TrimSpace(user.Endereco.Rua)
		user.Endereco.Numero = strings.TrimSpace(user.Endereco.Numero)
		user.Endereco.Bairro = strings.TrimSpace(user.Endereco.Bairro)
		user.Endereco.Cidade = strings.TrimSpace(user.Endereco.Cidade)
		user.Endereco.Estado = strings.TrimSpace(user.Endereco.Estado)
		user.Endereco.Pais = strings.TrimSpace(user.Endereco.Pais)
		user.Endereco.Complemento = strings.TrimSpace(user.Endereco.Complemento)

		if user.Endereco.Pais == "" {
			user.Endereco.Pais = "Brasil"
		}
	}

	if step == "create" && user.Senha != "" {
		hashedPassword, err := security.Hash(user.Senha)
		if err != nil {
			return err
		}
		user.Senha = string(hashedPassword)
	}

	return nil
}
