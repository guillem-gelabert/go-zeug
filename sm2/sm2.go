package sm2

import (
	"fmt"
	"math"
	"time"

	"github.com/guillem-gelabert/go-zeug/pkg/models"
)

func UpdateReviewedCard(card *models.Card, correct bool) (*models.Card, error) {
	easiness, err := getEasiness(correct, card.Easiness)
	if err != nil {
		return nil, err
	}
	card.Easiness = easiness

	cca, err := getCCA(correct, card.ConsecutiveCorrectAnswers)
	if err != nil {
		return nil, err
	}
	card.ConsecutiveCorrectAnswers = cca

	nextDueDate, err := getNextDueDate(correct, card.Easiness, card.ConsecutiveCorrectAnswers)
	if err != nil {
		return nil, err
	}
	card.NextDueDate = nextDueDate
	return card, nil
}

func getEasiness(correct bool, easiness float64) (float64, error) {
	if easiness < 1.3 {
		return 0, fmt.Errorf("Easiness must be higher than 1.3, got %v", easiness)
	}

	if correct {
		return easiness + 0.22, nil
	}
	return math.Max(1.3, easiness-0.8), nil
}

func getCCA(correct bool, cca int) (int, error) {
	if cca < 0 {
		return 0, fmt.Errorf("Consecutive Correct answers must be a positive number")
	}
	if correct {
		return cca + 1, nil
	}

	return 0, nil
}

func getNextDueDate(correct bool, easiness float64, cca int) (time.Time, error) {
	if easiness < 1.3 {
		return time.Time{}, fmt.Errorf("Easiness must be higher than 1.3, got %v", easiness)
	}
	if cca < 0 {
		return time.Time{}, fmt.Errorf("Consecutive Correct answers must be a positive number")
	}

	interval := 1
	if correct {
		a := math.Pow(easiness, float64(cca-1))
		interval = 6 * int(math.Floor(a))
	}

	return time.Now().AddDate(0, 0, interval), nil
}
