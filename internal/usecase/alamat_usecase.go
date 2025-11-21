package usecase

import (
	"errors"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/repository"
	"gorm.io/gorm"
)

// CreateAlamatInput represents payload to create a new alamat.
type CreateAlamatInput struct {
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	DetailAlamat string `json:"detail_alamat"`
}

// UpdateAlamatInput represents payload to update an alamat.
type UpdateAlamatInput struct {
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	DetailAlamat string `json:"detail_alamat"`
}

// AlamatUsecase handles alamat-related business logic.
type AlamatUsecase interface {
	GetMyAlamat(userID uint, filterJudul string) ([]domain.Alamat, error)
	GetByID(userID uint, id uint) (*domain.Alamat, error)
	Create(userID uint, in CreateAlamatInput) (*domain.Alamat, error)
	Update(userID uint, id uint, in UpdateAlamatInput) (*domain.Alamat, error)
	Delete(userID uint, id uint) error
}

type alamatUsecase struct {
	alamatRepo repository.AlamatRepository
}

// NewAlamatUsecase creates a new AlamatUsecase.
func NewAlamatUsecase(alamatRepo repository.AlamatRepository) AlamatUsecase {
	return &alamatUsecase{alamatRepo: alamatRepo}
}

// ErrAlamatNotFound indicates alamat not found.
var ErrAlamatNotFound = errors.New("alamat not found")

func (uc *alamatUsecase) GetMyAlamat(userID uint, filterJudul string) ([]domain.Alamat, error) {
	return uc.alamatRepo.GetAllByUser(userID, filterJudul)
}

func (uc *alamatUsecase) GetByID(userID uint, id uint) (*domain.Alamat, error) {
	alamat, err := uc.alamatRepo.GetByIDForUser(userID, id)
	if err != nil {
		return nil, err
	}
	if alamat == nil {
		return nil, ErrAlamatNotFound
	}
	return alamat, nil
}

func (uc *alamatUsecase) Create(userID uint, in CreateAlamatInput) (*domain.Alamat, error) {
	if in.JudulAlamat == "" || in.NamaPenerima == "" || in.NoTelp == "" || in.DetailAlamat == "" {
		return nil, errors.New("judul_alamat, nama_penerima, no_telp, detail_alamat wajib diisi")
	}

	alamat := &domain.Alamat{
		UserID:       userID,
		JudulAlamat:  in.JudulAlamat,
		NamaPenerima: in.NamaPenerima,
		NoTelp:       in.NoTelp,
		DetailAlamat: in.DetailAlamat,
	}

	if err := uc.alamatRepo.Create(alamat); err != nil {
		return nil, err
	}

	return alamat, nil
}

func (uc *alamatUsecase) Update(userID uint, id uint, in UpdateAlamatInput) (*domain.Alamat, error) {
	alamat, err := uc.alamatRepo.GetByIDForUser(userID, id)
	if err != nil {
		return nil, err
	}
	if alamat == nil {
		return nil, ErrAlamatNotFound
	}

	if in.JudulAlamat != "" {
		alamat.JudulAlamat = in.JudulAlamat
	}
	if in.NamaPenerima != "" {
		alamat.NamaPenerima = in.NamaPenerima
	}
	if in.NoTelp != "" {
		alamat.NoTelp = in.NoTelp
	}
	if in.DetailAlamat != "" {
		alamat.DetailAlamat = in.DetailAlamat
	}

	if err := uc.alamatRepo.Update(alamat); err != nil {
		return nil, err
	}

	return alamat, nil
}

func (uc *alamatUsecase) Delete(userID uint, id uint) error {
	if err := uc.alamatRepo.DeleteByIDForUser(userID, id); err != nil {
		if errors.Is(err, ErrAlamatNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAlamatNotFound
		}
		return err
	}
	return nil
}
