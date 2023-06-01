package model

import "time"

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
