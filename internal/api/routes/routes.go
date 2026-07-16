package routes

import (
	"backend/internal/api/handlers"
	"backend/internal/configs"
	"backend/internal/repositories"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, cfg *configs.Config, db *gorm.DB) {

	// Dependencies
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	orderRepo := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	// API Versioning
	v1 := r.Group("/v1")

	// User Routes
	users := v1.Group("/users")
	{
		users.GET("", userHandler.GetAll)
		users.GET("/:id", userHandler.GetById)
	}

	products := v1.Group("/products")
	{
		products.GET("", productHandler.GetAll)
		products.GET("/:id", productHandler.GetById)

	}

	orders := v1.Group("/orders")
	{
		orders.GET("", orderHandler.GetAll)
		orders.GET("/:id", orderHandler.GetById)
	}
}
