package models

import (
	"errors"
	"log"
	"strings"
	"time"
)

// Classes are the standard teaching groups of the system, it can only have a single Teacher assigned to it but may have multiple Students.
type Classes struct {
	ID       uint64     `json:"id_turma,omitempty"`
	Name     string     `json:"nome_turma,omitempty"`
	Teacher  Teachers   `json:"id_professor,omitempty"`
	Creation time.Time  `json:"data_criacao,omitempty"`
	Students []Students `json:"alunos"`
}

func (c *Classes) Prepare() error {
	if err := c.validate(); err != nil {
		return err
	}

	if err := c.format(); err != nil {
		return err
	}

	return nil
}

func (c *Classes) validate() error {
	log.Println("DEBUG: Entrou no validate")

	if c.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}

	return nil
}

func (c *Classes) format() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Name = strings.ToLower(c.Name)

	return nil
}
