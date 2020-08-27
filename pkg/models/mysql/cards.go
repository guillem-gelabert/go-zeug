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
