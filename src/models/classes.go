package models

import (
	"errors"
	"strings"
	"time"
)

type Class struct {
	ID             uint64     `json:"id,omitempty"`
	TeacherID      *uint64    `json:"teacher_id,omitempty"`
	Name           string     `json:"name"`
	RecurrenceDesc string     `json:"recurrence_desc,omitempty"`
	RecurrenceJSON string     `json:"recurrence_json,omitempty"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`

	Students []Student `json:"students,omitempty"`
}

type CreatePrivateClassRequest struct {
	StudentID uint64  `json:"student_id"`
	TeacherID *uint64 `json:"teacher_id,omitempty"`
}

func (c *Class) Prepare() error {
	c.Name = strings.TrimSpace(c.Name)
	c.RecurrenceDesc = strings.TrimSpace(c.RecurrenceDesc)
	c.RecurrenceJSON = strings.TrimSpace(c.RecurrenceJSON)

	if c.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
