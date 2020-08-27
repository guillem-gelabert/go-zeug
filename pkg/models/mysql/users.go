package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/guillem-gelabert/go-zeug/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// UserModel with embedded DB
type UserModel struct {
	DB *sql.DB
}

// Insert adds a user to the Database
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (displayName, email, password, createdAt, updatedAt)
		VALUES(?,?,?,UTC_TIMESTAMP(), UTC_TIMESTAMP())
		`
	_, err = m.DB.Exec(stmt, name, email, hashedPassword)
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

// Authenticate checks if email and password match hash and returns User ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	row := m.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email)
	if err := row.Scan(&id, &hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

// Get fetches details for a specific user based on their user ID
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}
	stmt := `
		SELECT
			id,
			lastSeenPriority,
			lastUpdate,
			displayName,
			newWordsPerSession,
			createdAt,
			updatedAt,
			email,
			password
		FROM users WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(
		&u.ID,
		&u.LastSeenPriority,
		&u.LastUpdate,
		&u.DisplayName,
		&u.NewWordsPerSession,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Email,
		&u.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return u, nil
}
