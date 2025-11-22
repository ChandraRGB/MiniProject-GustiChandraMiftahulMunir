package usecase

import (
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/repository"
)

// TokoListResult wraps paginated toko list.
type TokoListResult struct {
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
	Data  []domain.Toko `json:"data"`
}

// UpdateTokoInput represents fields allowed to update store profile.
type UpdateTokoInput struct {
	NamaToko string `json:"nama_toko"`
	UrlFoto  string `json:"url_foto"`
}

// TokoUsecase defines toko-related business logic.
type TokoUsecase interface {
	GetAll(limit, page int, nama string) (*TokoListResult, error)
	GetByID(id uint) (*domain.Toko, error)
	GetMyStore(userID uint) (*domain.Toko, error)
	UpdateMyStore(userID uint, in UpdateTokoInput) (*domain.Toko, error)
}

type tokoUsecase struct {
	tokoRepo repository.TokoRepository
}

// NewTokoUsecase creates a new TokoUsecase.
func NewTokoUsecase(tokoRepo repository.TokoRepository) TokoUsecase {
	return &tokoUsecase{tokoRepo: tokoRepo}
}

func (uc *tokoUsecase) GetAll(limit, page int, nama string) (*TokoListResult, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	tokos, err := uc.tokoRepo.GetAll(limit, page, nama)
	if err != nil {
		return nil, err
	}

	return &TokoListResult{
		Page:  page,
		Limit: limit,
		Data:  tokos,
	}, nil
}

func (uc *tokoUsecase) GetByID(id uint) (*domain.Toko, error) {
	return uc.tokoRepo.GetByID(id)
}

func (uc *tokoUsecase) GetMyStore(userID uint) (*domain.Toko, error) {
	return uc.tokoRepo.GetByUserID(userID)
}

func (uc *tokoUsecase) UpdateMyStore(userID uint, in UpdateTokoInput) (*domain.Toko, error) {
	toko, err := uc.tokoRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if toko == nil {
		return nil, nil
	}

	if in.NamaToko != "" {
		toko.NamaToko = in.NamaToko
	}
	if in.UrlFoto != "" {
		toko.UrlFoto = in.UrlFoto
	}

	if err := uc.tokoRepo.Update(toko); err != nil {
		return nil, err
	}

	return toko, nil
}
