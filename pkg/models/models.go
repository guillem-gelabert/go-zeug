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

// User represents a User as retrieved from the DB
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

// Word represents a Word as retrieved from the DB
type Word struct {
	ID          int
	Article     string
	Substantive string
	Priority    int
}

// Card represents a Card as retrieved from the DB
type Card struct {
	ID                        int
	WordID                    int
	UserID                    int
	Stage                     string
	NextDueDate               time.Time
	Easiness                  float64
	ConsecutiveCorrectAnswers int
}
