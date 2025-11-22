package http

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/usecase"
)

// ProductHandler handles HTTP requests for product resources.
type ProductHandler struct {
	productUC usecase.ProductUsecase
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(productUC usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUC: productUC}
}

// GetAllProduct handles GET /product.
func (h *ProductHandler) GetAllProduct(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	filter := usecase.ProductFilter{}
	filter.NamaProduk = c.Query("nama_produk")

	if v := c.Query("category_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filter.CategoryID = uint(id)
		}
	}
	if v := c.Query("toko_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filter.TokoID = uint(id)
		}
	}
	if v := c.Query("min_harga"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.MinHarga = n
		}
	}
	if v := c.Query("max_harga"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.MaxHarga = n
		}
	}

	result, err := h.productUC.GetAll(limit, page, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	products := make([]fiber.Map, 0, len(result.Data))
	for i := range result.Data {
		products = append(products, buildProductResponse(&result.Data[i]))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"data":  products,
			"page":  result.Page,
			"limit": result.Limit,
		},
	})
}

// GetProductByID handles GET /product/:id.
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
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

	product, err := h.productUC.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to GET data",
				"errors":  []string{"No Data Product"},
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
		"data":    buildProductResponse(product),
	})
}

// CreateProduct handles POST /product.
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
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

	namaProduk := c.FormValue("nama_produk")
	categoryIDStr := c.FormValue("category_id")
	hargaResellerStr := c.FormValue("harga_reseller")
	hargaKonsumenStr := c.FormValue("harga_konsumen")
	stokStr := c.FormValue("stok")
	deskripsi := c.FormValue("deskripsi")

	if namaProduk == "" || categoryIDStr == "" || hargaResellerStr == "" || hargaKonsumenStr == "" || stokStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"nama_produk, category_id, harga_reseller, harga_konsumen, stok wajib diisi"},
			"data":    nil,
		})
	}

	categoryIDInt, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"invalid category_id"},
			"data":    nil,
		})
	}
	hargaReseller, err := strconv.Atoi(hargaResellerStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"invalid harga_reseller"},
			"data":    nil,
		})
	}
	hargaKonsumen, err := strconv.Atoi(hargaKonsumenStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"invalid harga_konsumen"},
			"data":    nil,
		})
	}
	stok, err := strconv.Atoi(stokStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"invalid stok"},
			"data":    nil,
		})
	}

	photoFilenames, err := saveProductPhotos(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	in := usecase.CreateProductInput{
		NamaProduk:    namaProduk,
		CategoryID:    uint(categoryIDInt),
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok:          stok,
		Deskripsi:     deskripsi,
	}

	product, err := h.productUC.Create(userID, in, photoFilenames)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    product.ID,
	})
}

// UpdateProduct handles PUT /product/:id.
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
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
			"message": "Failed to PUT data",
			"errors":  []string{"invalid id"},
			"data":    nil,
		})
	}

	var in usecase.UpdateProductInput

	if v := c.FormValue("nama_produk"); v != "" {
		in.NamaProduk = &v
	}
	if v := c.FormValue("category_id"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			vv := uint(val)
			in.CategoryID = &vv
		}
	}
	if v := c.FormValue("harga_reseller"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			in.HargaReseller = &val
		}
	}
	if v := c.FormValue("harga_konsumen"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			in.HargaKonsumen = &val
		}
	}
	if v := c.FormValue("stok"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			in.Stok = &val
		}
	}
	if v := c.FormValue("deskripsi"); v != "" {
		in.Deskripsi = &v
	}

	photoFilenames, err := saveProductPhotos(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	product, err := h.productUC.Update(userID, uint(id), in, photoFilenames)
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to PUT data",
				"errors":  []string{"record not found"},
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	_ = product

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to PUT data",
		"errors":  nil,
		"data":    "",
	})
}

// DeleteProduct handles DELETE /product/:id.
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
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
			"message": "Failed to DELETE data",
			"errors":  []string{"invalid id"},
			"data":    nil,
		})
	}

	if err := h.productUC.Delete(userID, uint(id)); err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to DELETE data",
				"errors":  []string{"record not found"},
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to DELETE data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to DELETE data",
		"errors":  nil,
		"data":    "",
	})
}

// saveProductPhotos saves uploaded files under the "photos" field
// and returns their stored filenames.
func saveProductPhotos(c *fiber.Ctx) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		// Not a multipart request or no files - treat as no photos.
		return []string{}, nil
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return []string{}, nil
	}

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return nil, err
	}

	var filenames []string
	for _, fh := range files {
		if fh == nil {
			continue
		}
		filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(fh.Filename))
		fullpath := filepath.Join("uploads", filename)
		if err := c.SaveFile(fh, fullpath); err != nil {
			return nil, err
		}
		filenames = append(filenames, filename)
	}

	return filenames, nil
}

// buildProductResponse maps domain.Produk into the JSON shape used by Postman.
func buildProductResponse(p *domain.Produk) fiber.Map {
	hargaReseller, _ := strconv.Atoi(p.HargaReseller)
	hargaKonsumen, _ := strconv.Atoi(p.HargaKonsumen)

	photos := make([]fiber.Map, 0, len(p.FotoProduk))
	for _, photo := range p.FotoProduk {
		photos = append(photos, fiber.Map{
			"id":         photo.ID,
			"product_id": photo.ProdukID,
			"url":        photo.URL,
		})
	}

	return fiber.Map{
		"id":             p.ID,
		"nama_produk":    p.NamaProduk,
		"slug":           p.Slug,
		"harga_reseler":  hargaReseller,
		"harga_konsumen": hargaKonsumen,
		"stok":           p.Stok,
		"deskripsi":      p.Deskripsi,
		"toko": fiber.Map{
			"id":        p.Toko.ID,
			"nama_toko": p.Toko.NamaToko,
			"url_foto":  p.Toko.UrlFoto,
		},
		"category": fiber.Map{
			"id":            p.Category.ID,
			"nama_category": p.Category.Nama,
		},
		"photos": photos,
	}
}
