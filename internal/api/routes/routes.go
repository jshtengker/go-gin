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
	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewRefreshTokenRepository(db)
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenRepo, cfg)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo)

	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)

	authMiddleware := middlewares.Authenticate(cfg.JWTSecret)

	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authMiddleware, authHandler.Logout)
		auth.GET("/me", authMiddleware, authHandler.GetProfile)
	}

	users := v1.Group("/users", authMiddleware)
	{
		users.GET("", userHandler.GetAll)
		users.GET("/:id", userHandler.GetByID)
	}

	products := v1.Group("/products", authMiddleware)
	{
		products.GET("", productHandler.GetAll)
		products.GET("/:id", productHandler.GetByID)
	}

	orders := v1.Group("/orders", authMiddleware)
	{
		orders.GET("", orderHandler.GetAll)
		orders.GET("/:id", orderHandler.GetByID)
	}
}
