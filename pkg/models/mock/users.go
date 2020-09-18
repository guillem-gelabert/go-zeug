package mock

import (
	"time"

	"github.com/guillem-gelabert/go-zeug/pkg/models"
)

var mockUser = &models.User{
	ID:                 1,
	DisplayName:        "abigail_adams",
	Email:              "abigail_adams@whitehouse.gov",
	NewWordsPerSession: 10,
	LastSeenPriority:   0,
	Password:           []byte("cr42yp455w0rd"),
	LastUpdate:         time.Time{},
	CreatedAt:          time.Time{},
	UpdatedAt:          time.Time{},
}

// UserModel is a mock of models/mysql/users.UserModel
type UserModel struct{}

// Insert is a mock of models/mysql/users.Insert
func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

// Authenticate is a mock of models/mysql/users.Authenticate
func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "abigail_adams@whitehouse.gov":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

// Get is a mock of models/mysql/users.Get
func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// Update is a mock of models/mysql/users.Update
func (m *UserModel) Update(u *models.User) error {
	switch u.ID {
	case 1:
		return nil
	default:
		return models.ErrNoRecord
	}
}
