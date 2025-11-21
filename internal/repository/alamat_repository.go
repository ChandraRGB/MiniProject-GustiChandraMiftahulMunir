package repository

import (
	"errors"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"gorm.io/gorm"
)

// AlamatRepository defines methods to interact with alamat (address) data.
type AlamatRepository interface {
	GetAllByUser(userID uint, judul string) ([]domain.Alamat, error)
	GetByIDForUser(userID uint, id uint) (*domain.Alamat, error)
	Create(alamat *domain.Alamat) error
	Update(alamat *domain.Alamat) error
	DeleteByIDForUser(userID uint, id uint) error
}

type alamatRepository struct {
	db *gorm.DB
}

// NewAlamatRepository creates a new AlamatRepository implementation.
func NewAlamatRepository(db *gorm.DB) AlamatRepository {
	return &alamatRepository{db: db}
}

func (r *alamatRepository) GetAllByUser(userID uint, judul string) ([]domain.Alamat, error) {
	var list []domain.Alamat

	query := r.db.Where("id_user = ?", userID)
	if judul != "" {
		query = query.Where("judul_alamat LIKE ?", "%"+judul+"%")
	}

	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *alamatRepository) GetByIDForUser(userID uint, id uint) (*domain.Alamat, error) {
	var alamat domain.Alamat
	if err := r.db.Where("id_user = ? AND id = ?", userID, id).First(&alamat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &alamat, nil
}

func (r *alamatRepository) Create(alamat *domain.Alamat) error {
	return r.db.Create(alamat).Error
}

func (r *alamatRepository) Update(alamat *domain.Alamat) error {
	return r.db.Save(alamat).Error
}

func (r *alamatRepository) DeleteByIDForUser(userID uint, id uint) error {
	result := r.db.Where("id_user = ? AND id = ?", userID, id).Delete(&domain.Alamat{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
