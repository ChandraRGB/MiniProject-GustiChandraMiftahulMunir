package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/usecase"
)

// TrxHandler handles HTTP requests for transaction resources.
type TrxHandler struct {
	trxUC usecase.TrxUsecase
}

// NewTrxHandler creates a new TrxHandler.
func NewTrxHandler(trxUC usecase.TrxUsecase) *TrxHandler {
	return &TrxHandler{trxUC: trxUC}
}

type postTrxRequest struct {
	MethodBayar string `json:"method_bayar"`
	AlamatKirim uint   `json:"alamat_kirim"`
	DetailTrx   []struct {
		ProductID uint `json:"product_id"`
		Kuantitas int  `json:"kuantitas"`
	} `json:"detail_trx"`
}

// GetAllTrx handles GET /trx.
func (h *TrxHandler) GetAllTrx(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  []string{"invalid user id in token"},
			"data":    nil,
		})
	}

	trxs, err := h.trxUC.GetAll(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	list := make([]fiber.Map, 0, len(trxs))
	for i := range trxs {
		list = append(list, buildTrxResponse(&trxs[i]))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"data":  list,
			"page":  0,
			"limit": 0,
		},
	})
}

// GetTrxByID handles GET /trx/:id.
func (h *TrxHandler) GetTrxByID(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  []string{"invalid user id in token"},
			"data":    nil,
		})
	}

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"invalid id"},
			"data":    nil,
		})
	}

	trx, err := h.trxUC.GetByID(userID, uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrTrxNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to GET data",
				"errors":  []string{"No Data Trx"},
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    buildTrxResponse(trx),
	})
}

// PostTrx handles POST /trx.
func (h *TrxHandler) PostTrx(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  []string{"invalid user id in token"},
			"data":    nil,
		})
	}

	var req postTrxRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"invalid request body"},
			"data":    nil,
		})
	}

	in := usecase.CreateTrxInput{
		MethodBayar: req.MethodBayar,
		AlamatKirim: req.AlamatKirim,
	}
	for _, d := range req.DetailTrx {
		in.DetailTrx = append(in.DetailTrx, usecase.TrxItemInput{
			ProductID: d.ProductID,
			Kuantitas: d.Kuantitas,
		})
	}

	trx, err := h.trxUC.Create(userID, in)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		var errs []string

		if errors.Is(err, usecase.ErrTrxAlamatNotFound) {
			statusCode = fiber.StatusBadRequest
			errs = append(errs, "alamat_kirim tidak ditemukan")
		} else if errors.Is(err, usecase.ErrTrxProductNotFound) {
			statusCode = fiber.StatusBadRequest
			errs = append(errs, "product tidak ditemukan")
		} else if errors.Is(err, usecase.ErrTrxInsufficientStock) {
			statusCode = fiber.StatusBadRequest
			errs = append(errs, "stok tidak cukup")
		} else if errors.Is(err, usecase.ErrTrxEmptyDetail) {
			statusCode = fiber.StatusBadRequest
			errs = append(errs, "detail_trx tidak boleh kosong")
		} else {
			errs = append(errs, err.Error())
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  errs,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    trx.ID,
	})
}

// buildTrxResponse maps domain.Trx to the JSON structure used in the Postman examples.
func buildTrxResponse(trx *domain.Trx) fiber.Map {
	// Build alamat_kirim object
	alamat := fiber.Map{
		"id":             trx.Alamat.ID,
		"judul_alamat":   trx.Alamat.JudulAlamat,
		"nama_penerima":  trx.Alamat.NamaPenerima,
		"no_telp":        trx.Alamat.NoTelp,
		"detail_alamat":  trx.Alamat.DetailAlamat,
	}

	// Build detail_trx list
	details := make([]fiber.Map, 0, len(trx.DetailTrx))
	for _, d := range trx.DetailTrx {
		productMap := buildProductFromLog(&d.LogProduk)

		tokoMap := fiber.Map{
			"id":        d.Toko.ID,
			"nama_toko": d.Toko.NamaToko,
			"url_foto":  d.Toko.UrlFoto,
		}

		details = append(details, fiber.Map{
			"product":     productMap,
			"toko":        tokoMap,
			"kuantitas":   d.Kuantitas,
			"harga_total": d.HargaTotal,
		})
	}

	return fiber.Map{
		"id":           trx.ID,
		"harga_total":  trx.HargaTotal,
		"kode_invoice": trx.KodeInvoice,
		"method_bayar": trx.MethodBayar,
		"alamat_kirim": alamat,
		"detail_trx":   details,
	}
}

// buildProductFromLog maps a LogProduk snapshot into the Product JSON shape used in responses.
func buildProductFromLog(log *domain.LogProduk) fiber.Map {
	hargaReseller, _ := strconv.Atoi(log.HargaReseller)
	hargaKonsumen, _ := strconv.Atoi(log.HargaKonsumen)

	photos := make([]fiber.Map, 0, len(log.Produk.FotoProduk))
	for _, p := range log.Produk.FotoProduk {
		photos = append(photos, fiber.Map{
			"id":         p.ID,
			"product_id": p.ProdukID,
			"url":        p.URL,
		})
	}

	tokoMap := fiber.Map{
		"nama_toko": log.Toko.NamaToko,
		"url_foto":  log.Toko.UrlFoto,
	}

	categoryMap := fiber.Map{
		"id":            log.Category.ID,
		"nama_category": log.Category.Nama,
	}

	return fiber.Map{
		"id":             log.ProdukID,
		"nama_produk":    log.NamaProduk,
		"slug":           log.Slug,
		"harga_reseler":  hargaReseller,
		"harga_konsumen": hargaKonsumen,
		"deskripsi":      log.Deskripsi,
		"toko":           tokoMap,
		"category":       categoryMap,
		"photos":         photos,
	}
}
