package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/usecase"
)

// UserHandler handles HTTP requests related to the authenticated user.
type UserHandler struct {
	userUC usecase.UserUsecase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userUC usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

// GetProfile handles GET /user.
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  []string{"invalid user id in token"},
			"data":    nil,
		})
	}

	user, err := h.userUC.GetProfile(userID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		var errs []string

		if errors.Is(err, usecase.ErrUserNotFound) {
			statusCode = fiber.StatusNotFound
			errs = append(errs, "user not found")
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
			"id":            user.ID,
			"nama":          user.Nama,
			"no_telp":       user.NoTelp,
			"tanggal_Lahir": user.TanggalLahir,
			"tentang":       user.Tentang,
			"pekerjaan":     user.Pekerjaan,
			"email":         user.Email,
			"id_provinsi":   user.IDProvinsi,
			"id_kota":       user.IDKota,
		},
	})
}

// UpdateProfile handles PUT /user.
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Unauthorized",
			"errors":  []string{"invalid user id in token"},
			"data":    nil,
		})
	}

	var in usecase.UpdateUserInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"invalid request body"},
			"data":    nil,
		})
	}

	user, err := h.userUC.UpdateProfile(userID, in)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		var errs []string

		if errors.Is(err, usecase.ErrUserNotFound) {
			statusCode = fiber.StatusNotFound
			errs = append(errs, "user not found")
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
			"id":            user.ID,
			"nama":          user.Nama,
			"no_telp":       user.NoTelp,
			"tanggal_Lahir": user.TanggalLahir,
			"tentang":       user.Tentang,
			"pekerjaan":     user.Pekerjaan,
			"email":         user.Email,
			"id_provinsi":   user.IDProvinsi,
			"id_kota":       user.IDKota,
		},
	})
}
