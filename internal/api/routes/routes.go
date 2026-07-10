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

	// API Versioning
	v1 := r.Group("/v1")

	// User Routes
	users := v1.Group("/users")
	{
		users.GET("", userHandler.GetAll)
	}
}
