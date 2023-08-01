package models

import (
	"errors"
	"strings"
)

type Students struct {
	ID     int64  `json:"id"`
	Name   string `json:"nome"`
	Active uint64 `json:"ativo,omitempty"`
}

func (s *Students) Prepare() error {
	if err := s.validate(); err != nil {
		return err
	}

	if err := s.format(); err != nil {
		return err
	}

	return nil
}

func (s *Students) validate() error {
	if s.Name == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}

	return nil
}

func (s *Students) format() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Name = strings.ToLower(s.Name)

	return nil
}
