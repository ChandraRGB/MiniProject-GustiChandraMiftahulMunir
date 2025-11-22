package usecase

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/repository"
)

// TrxItemInput represents a single item in the transaction request.
type TrxItemInput struct {
	ProductID uint `json:"product_id"`
	Kuantitas int  `json:"kuantitas"`
}

// CreateTrxInput represents the payload to create a new transaction.
type CreateTrxInput struct {
	MethodBayar string         `json:"method_bayar"`
	AlamatKirim uint           `json:"alamat_kirim"`
	DetailTrx   []TrxItemInput `json:"detail_trx"`
}

// TrxUsecase defines transaction-related business logic.
type TrxUsecase interface {
	GetAll(userID uint) ([]domain.Trx, error)
	GetByID(userID, trxID uint) (*domain.Trx, error)
	Create(userID uint, in CreateTrxInput) (*domain.Trx, error)
}

type trxUsecase struct {
	trxRepo     repository.TrxRepository
	alamatRepo  repository.AlamatRepository
	productRepo repository.ProductRepository
}

// NewTrxUsecase creates a new TrxUsecase.
func NewTrxUsecase(trxRepo repository.TrxRepository, alamatRepo repository.AlamatRepository, productRepo repository.ProductRepository) TrxUsecase {
	return &trxUsecase{trxRepo: trxRepo, alamatRepo: alamatRepo, productRepo: productRepo}
}

var (
	// ErrTrxNotFound indicates transaction not found.
	ErrTrxNotFound = errors.New("trx not found")
	// ErrTrxAlamatNotFound indicates shipping address not found for user.
	ErrTrxAlamatNotFound = errors.New("alamat pengiriman not found")
	// ErrTrxProductNotFound indicates one of the products in the trx was not found.
	ErrTrxProductNotFound = errors.New("product not found")
	// ErrTrxInsufficientStock indicates product stock is insufficient.
	ErrTrxInsufficientStock = errors.New("insufficient stock")
	// ErrTrxEmptyDetail indicates empty detail_trx payload.
	ErrTrxEmptyDetail = errors.New("detail_trx empty")
)

func (uc *trxUsecase) GetAll(userID uint) ([]domain.Trx, error) {
	return uc.trxRepo.GetAllByUser(userID)
}

func (uc *trxUsecase) GetByID(userID, trxID uint) (*domain.Trx, error) {
	trx, err := uc.trxRepo.GetByIDForUser(userID, trxID)
	if err != nil {
		return nil, err
	}
	if trx == nil {
		return nil, ErrTrxNotFound
	}
	return trx, nil
}

func (uc *trxUsecase) Create(userID uint, in CreateTrxInput) (*domain.Trx, error) {
	if in.MethodBayar == "" || in.AlamatKirim == 0 {
		return nil, errors.New("method_bayar and alamat_kirim wajib diisi")
	}
	if len(in.DetailTrx) == 0 {
		return nil, ErrTrxEmptyDetail
	}

	// Ensure alamat pengiriman belongs to the user
	alamat, err := uc.alamatRepo.GetByIDForUser(userID, in.AlamatKirim)
	if err != nil {
		return nil, err
	}
	if alamat == nil {
		return nil, ErrTrxAlamatNotFound
	}

	var (
		logs           []domain.LogProduk
		details        []domain.DetailTrx
		updatedProduct []*domain.Produk
		totalHarga     int
	)

	for _, item := range in.DetailTrx {
		if item.ProductID == 0 || item.Kuantitas <= 0 {
			return nil, errors.New("product_id dan kuantitas wajib diisi dan > 0")
		}

		produk, err := uc.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, err
		}
		if produk == nil {
			return nil, ErrTrxProductNotFound
		}

		hargaKonsumen, err := strconv.Atoi(produk.HargaKonsumen)
		if err != nil {
			return nil, fmt.Errorf("invalid harga_konsumen for product %d", produk.ID)
		}

		if produk.Stok < item.Kuantitas {
			return nil, ErrTrxInsufficientStock
		}

		lineTotal := hargaKonsumen * item.Kuantitas
		totalHarga += lineTotal

		// Prepare updated stock
		produk.Stok = produk.Stok - item.Kuantitas
		updatedProduct = append(updatedProduct, produk)

		// Prepare log_produk snapshot
		log := domain.LogProduk{
			ProdukID:      produk.ID,
			NamaProduk:    produk.NamaProduk,
			Slug:          produk.Slug,
			HargaReseller: produk.HargaReseller,
			HargaKonsumen: produk.HargaKonsumen,
			Deskripsi:     produk.Deskripsi,
			TokoID:        produk.TokoID,
			CategoryID:    produk.CategoryID,
		}
		logs = append(logs, log)

		detail := domain.DetailTrx{
			TokoID:     produk.TokoID,
			Kuantitas:  item.Kuantitas,
			HargaTotal: lineTotal,
		}
		details = append(details, detail)
	}

	trx := &domain.Trx{
		UserID:             userID,
		AlamatPengirimanID: alamat.ID,
		HargaTotal:         totalHarga,
		KodeInvoice:        fmt.Sprintf("INV-%d", time.Now().Unix()),
		MethodBayar:        in.MethodBayar,
	}

	if err := uc.trxRepo.CreateWithDetails(trx, logs, details, updatedProduct); err != nil {
		return nil, err
	}

	// Reload the trx with all relations for response
	created, err := uc.trxRepo.GetByIDForUser(userID, trx.ID)
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, ErrTrxNotFound
	}

	return created, nil
}
