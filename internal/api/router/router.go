package router

import (
	"backend/internal/api/routes"
	"backend/internal/configs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(cfg *configs.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	routes.Register(r, cfg, db)

	return r
}
