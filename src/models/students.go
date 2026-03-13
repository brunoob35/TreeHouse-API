package models

import (
	"errors"
	"strings"
	"time"
)

// Student represents a student entity in the system.
// Each student may optionally have a book assigned and
// a literacy level registered.
type Student struct {
	ID            uint64     `json:"id,omitempty"`
	Nome          string     `json:"nome"`
	Livro         string     `json:"livro,omitempty"`
	Alfabetizacao string     `json:"alfabetizacao,omitempty"`
	Nascimento    *time.Time `json:"nascimento,omitempty"`
	Ativo         bool       `json:"ativo"`
	CreatedAt     time.Time  `json:"created_at,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at,omitempty"`
}

// Prepare formats and validates student data before persistence.
// The "step" parameter defines the context of the operation:
//   - "create" → preparing a new record insertion
//   - "update" → preparing a record update
//
// The function first normalizes the data and then validates it.
func (student *Student) Prepare(step string) error {
	student.format()

	if err := student.validate(step); err != nil {
		return err
	}

	return nil
}

// validate ensures that required student data is present.
// Currently the only mandatory field is "Nome".
func (student *Student) validate(step string) error {
	if student.Nome == "" {
		return errors.New("o nome é obrigatório")
	}

	return nil
}

// format normalizes incoming student fields.
// This function trims whitespace from all string fields
// to ensure consistent data storage.
func (student *Student) format() {
	student.Nome = strings.TrimSpace(student.Nome)
	student.Livro = strings.TrimSpace(student.Livro)
	student.Alfabetizacao = strings.TrimSpace(student.Alfabetizacao)
}
