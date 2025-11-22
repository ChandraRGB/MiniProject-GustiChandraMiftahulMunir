package http

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/middleware"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/repository"
	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/usecase"
)

// RegisterRoutes registers all HTTP routes for the application.
func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tokoRepo := repository.NewTokoRepository(db)
	alamatRepo := repository.NewAlamatRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	fotoProdukRepo := repository.NewFotoProdukRepository(db)
	trxRepo := repository.NewTrxRepository(db)

	// Initialize usecases
	authUC := usecase.NewAuthUsecase(userRepo, tokoRepo)
	userUC := usecase.NewUserUsecase(userRepo)
	alamatUC := usecase.NewAlamatUsecase(alamatRepo)
	tokoUC := usecase.NewTokoUsecase(tokoRepo)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)
	productUC := usecase.NewProductUsecase(productRepo, fotoProdukRepo, tokoRepo)
	trxUC := usecase.NewTrxUsecase(trxRepo, alamatRepo, productRepo)
	provinceCityUC := usecase.NewProvinceCityUsecase()

	// Initialize handlers
	authHandler := NewAuthHandler(authUC)
	userHandler := NewUserHandler(userUC)
	alamatHandler := NewAlamatHandler(alamatUC)
	tokoHandler := NewTokoHandler(tokoUC)
	categoryHandler := NewCategoryHandler(categoryUC)
	productHandler := NewProductHandler(productUC)
	trxHandler := NewTrxHandler(trxUC)
	provinceCityHandler := NewProvinceCityHandler(provinceCityUC)

	// Auth routes based on Postman collection
	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	// User routes (protected with JWT middleware)
	userGroup := app.Group("/user", middleware.JWTMiddleware())
	userGroup.Get("/", userHandler.GetProfile)
	userGroup.Put("/", userHandler.UpdateProfile)

	alamatGroup := userGroup.Group("/alamat")
	alamatGroup.Get("/", alamatHandler.GetMyAlamat)
	alamatGroup.Get("/:id", alamatHandler.GetAlamatByID)
	alamatGroup.Post("/", alamatHandler.CreateAlamat)
	alamatGroup.Put("/:id", alamatHandler.UpdateAlamat)
	alamatGroup.Delete("/:id", alamatHandler.DeleteAlamat)

	// Toko routes (public listing/detail, and update for logged-in user)
	app.Get("/toko", tokoHandler.GetAllToko)
	app.Get("/toko/my", middleware.JWTMiddleware(), tokoHandler.GetMyToko)
	app.Get("/toko/:id", tokoHandler.GetTokoByID)
	// Support both PUT /toko and PUT /toko/:id_toko (as in Postman collection)
	app.Put("/toko", middleware.JWTMiddleware(), tokoHandler.UpdateMyToko)
	app.Put("/toko/:id_toko", middleware.JWTMiddleware(), tokoHandler.UpdateMyToko)

	// Category routes (admin only)
	categoryGroup := app.Group("/category", middleware.JWTMiddleware(), middleware.AdminOnly())
	categoryGroup.Get("/", categoryHandler.GetAll)
	categoryGroup.Get("/:id", categoryHandler.GetByID)
	categoryGroup.Post("/", categoryHandler.Create)
	categoryGroup.Put("/:id", categoryHandler.Update)
	categoryGroup.Delete("/:id", categoryHandler.Delete)

	// Product routes
	app.Get("/product", productHandler.GetAllProduct)
	app.Get("/product/:id", productHandler.GetProductByID)
	productGroup := app.Group("/product", middleware.JWTMiddleware())
	productGroup.Post("/", productHandler.CreateProduct)
	productGroup.Put("/:id", productHandler.UpdateProduct)
	productGroup.Delete("/:id", productHandler.DeleteProduct)

	// Trx routes (protected with JWT middleware)
	trxGroup := app.Group("/trx", middleware.JWTMiddleware())
	trxGroup.Get("/", trxHandler.GetAllTrx)
	trxGroup.Get("/:id", trxHandler.GetTrxByID)
	trxGroup.Post("/", trxHandler.PostTrx)

	// Province & City routes (public, proxy to EMSIFA API)
	provCityGroup := app.Group("/provcity")
	provCityGroup.Get("/listprovincies", provinceCityHandler.GetListProvince)
	provCityGroup.Get("/listcities/:prov_id", provinceCityHandler.GetListCities)

	// TODO: register other feature routes (user, toko, alamat, kategori, produk, trx)
}