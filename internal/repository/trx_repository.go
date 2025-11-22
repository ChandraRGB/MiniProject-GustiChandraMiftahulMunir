package repository

import (
	"errors"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"gorm.io/gorm"
)

// TrxRepository defines DB operations for transaksi and related details.
type TrxRepository interface {
	CreateWithDetails(trx *domain.Trx, logs []domain.LogProduk, details []domain.DetailTrx, products []*domain.Produk) error
	GetAllByUser(userID uint) ([]domain.Trx, error)
	GetByIDForUser(userID, trxID uint) (*domain.Trx, error)
}

type trxRepository struct {
	db *gorm.DB
}

// NewTrxRepository creates a new TrxRepository.
func NewTrxRepository(db *gorm.DB) TrxRepository {
	return &trxRepository{db: db}
}

func (r *trxRepository) CreateWithDetails(trx *domain.Trx, logs []domain.LogProduk, details []domain.DetailTrx, products []*domain.Produk) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trx).Error; err != nil {
			return err
		}

		if len(logs) != len(details) {
			return errors.New("logs and details length mismatch")
		}

		// Create log_produk entries
		if len(logs) > 0 {
			if err := tx.Create(&logs).Error; err != nil {
				return err
			}

			// Attach foreign keys and create detail_trx entries
			for i := range details {
				details[i].TrxID = trx.ID
				details[i].LogProdukID = logs[i].ID
			}

			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		// Update product stocks
		for _, p := range products {
			if err := tx.Model(&domain.Produk{}).
				Where("id = ?", p.ID).
				Update("stok", p.Stok).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *trxRepository) GetAllByUser(userID uint) ([]domain.Trx, error) {
	var trxs []domain.Trx
	if err := r.db.
		Where("id_user = ?", userID).
		Preload("Alamat").
		Preload("DetailTrx.LogProduk.Category").
		Preload("DetailTrx.LogProduk.Toko").
		Preload("DetailTrx.LogProduk.Produk.FotoProduk").
		Preload("DetailTrx.Toko").
		Find(&trxs).Error; err != nil {
		return nil, err
	}
	return trxs, nil
}

func (r *trxRepository) GetByIDForUser(userID, trxID uint) (*domain.Trx, error) {
	var trx domain.Trx
	if err := r.db.
		Where("id_user = ? AND id = ?", userID, trxID).
		Preload("Alamat").
		Preload("DetailTrx.LogProduk.Category").
		Preload("DetailTrx.LogProduk.Toko").
		Preload("DetailTrx.LogProduk.Produk.FotoProduk").
		Preload("DetailTrx.Toko").
		First(&trx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &trx, nil
}
