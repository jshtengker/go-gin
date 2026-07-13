package main

import (
	"backend/internal/api/requests"
	"backend/internal/api/router"
	"backend/internal/configs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func getDataParam(c *gin.Context) {
	name := c.Param("name")
	age := c.Param("age")
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"name":    name,
		"age":     age,
	})
}

func getDataQuery(c *gin.Context) {
	name := c.Query("name")
	age := c.Query("age")
	c.JSON(http.StatusOK, requests.APIResponse{
		Success: true,
		Message: "user found",
		Data: requests.User{
			Name: name,
			Age:  age,
		},
	})
}

func main() {

	configs.LoadEnv()

	cfg := configs.Load()

	db, err := configs.ConnPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	log.Println("connected to PostgreSQL")

	r := router.Setup(cfg, db)

	log.Println("Server started on :9001")

	if err := r.Run(":9001"); err != nil {
		log.Fatal(err)
	}
}
