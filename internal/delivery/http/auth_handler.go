package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/usecase"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	authUC usecase.AuthUsecase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authUC usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var in usecase.RegisterInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"invalid request body"},
			"data":    nil,
		})
	}

	if err := h.authUC.Register(in); err != nil {
		statusCode := fiber.StatusInternalServerError
		var errs []string

		if errors.Is(err, usecase.ErrEmailOrPhoneExists) {
			statusCode = fiber.StatusBadRequest
			errs = append(errs, "email atau no_telp sudah digunakan")
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
		"data":    "Register Succeed",
	})
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var in usecase.LoginInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"invalid request body"},
			"data":    nil,
		})
	}

	res, err := h.authUC.Login(in)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		var errs []string

		if errors.Is(err, usecase.ErrInvalidCredentials) {
			statusCode = fiber.StatusUnauthorized
			errs = append(errs, "No Telp atau kata sandi salah")
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

	// response: token + sebagian data user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data": fiber.Map{
			"nama":          res.User.Nama,
			"no_telp":       res.User.NoTelp,
			"tanggal_Lahir": res.User.TanggalLahir,
			"tentang":       res.User.Tentang,
			"pekerjaan":     res.User.Pekerjaan,
			"email":         res.User.Email,
			"id_provinsi":   res.User.IDProvinsi,
			"id_kota":       res.User.IDKota,
			"token":         res.Token,
		},
	})
}