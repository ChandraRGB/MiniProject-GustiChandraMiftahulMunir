package repository

import (
	"errors"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"gorm.io/gorm"
)

// UserRepository defines methods to interact with the users table.
type UserRepository interface {
	Create(user *domain.User) error
	FindByNoTelp(noTelp string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Update(user *domain.User) error
	IsEmailOrNoTelpExists(email, noTelp string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository implementation.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByNoTelp(noTelp string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) IsEmailOrNoTelpExists(email, noTelp string) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.User{}).
		Where("email = ? OR notelp = ?", email, noTelp).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}