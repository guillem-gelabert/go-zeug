package mock

import (
	"time"

	"github.com/guillem-gelabert/go-zeug/pkg/models"
	dto "github.com/guillem-gelabert/go-zeug/web/dtos"
)

// CardModel is a mock of models/mysql/cards.CardModel
type CardModel struct{}

// Create is a mock of models/mysql/cards.Create
func (m *CardModel) Create(uid int, w *models.Word) (*models.Card, error) {
	switch uid {
	case 2:
		return nil, models.ErrDuplicateEmail
	default:
		return &models.Card{
			ID:                        1,
			WordID:                    1,
			UserID:                    1,
			Stage:                     "UNSEEN",
			NextDueDate:               time.Now(),
			Easiness:                  2.5,
			ConsecutiveCorrectAnswers: 1,
		}, nil
	}
}

// GetByID is a mock of models/mysql/cards.GetByID
func (m *CardModel) GetByID(cid int) (*models.Card, error) {
	switch cid {
	case 2:
		return nil, models.ErrNoRecord
	default:
		return &models.Card{
			ID:                        1,
			WordID:                    1,
			UserID:                    1,
			Stage:                     "UNSEEN",
			NextDueDate:               time.Now(),
			Easiness:                  2.5,
			ConsecutiveCorrectAnswers: 1,
		}, nil
	}
}

// Update is a mock of models/mysql/cards.Update
func (m *CardModel) Update(cid int, correct bool) error {
	return nil
}

// NextSession is a mock of models/mysql/cards.NextSession
func (m *CardModel) NextSession(user *models.User) ([]*dto.CardDTO, error) {
	var r []*dto.CardDTO
	switch user.ID {
	case 2:
		return nil, models.ErrNoRecord
	case 3:
		return nil, models.ErrNoRecord
	default:
		for i, word := range MockWords {
			stage := "UNSEEN"
			if i%2 == 0 {
				stage = "SEEN"
			}

			r = append(r, &dto.CardDTO{
				ID:          i,
				WordID:      i,
				UserID:      user.ID,
				Stage:       stage,
				Article:     word[0],
				Substantive: word[1],
			})
		}

		return r, nil
	}
}

// GetDueBy is a mock of models/mysql/cards.GetDueBy
func (m *CardModel) GetDueBy(uid int, t time.Time) ([]*models.Card, error) {
	return nil, nil
}
