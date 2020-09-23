package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/guillem-gelabert/go-zeug/pkg/models"
	"github.com/guillem-gelabert/go-zeug/sm2"
	dto "github.com/guillem-gelabert/go-zeug/web/dtos"
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
	nextDueDate,
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
		c := models.Card{}
		err := rows.Scan(
			&c.ID,
			&c.WordID,
			&c.UserID,
			&c.Stage,
			&c.NextDueDate,
			&c.Easiness,
			&c.ConsecutiveCorrectAnswers,
		)
		if err != nil {
			return nil, err
		}

		cards = append(cards, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

// Create adds a card to a user from a word
func (m *CardModel) Create(uid int, w *models.Word) (*models.Card, error) {
	stmt := `INSERT INTO cards (wordId,	userId) VALUES (?,?)`

	_, err := m.DB.Exec(stmt, w.ID, uid)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return nil, models.ErrDuplicateEmail
			}
		}
		return nil, err
	}

	card := models.Card{}
	row := m.DB.QueryRow(
		`SELECT
			id,
			wordID,
			userID,
			stage,
			nextDueDate,
			easiness,
			consecutiveCorrectAnswers
		FROM cards
		WHERE id = LAST_INSERT_ID()`)

	err = row.Scan(
		&card.ID,
		&card.WordID,
		&card.UserID,
		&card.Stage,
		&card.NextDueDate,
		&card.Easiness,
		&card.ConsecutiveCorrectAnswers,
	)

	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (m *CardModel) GetByID(cid int) (*models.Card, error) {
	stmt := `SELECT
			id,
			wordID,
			userID,
			stage,
			nextDueDate,
			easiness,
			consecutiveCorrectAnswers
		FROM cards
		WHERE id = ?`

	card := models.Card{}
	row := m.DB.QueryRow(stmt, cid)
	err := row.Scan(
		&card.ID,
		&card.WordID,
		&card.UserID,
		&card.Stage,
		&card.NextDueDate,
		&card.Easiness,
		&card.ConsecutiveCorrectAnswers,
	)

	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (m *CardModel) Update(cid int, correct bool) error {
	c, err := m.GetByID(cid)
	if err != nil {
		return err
	}
	c, _ = sm2.UpdateReviewedCard(c, correct)

	stmt := `UPDATE cards
		SET easiness = ?,
		consecutiveCorrectAnswers = ?,
		nextDueDate = ?
		WHERE id = ?`

	_, err = m.DB.Exec(stmt, c.Easiness, c.ConsecutiveCorrectAnswers, c.NextDueDate, c.ID)
	if err != nil {
		return err
	}

	return nil
}

// NextSession adds words to cards table, updates the user's last seen priority and returns userDtos
func (m *CardModel) NextSession(user *models.User) ([]*dto.CardDTO, error) {

	cardInsert := sq.Insert("cards").
		Columns("wordId", "userId")

	for i := 0; i < user.NewWordsPerSession; i++ {
		cardInsert = cardInsert.Values(user.LastSeenPriority+1+i, user.ID)
	}

	_, err := cardInsert.RunWith(m.DB).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	userInsert := sq.Update("users").
		Where("id = ?", user.ID).
		Set("lastSeenPriority", user.LastSeenPriority+user.NewWordsPerSession).
		Set("lastUpdate", time.Now())

	if _, err = userInsert.RunWith(m.DB).Exec(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	cardSelect := sq.Select("cards.id", "article", "substantive", "stage", "wordId").
		From("cards").
		Join("words ON (wordId = words.id)").
		Limit(uint64(user.NewWordsPerSession)).
		Offset(uint64(user.LastSeenPriority))

	rows, err := cardSelect.RunWith(m.DB).Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	var cards []*dto.CardDTO
	for rows.Next() {
		card := &dto.CardDTO{
			UserID: user.ID,
		}
		err = rows.Scan(&card.ID, &card.Article, &card.Substantive, &card.Stage, &card.WordID)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}
