package usecase

import (
	"errors"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/repository"
)

// UpdateUserInput represents payload to update user profile.
type UpdateUserInput struct {
	Nama       string `json:"nama"`
	TanggalLahir string `json:"tanggal_Lahir"`
	Tentang    string `json:"tentang"`
	Pekerjaan  string `json:"pekerjaan"`
}

// UserUsecase handles business logic related to user account.
type UserUsecase interface {
	GetProfile(userID uint) (*domain.User, error)
	UpdateProfile(userID uint, in UpdateUserInput) (*domain.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

// ErrUserNotFound is returned when user is not found.
var ErrUserNotFound = errors.New("user not found")

func (uc *userUsecase) GetProfile(userID uint) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (uc *userUsecase) UpdateProfile(userID uint, in UpdateUserInput) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if in.Nama != "" {
		user.Nama = in.Nama
	}
	if in.Tentang != "" {
		user.Tentang = in.Tentang
	}
	if in.Pekerjaan != "" {
		user.Pekerjaan = in.Pekerjaan
	}
	// tanggal_Lahir dari string ke date bisa ditambah nanti jika diperlukan

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
