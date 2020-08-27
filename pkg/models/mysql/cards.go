package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/guillem-gelabert/go-zeug/pkg/models"
)

// CardModel with embedded DB
type CardModel struct {
	DB *sql.DB
}

// GetDueBy returns due cards at one given date
func (m *CardModel) GetDueBy(uid int, t time.Time) ([]*models.Card, error) {
	stmt := `
	SELECT
	id,
	wordId,
	userId,
	stage,
	nextDudeDate,
	easiness,
	consecutiveCorrectAnswers
	FROM cards
	WHERE userId = ? AND nextDueDate >= ?
	`
	rows, err := m.DB.Query(stmt, uid, t)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("No due cards")
		}
		return nil, err
	}

	var cards []*models.Card
	for rows.Next() {
		var c *models.Card
		err := rows.Scan(c)
		if err != nil {
			return nil, err
		}

		cards = append(cards, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

// Create adds a card to a user from a word
func (m *CardModel) Create(uid int, w *models.Word) error {
	stmt := `INSERT INTO cards (wordId,	userId) VALUES (?,?)`

	_, err := m.DB.Exec(stmt, w.ID, uid)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}
