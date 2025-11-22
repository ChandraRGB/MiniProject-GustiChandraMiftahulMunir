package usecase

import (
	"errors"
	"strconv"
	"strings"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/repository"
)

// Re-export ProductFilter so delivery layer can use it without depending on repository.
type ProductFilter = repository.ProductFilter

// ProductListResult wraps paginated product list.
type ProductListResult struct {
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Data  []domain.Produk `json:"data"`
}

// CreateProductInput represents required fields to create a product.
type CreateProductInput struct {
	NamaProduk    string
	CategoryID    uint
	HargaReseller int
	HargaKonsumen int
	Stok          int
	Deskripsi     string
}

// UpdateProductInput represents optional fields to update a product.
type UpdateProductInput struct {
	NamaProduk    *string
	CategoryID    *uint
	HargaReseller *int
	HargaKonsumen *int
	Stok          *int
	Deskripsi     *string
}

// ProductUsecase defines product-related business logic.
type ProductUsecase interface {
	GetAll(limit, page int, filter ProductFilter) (*ProductListResult, error)
	GetByID(id uint) (*domain.Produk, error)
	Create(userID uint, in CreateProductInput, photoFilenames []string) (*domain.Produk, error)
	Update(userID uint, productID uint, in UpdateProductInput, photoFilenames []string) (*domain.Produk, error)
	Delete(userID uint, productID uint) error
}

type productUsecase struct {
	productRepo repository.ProductRepository
	fotoRepo    repository.FotoProdukRepository
	tokoRepo    repository.TokoRepository
}

// NewProductUsecase creates a new ProductUsecase.
func NewProductUsecase(productRepo repository.ProductRepository, fotoRepo repository.FotoProdukRepository, tokoRepo repository.TokoRepository) ProductUsecase {
	return &productUsecase{productRepo: productRepo, fotoRepo: fotoRepo, tokoRepo: tokoRepo}
}

var (
	// ErrProductNotFound indicates product not found.
	ErrProductNotFound = errors.New("product not found")
)

func (uc *productUsecase) GetAll(limit, page int, filter ProductFilter) (*ProductListResult, error) {
	products, err := uc.productRepo.GetAll(limit, page, filter)
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	return &ProductListResult{
		Page:  page,
		Limit: limit,
		Data:  products,
	}, nil
}

func (uc *productUsecase) GetByID(id uint) (*domain.Produk, error) {
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (uc *productUsecase) Create(userID uint, in CreateProductInput, photoFilenames []string) (*domain.Produk, error) {
	if in.NamaProduk == "" || in.CategoryID == 0 || in.HargaReseller <= 0 || in.HargaKonsumen <= 0 || in.Stok < 0 {
		return nil, errors.New("nama_produk, category_id, harga_reseller, harga_konsumen, stok wajib diisi")
	}

	// Ensure user has a toko
	toko, err := uc.tokoRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if toko == nil {
		return nil, errors.New("toko not found for user")
	}

	product := &domain.Produk{
		NamaProduk:    in.NamaProduk,
		Slug:          slugify(in.NamaProduk),
		HargaReseller: strconv.Itoa(in.HargaReseller),
		HargaKonsumen: strconv.Itoa(in.HargaKonsumen),
		Stok:          in.Stok,
		Deskripsi:     in.Deskripsi,
		TokoID:        toko.ID,
		CategoryID:    in.CategoryID,
	}

	if err := uc.productRepo.Create(product); err != nil {
		return nil, err
	}

	// Save product photos
	if len(photoFilenames) > 0 {
		photos := make([]domain.FotoProduk, 0, len(photoFilenames))
		for _, filename := range photoFilenames {
			if filename == "" {
				continue
			}
			photos = append(photos, domain.FotoProduk{
				ProdukID: product.ID,
				URL:      filename,
			})
		}
		if err := uc.fotoRepo.CreateMany(photos); err != nil {
			return nil, err
		}
	}

	return product, nil
}

func (uc *productUsecase) Update(userID uint, productID uint, in UpdateProductInput, photoFilenames []string) (*domain.Produk, error) {
	// Ensure user has a toko
	toko, err := uc.tokoRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if toko == nil {
		return nil, errors.New("toko not found for user")
	}

	product, err := uc.productRepo.GetByIDForToko(toko.ID, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	if in.NamaProduk != nil && *in.NamaProduk != "" {
		product.NamaProduk = *in.NamaProduk
		product.Slug = slugify(*in.NamaProduk)
	}
	if in.CategoryID != nil && *in.CategoryID != 0 {
		product.CategoryID = *in.CategoryID
	}
	if in.HargaReseller != nil {
		product.HargaReseller = strconv.Itoa(*in.HargaReseller)
	}
	if in.HargaKonsumen != nil {
		product.HargaKonsumen = strconv.Itoa(*in.HargaKonsumen)
	}
	if in.Stok != nil {
		product.Stok = *in.Stok
	}
	if in.Deskripsi != nil {
		product.Deskripsi = *in.Deskripsi
	}

	if err := uc.productRepo.Update(product); err != nil {
		return nil, err
	}

	// If new photos are provided, replace existing photos
	if len(photoFilenames) > 0 {
		if err := uc.fotoRepo.DeleteByProdukID(product.ID); err != nil {
			return nil, err
		}

		photos := make([]domain.FotoProduk, 0, len(photoFilenames))
		for _, filename := range photoFilenames {
			if filename == "" {
				continue
			}
			photos = append(photos, domain.FotoProduk{
				ProdukID: product.ID,
				URL:      filename,
			})
		}
		if err := uc.fotoRepo.CreateMany(photos); err != nil {
			return nil, err
		}
	}

	return product, nil
}

func (uc *productUsecase) Delete(userID uint, productID uint) error {
	// Ensure user has a toko
	toko, err := uc.tokoRepo.GetByUserID(userID)
	if err != nil {
		return err
	}
	if toko == nil {
		return errors.New("toko not found for user")
	}

	product, err := uc.productRepo.GetByIDForToko(toko.ID, productID)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	if err := uc.fotoRepo.DeleteByProdukID(product.ID); err != nil {
		return err
	}

	return uc.productRepo.Delete(product.ID)
}

// slugify converts a product name into a URL-friendly slug.
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	// naive implementation is enough for this case
	return s
}
