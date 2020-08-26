package mysql

import (
	"database/sql"

	"github.com/guillem-gelabert/go-zeug/pkg/models"
)

// WordModel wraps a sql.DB connection pool
type WordModel struct {
	DB *sql.DB
}

// Next returns the next n words from an offset m
func (m *WordModel) Next(offset, next int) ([]*models.Word, error) {
	var words []*models.Word
	stmt := `SELECT * FROM words LIMIT ? OFFSET?;`
	rows, err := m.DB.Query(stmt, offset, next)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		w := &models.Word{}
		err = rows.Scan(&w.ID, &w.Article, &w.Substantive, &w.Priority)
		if err != nil {
			return nil, err
		}

		words = append(words, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
