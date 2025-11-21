package repository

import (
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"gorm.io/gorm"
)

// TokoRepository defines methods to interact with the toko table.
type TokoRepository interface {
	Create(toko *domain.Toko) error
	GetAll(limit, page int, nama string) ([]domain.Toko, error)
	GetByID(id uint) (*domain.Toko, error)
	GetByUserID(userID uint) (*domain.Toko, error)
	Update(toko *domain.Toko) error
}

type tokoRepository struct {
	db *gorm.DB
}

// NewTokoRepository creates a new TokoRepository implementation.
func NewTokoRepository(db *gorm.DB) TokoRepository {
	return &tokoRepository{db: db}
}

func (r *tokoRepository) Create(toko *domain.Toko) error {
	return r.db.Create(toko).Error
}

func (r *tokoRepository) GetAll(limit, page int, nama string) ([]domain.Toko, error) {
	var list []domain.Toko
	query := r.db.Model(&domain.Toko{})
	if nama != "" {
		query = query.Where("nama_toko LIKE ?", "%"+nama+"%")
	}
	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *tokoRepository) GetByID(id uint) (*domain.Toko, error) {
	var toko domain.Toko
	if err := r.db.First(&toko, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &toko, nil
}

func (r *tokoRepository) GetByUserID(userID uint) (*domain.Toko, error) {
	var toko domain.Toko
	if err := r.db.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &toko, nil
}

func (r *tokoRepository) Update(toko *domain.Toko) error {
	return r.db.Save(toko).Error
}