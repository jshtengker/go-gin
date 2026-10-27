package routes

import (
	"backend/internal/api/handlers"
	"backend/internal/api/middlewares"
	"backend/internal/configs"
	"backend/internal/repositories"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(api *gin.RouterGroup, cfg *configs.Config, db *gorm.DB) {

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewRefreshTokenRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenRepo, cfg)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Auth Middleware
	authMiddleware := middlewares.Authenticate(cfg.JWTSecret)

	// API Versioning
	v1 := api.Group("/v1")

	// Authentication Routes
	auth := v1.Group("/auth")
	{
		// Public auth endpoints
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)

		// Protected auth endpoints
		auth.POST("/logout", authMiddleware, authHandler.Logout)
		auth.GET("/me", authMiddleware, authHandler.GetProfile)
	}

	// User Routes (Protected)
	users := v1.Group("/users", authMiddleware)
	{
		users.GET("", userHandler.GetAll)
		users.GET("/:id", userHandler.GetByID)
	}

	// Product Routes (Protected)
	products := v1.Group("/products", authMiddleware)
	{
		products.GET("", productHandler.GetAll)
		products.GET("/:id", productHandler.GetByID)
	}

	// Order Routes (Protected)
	orders := v1.Group("/orders", authMiddleware)
	{
		orders.GET("", orderHandler.GetAll)
		orders.GET("/:id", orderHandler.GetByID)
	}
}
