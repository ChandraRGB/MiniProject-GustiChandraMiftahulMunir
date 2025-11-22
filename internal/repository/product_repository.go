package repository

import (
	"errors"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"gorm.io/gorm"
)

// ProductFilter represents filters for listing products.
type ProductFilter struct {
	NamaProduk string
	CategoryID uint
	TokoID     uint
	MinHarga   int
	MaxHarga   int
}

// ProductRepository defines DB operations for produk.
type ProductRepository interface {
	GetAll(limit, page int, filter ProductFilter) ([]domain.Produk, error)
	GetByID(id uint) (*domain.Produk, error)
	GetByIDForToko(tokoID, productID uint) (*domain.Produk, error)
	Create(product *domain.Produk) error
	Update(product *domain.Produk) error
	Delete(id uint) error
}

// FotoProdukRepository defines DB operations for foto_produk.
type FotoProdukRepository interface {
	CreateMany(photos []domain.FotoProduk) error
	DeleteByProdukID(produkID uint) error
}

type productRepository struct {
	db *gorm.DB
}

type fotoProdukRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// NewFotoProdukRepository creates a new FotoProdukRepository.
func NewFotoProdukRepository(db *gorm.DB) FotoProdukRepository {
	return &fotoProdukRepository{db: db}
}

func (r *productRepository) GetAll(limit, page int, filter ProductFilter) ([]domain.Produk, error) {
	var products []domain.Produk

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	db := r.db.Model(&domain.Produk{}).
		Preload("Toko").
		Preload("Category").
		Preload("FotoProduk")

	if filter.NamaProduk != "" {
		db = db.Where("nama_produk LIKE ?", "%"+filter.NamaProduk+"%")
	}
	if filter.CategoryID != 0 {
		db = db.Where("id_category = ?", filter.CategoryID)
	}
	if filter.TokoID != 0 {
		db = db.Where("id_toko = ?", filter.TokoID)
	}
	if filter.MinHarga > 0 {
		// filter based on harga_konsumen numeric value
		db = db.Where("CAST(harga_konsumen AS UNSIGNED) >= ?", filter.MinHarga)
	}
	if filter.MaxHarga > 0 {
		db = db.Where("CAST(harga_konsumen AS UNSIGNED) <= ?", filter.MaxHarga)
	}

	if err := db.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) GetByID(id uint) (*domain.Produk, error) {
	var product domain.Produk
	if err := r.db.Preload("Toko").Preload("Category").Preload("FotoProduk").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetByIDForToko(tokoID, productID uint) (*domain.Produk, error) {
	var product domain.Produk
	if err := r.db.Where("id_toko = ?", tokoID).
		Preload("Toko").Preload("Category").Preload("FotoProduk").
		First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *domain.Produk) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *domain.Produk) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Produk{}, id).Error
}

func (r *fotoProdukRepository) CreateMany(photos []domain.FotoProduk) error {
	if len(photos) == 0 {
		return nil
	}
	return r.db.Create(&photos).Error
}

func (r *fotoProdukRepository) DeleteByProdukID(produkID uint) error {
	return r.db.Where("id_produk = ?", produkID).Delete(&domain.FotoProduk{}).Error
}
