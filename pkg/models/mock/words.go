package mock

import "github.com/guillem-gelabert/go-zeug/pkg/models"

// WordModel is a mock of models/mysql/words.WordModel
type WordModel struct{}

// MockWords is a collection of words to be mapped to mock Words or Cards
var MockWords [][2]string = [][2]string{
	{"das", "Haus"},
	{"das", "Jahr"},
	{"das", "Prozent"},
	{"der", "Euro"},
	{"die", "Zeit"},
	{"die", "Kategorie"},
	{"die", "Stadt"},
	{"das", "Ende"},
	{"die", "Frau"},
	{"das", "Leben"},
	{"das", "Leben"},
}

// Next is a mock of models/mysql/words.Next
func (m *WordModel) Next(offset, next int) ([]*models.Word, error) {
	var r []*models.Word
	for i := offset; i < offset+next; i++ {
		r = append(r, &models.Word{
			ID:          i,
			Priority:    i,
			Article:     MockWords[offset%len(MockWords)][0],
			Substantive: MockWords[offset%len(MockWords)][1],
		})
	}

	return r, nil
}
