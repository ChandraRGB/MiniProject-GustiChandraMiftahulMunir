package usecase

import (
	"errors"
	"time"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/helper"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/repository"
)

// RegisterInput represents expected payload for /auth/register.
type RegisterInput struct {
	Nama         string `json:"nama"`
	KataSandi    string `json:"kata_sandi"`
	NoTelp       string `json:"no_telp"`
	TanggalLahir string `json:"tanggal_Lahir"`
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

// LoginInput represents expected payload for /auth/login.
type LoginInput struct {
	NoTelp    string `json:"no_telp"`
	KataSandi string `json:"kata_sandi"`
}

// LoginResult is returned after successful login.
type LoginResult struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

// AuthUsecase exposes authentication use cases.
type AuthUsecase interface {
	Register(in RegisterInput) error
	Login(in LoginInput) (*LoginResult, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
	tokoRepo repository.TokoRepository
}

// NewAuthUsecase constructs a new AuthUsecase implementation.
func NewAuthUsecase(userRepo repository.UserRepository, tokoRepo repository.TokoRepository) AuthUsecase {
	return &authUsecase{userRepo: userRepo, tokoRepo: tokoRepo}
}

var (
	// ErrEmailOrPhoneExists indicates duplicate email or phone.
	ErrEmailOrPhoneExists = errors.New("email or phone already exists")
	// ErrInvalidCredentials indicates login credential mismatch.
	ErrInvalidCredentials = errors.New("no telp atau kata sandi salah")
)

func (uc *authUsecase) Register(in RegisterInput) error {
	if in.Nama == "" || in.KataSandi == "" || in.NoTelp == "" || in.Email == "" {
		return errors.New("nama, kata_sandi, no_telp, dan email wajib diisi")
	}

	exists, err := uc.userRepo.IsEmailOrNoTelpExists(in.Email, in.NoTelp)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailOrPhoneExists
	}

	hashedPassword, err := helper.HashPassword(in.KataSandi)
	if err != nil {
		return err
	}

	var tglLahir time.Time
	if in.TanggalLahir != "" {
		// format dari Postman: 02/01/2006
		parsed, err := time.Parse("02/01/2006", in.TanggalLahir)
		if err == nil {
			tglLahir = parsed
		}
	}

	user := &domain.User{
		Nama:         in.Nama,
		KataSandi:    hashedPassword,
		NoTelp:       in.NoTelp,
		TanggalLahir: tglLahir,
		Pekerjaan:    in.Pekerjaan,
		Email:        in.Email,
		IDProvinsi:   in.IDProvinsi,
		IDKota:       in.IDKota,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return err
	}

	// Automatically create store for new user
	toko := &domain.Toko{
		UserID:   user.ID,
		NamaToko: in.Nama + " Store",
	}

	if err := uc.tokoRepo.Create(toko); err != nil {
		return err
	}

	return nil
}

func (uc *authUsecase) Login(in LoginInput) (*LoginResult, error) {
	if in.NoTelp == "" || in.KataSandi == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := uc.userRepo.FindByNoTelp(in.NoTelp)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := helper.CheckPasswordHash(user.KataSandi, in.KataSandi); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := helper.GenerateJWT(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
		User:  user,
	}, nil
}