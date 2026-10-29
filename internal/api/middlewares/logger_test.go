package middlewares_test

import (
	"backend/internal/api/middlewares"
	"backend/pkg/logger"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJSONLogger_Development(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init("development")

	r := gin.New()
	r.Use(middlewares.RequestID())
	r.Use(middlewares.JSONLogger())

	r.GET("/test-ok", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/test-warn", func(c *gin.Context) {
		_ = c.Error(errors.New("bad request error"))
		c.String(http.StatusBadRequest, "bad request")
	})
	r.GET("/test-err", func(c *gin.Context) {
		c.String(http.StatusInternalServerError, "server error")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test-ok?param=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	reqWarn, _ := http.NewRequest(http.MethodGet, "/test-warn", nil)
	wWarn := httptest.NewRecorder()
	r.ServeHTTP(wWarn, reqWarn)
	assert.Equal(t, http.StatusBadRequest, wWarn.Code)

	reqErr, _ := http.NewRequest(http.MethodGet, "/test-err", nil)
	wErr := httptest.NewRecorder()
	r.ServeHTTP(wErr, reqErr)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)
}

func TestJSONLogger_Production(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init("production")

	r := gin.New()
	r.Use(middlewares.RequestID())
	r.Use(middlewares.JSONLogger())

	r.GET("/test-prod", func(c *gin.Context) {
		c.String(http.StatusOK, "prod ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test-prod", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
