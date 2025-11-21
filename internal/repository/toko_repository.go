package repository

import (
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"gorm.io/gorm"
)

// TokoRepository defines methods to interact with the toko table.
type TokoRepository interface {
	Create(toko *domain.Toko) error
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