package router

import (
	"backend/internal/api/middlewares"
	"backend/internal/api/routes"
	"backend/internal/configs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(cfg *configs.Config, db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Global Middlewares
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestID())
	r.Use(middlewares.JSONLogger())

	api := r.Group("/api")
	routes.Register(api, cfg, db)

	return r
}
