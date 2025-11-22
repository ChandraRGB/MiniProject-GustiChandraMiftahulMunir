package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/usecase"
)

// ProvinceCityHandler handles HTTP requests for province and city resources.
type ProvinceCityHandler struct {
	uc usecase.ProvinceCityUsecase
}

// NewProvinceCityHandler creates a new ProvinceCityHandler.
func NewProvinceCityHandler(uc usecase.ProvinceCityUsecase) *ProvinceCityHandler {
	return &ProvinceCityHandler{uc: uc}
}

// GetListProvince handles GET /provcity/listprovincies.
func (h *ProvinceCityHandler) GetListProvince(c *fiber.Ctx) error {
	provinces, err := h.uc.GetProvinces()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to get data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to get data",
		"errors":  nil,
		"data":    provinces,
	})
}

// GetListCities handles GET /provcity/listcities/:prov_id.
func (h *ProvinceCityHandler) GetListCities(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	if provID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to get data",
			"errors":  []string{"prov_id is required"},
			"data":    nil,
		})
	}

	cities, err := h.uc.GetCitiesByProvince(provID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to get data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to get data",
		"errors":  nil,
		"data":    cities,
	})
}
