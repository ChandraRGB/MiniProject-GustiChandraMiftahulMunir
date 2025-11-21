package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/usecase"
)

// AlamatHandler handles HTTP requests related to alamat kirim.
type AlamatHandler struct {
	alamatUC usecase.AlamatUsecase
}

// NewAlamatHandler creates a new AlamatHandler.
func NewAlamatHandler(alamatUC usecase.AlamatUsecase) *AlamatHandler {
	return &AlamatHandler{alamatUC: alamatUC}
}

// GetMyAlamat handles GET /user/alamat.
func (h *AlamatHandler) GetMyAlamat(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  []string{"invalid user id in token"},
			"data":    nil,
		})
	}

	filterJudul := c.Query("judul_alamat")

	list, err := h.alamatUC.GetMyAlamat(userID, filterJudul)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	data := make([]fiber.Map, 0, len(list))
	for _, a := range list {
		data = append(data, fiber.Map{
			"id":            a.ID,
			"judul_alamat":  a.JudulAlamat,
			"nama_penerima": a.NamaPenerima,
			"no_telp":       a.NoTelp,
			"detail_alamat": a.DetailAlamat,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    data,
	})
}

// GetAlamatByID handles GET /user/alamat/:id.
func (h *AlamatHandler) GetAlamatByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
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

	alamat, err := h.alamatUC.GetByID(userID, uint(id))
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		var errs []string

		if errors.Is(err, usecase.ErrAlamatNotFound) {
			statusCode = fiber.StatusNotFound
			errs = append(errs, "record not found")
		} else {
			errs = append(errs, err.Error())
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  errs,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"id":            alamat.ID,
			"judul_alamat":  alamat.JudulAlamat,
			"nama_penerima": alamat.NamaPenerima,
			"no_telp":       alamat.NoTelp,
			"detail_alamat": alamat.DetailAlamat,
		},
	})
}

// CreateAlamat handles POST /user/alamat.
func (h *AlamatHandler) CreateAlamat(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  []string{"invalid user id in token"},
			"data":    nil,
		})
	}

	var in usecase.CreateAlamatInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"invalid request body"},
			"data":    nil,
		})
	}

	alamat, err := h.alamatUC.Create(userID, in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	// Sesuai contoh Postman, data mengembalikan id alamat
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    alamat.ID,
	})
}

// UpdateAlamat handles PUT /user/alamat/:id.
func (h *AlamatHandler) UpdateAlamat(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
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

	var in usecase.UpdateAlamatInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"invalid request body"},
			"data":    nil,
		})
	}

	alamat, err := h.alamatUC.Update(userID, uint(id), in)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		var errs []string

		if errors.Is(err, usecase.ErrAlamatNotFound) {
			statusCode = fiber.StatusNotFound
			errs = append(errs, "record not found")
		} else {
			errs = append(errs, err.Error())
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  errs,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to PUT data",
		"errors":  nil,
		"data": fiber.Map{
			"id":            alamat.ID,
			"judul_alamat":  alamat.JudulAlamat,
			"nama_penerima": alamat.NamaPenerima,
			"no_telp":       alamat.NoTelp,
			"detail_alamat": alamat.DetailAlamat,
		},
	})
}

// DeleteAlamat handles DELETE /user/alamat/:id.
func (h *AlamatHandler) DeleteAlamat(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
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

	if err := h.alamatUC.Delete(userID, uint(id)); err != nil {
		statusCode := fiber.StatusInternalServerError
		var errs []string

		if errors.Is(err, usecase.ErrAlamatNotFound) {
			statusCode = fiber.StatusNotFound
			errs = append(errs, "record not found")
		} else {
			errs = append(errs, err.Error())
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to DELETE data",
			"errors":  errs,
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
