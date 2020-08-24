package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type User struct {
	ID                 int
	DisplayName        string
	Email              string
	NewWordsPerSession int
	LastSeenPriority   int
	Password           []byte
	LastUpdate         time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
