package models

import (
	"errors"
	"strings"
	"time"
)

type Lesson struct {
	ID         uint64    `json:"id,omitempty"`
	StatusID   uint64    `json:"status_id,omitempty"`
	TeacherID  *uint64   `json:"teacher_id,omitempty"`
	ClassID    uint64    `json:"class_id"`
	Subject    string    `json:"subject,omitempty"`
	Vocabulary string    `json:"vocabulary,omitempty"`
	Balance    string    `json:"balance,omitempty"`
	Notes      string    `json:"notes,omitempty"`
	LessonDate time.Time `json:"lesson_date"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`

	Students []Student `json:"students,omitempty"`
}

func (l *Lesson) Prepare() error {
	l.Subject = strings.TrimSpace(l.Subject)
	l.Vocabulary = strings.TrimSpace(l.Vocabulary)
	l.Balance = strings.TrimSpace(l.Balance)
	l.Notes = strings.TrimSpace(l.Notes)

	if l.ClassID == 0 {
		return errors.New("class_id is required")
	}

	if l.LessonDate.IsZero() {
		return errors.New("lesson_date is required")
	}

	return nil
}

func (l *Lesson) PrepareUpdate() error {
	l.Subject = strings.TrimSpace(l.Subject)
	l.Vocabulary = strings.TrimSpace(l.Vocabulary)
	l.Balance = strings.TrimSpace(l.Balance)
	l.Notes = strings.TrimSpace(l.Notes)

	if l.LessonDate.IsZero() {
		return errors.New("lesson_date is required")
	}

	return nil
}
