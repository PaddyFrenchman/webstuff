package core

import (
	"go-web/internal/models"

	"github.com/jinzhu/gorm"
)

// UserStore implements the UserStore interface
type UserStore struct {
	DB *gorm.DB
}

func (s *UserStore) First(m *models.User) error {
	return s.DB.First(m).Error
}

func (s *UserStore) Create(m *models.User) error {
	return s.DB.Create(m).Error
}

func (s *UserStore) Find(m *[]models.User) error {
	return s.DB.Find(m).Error
}

func (s *UserStore) Ping() error {
	return s.DB.DB().Ping()
}
