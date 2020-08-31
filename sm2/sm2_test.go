package sm2

import (
	"reflect"
	"testing"
	"time"

	"github.com/guillem-gelabert/go-zeug/pkg/models"
)

func TestGetEasiness(t *testing.T) {
	testCases := []struct {
		desc             string
		correct          bool
		previousEasiness float64
		expected         float64
	}{
		{
			desc:             "with correct answer and default easiness",
			correct:          true,
			previousEasiness: 2,
			expected:         2.22,
		},
		{
			desc:             "with incorrect answer and almost min easiness",
			correct:          false,
			previousEasiness: 1.5,
			expected:         1.3,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, err := getEasiness(tC.correct, tC.previousEasiness)
			if err != nil {
				t.Fatal(err)
			}

			if actual != tC.expected {
				t.Errorf("Expected %f; got %f", tC.expected, actual)
			}
		})
	}

	t.Run("With lower than min easiness", func(t *testing.T) {
		_, err := getEasiness(false, 1)
		if err == nil {
			t.Errorf("Expected to fail")
		}
	})
}

func TestGetNextDueDate(t *testing.T) {
	testCases := []struct {
		desc     string
		correct  bool
		easiness float64
		cca      int
		expected time.Time
	}{
		{
			desc:     "first correct answer in new card should reschedule to same day",
			correct:  true,
			easiness: 2,
			cca:      0,
			expected: time.Now(),
		},
		{
			desc:     "second correct answer in new card should reschedule to 6 days",
			correct:  true,
			easiness: 2.2,
			cca:      1,
			expected: time.Now().AddDate(0, 0, 6),
		},
		{
			desc:     "wrong answer should reschedule to next day",
			correct:  false,
			easiness: 2.2,
			cca:      1,
			expected: time.Now().AddDate(0, 0, 1),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, err := getNextDueDate(tC.correct, tC.easiness, tC.cca)
			if err != nil {
				t.Fatal(err)
			}

			actualDate := actual.Format("2006-01-02")
			expectedDate := tC.expected.Format("2006-01-02")

			if actualDate != expectedDate {
				t.Errorf("Expected %q; got %q", expectedDate, actualDate)
			}
		})
	}

	t.Run("With negative CCA", func(t *testing.T) {
		_, err := getNextDueDate(true, 2, -1)
		if err == nil {
			t.Errorf("Expected to fail")
		}
	})

	t.Run("With sub min easiness", func(t *testing.T) {
		_, err := getNextDueDate(true, 1.2, 2)
		if err == nil {
			t.Errorf("Expected to fail")
		}
	})
}

func TestCCA(t *testing.T) {
	testCases := []struct {
		desc        string
		correct     bool
		previousCCA int
		expected    int
	}{
		{
			desc:        "with correct answer",
			correct:     true,
			previousCCA: 2,
			expected:    3,
		},
		{
			desc:        "with incorrect answer",
			correct:     false,
			previousCCA: 15,
			expected:    0,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, err := getCCA(tC.correct, tC.previousCCA)
			if err != nil {
				t.Fatal(err)
			}

			if actual != tC.expected {
				t.Errorf("Expected %d; got %d", tC.expected, actual)
			}
		})
	}

	t.Run("With negative CCA", func(t *testing.T) {
		_, err := getCCA(false, -23)
		if err == nil {
			t.Errorf("Expected to fail")
		}
	})
}

func TestUpdateReviewedCard(t *testing.T) {
	testCases := []struct {
		desc     string
		correct  bool
		target   *models.Card
		expected *models.Card
	}{
		{
			desc:    "with correct answer",
			correct: true,
			target: &models.Card{
				Easiness:                  2,
				ConsecutiveCorrectAnswers: 0,
			},
			expected: &models.Card{
				Easiness:                  2,
				ConsecutiveCorrectAnswers: 1,
				NextDueDate:               time.Now().AddDate(0, 0, 1),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, err := UpdateReviewedCard(tC.target, tC.correct)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tC.expected, actual) {
				t.Errorf("Expected %d; got %d", tC.expected, actual)
			}
		})
	}

	t.Run("With negative CCA", func(t *testing.T) {
		_, err := getCCA(false, -23)
		if err == nil {
			t.Errorf("Expected to fail")
		}
	})
}
