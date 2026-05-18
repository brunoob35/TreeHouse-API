package models

import (
	"errors"
	"strings"
	"time"
)

type Class struct {
	ID               uint64     `json:"id,omitempty"`
	TeacherID        *uint64    `json:"teacher_id,omitempty"`
	IDEndereco       *uint64    `json:"id_endereco,omitempty"`
	Name             string     `json:"name"`
	RecurrenceDesc   string     `json:"recurrence_desc,omitempty"`
	RecurrenceJSON   string     `json:"recurrence_json,omitempty"`
	StudentCount     uint64     `json:"student_count,omitempty"`
	LessonsTotal     uint64     `json:"lessons_total,omitempty"`
	LessonsCompleted uint64     `json:"lessons_completed,omitempty"`
	CreatedAt        time.Time  `json:"created_at,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`

	Endereco *Address  `json:"endereco,omitempty"`
	Students []Student `json:"students,omitempty"`
}

type CreatePrivateClassRequest struct {
	StudentID      uint64   `json:"student_id"`
	TeacherID      *uint64  `json:"teacher_id,omitempty"`
	Name           string   `json:"name,omitempty"`
	RecurrenceDesc string   `json:"recurrence_desc,omitempty"`
	RecurrenceJSON string   `json:"recurrence_json,omitempty"`
	Endereco       *Address `json:"endereco,omitempty"`
}

func (c *Class) Prepare() error {
	c.Name = strings.TrimSpace(c.Name)
	c.RecurrenceDesc = strings.TrimSpace(c.RecurrenceDesc)
	c.RecurrenceJSON = strings.TrimSpace(c.RecurrenceJSON)

	if c.Endereco != nil {
		c.Endereco.CEP = strings.TrimSpace(c.Endereco.CEP)
		c.Endereco.Rua = strings.TrimSpace(c.Endereco.Rua)
		c.Endereco.Numero = strings.TrimSpace(c.Endereco.Numero)
		c.Endereco.Bairro = strings.TrimSpace(c.Endereco.Bairro)
		c.Endereco.Cidade = strings.TrimSpace(c.Endereco.Cidade)
		c.Endereco.Estado = strings.TrimSpace(c.Endereco.Estado)
		c.Endereco.Pais = strings.TrimSpace(c.Endereco.Pais)
		c.Endereco.Complemento = strings.TrimSpace(c.Endereco.Complemento)

		if c.Endereco.CEP == "" &&
			c.Endereco.Rua == "" &&
			c.Endereco.Numero == "" &&
			c.Endereco.Bairro == "" &&
			c.Endereco.Cidade == "" &&
			c.Endereco.Estado == "" &&
			c.Endereco.Complemento == "" {
			c.Endereco = nil
		} else if c.Endereco.Pais == "" {
			c.Endereco.Pais = "Brasil"
		}
	}

	if c.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
