package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/usecase"
)

// TokoHandler handles HTTP requests for toko resources.
type TokoHandler struct {
	tokoUC usecase.TokoUsecase
}

// NewTokoHandler creates a new TokoHandler.
func NewTokoHandler(tokoUC usecase.TokoUsecase) *TokoHandler {
	return &TokoHandler{tokoUC: tokoUC}
}

// GetAllToko handles GET /toko.
func (h *TokoHandler) GetAllToko(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	nama := c.Query("nama")

	result, err := h.tokoUC.GetAll(limit, page, nama)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	// Sesuai contoh Postman, data dikemas dalam object page+limit+data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"page":  result.Page,
			"limit": result.Limit,
			"data":  result.Data,
		},
	})
}

// GetTokoByID handles GET /toko/:id.
func (h *TokoHandler) GetTokoByID(c *fiber.Ctx) error {
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

	toko, err := h.tokoUC.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}
	if toko == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"id":        toko.ID,
			"nama_toko": toko.NamaToko,
			"url_foto":  toko.UrlFoto,
		},
	})
}

// GetMyToko handles GET /toko/my to get the store owned by logged-in user.
func (h *TokoHandler) GetMyToko(c *fiber.Ctx) error {
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

	toko, err := h.tokoUC.GetMyStore(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}
	if toko == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"id":        toko.ID,
			"nama_toko": toko.NamaToko,
			"url_foto":  toko.UrlFoto,
		},
	})
}

// UpdateMyToko handles PUT /toko (update store owned by logged-in user).
func (h *TokoHandler) UpdateMyToko(c *fiber.Ctx) error {
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

	var in usecase.UpdateTokoInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"invalid request body"},
			"data":    nil,
		})
	}

	toko, err := h.tokoUC.UpdateMyStore(userID, in)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}
	if toko == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to PUT data",
		"errors":  nil,
		"data": fiber.Map{
			"id":        toko.ID,
			"nama_toko": toko.NamaToko,
			"url_foto":  toko.UrlFoto,
		},
	})
}
