package handlers

import (
	"backend/internal/api/responses"
	"backend/internal/mocks"
	"backend/pkg/apperrors"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestProductRouter(mockService *mocks.MockProductService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewProductHandler(mockService)

	r.GET("/products", handler.GetAll)
	r.GET("/products/:id", handler.GetByID)
	return r
}

func TestProductHandler_GetAll(t *testing.T) {
	mockService := new(mocks.MockProductService)
	router := setupTestProductRouter(mockService)

	t.Run("Success", func(t *testing.T) {
		products := []responses.ProductResponse{
			{ID: uuid.New(), Name: "Widget A", Price: "10.00"},
		}
		mockService.On("GetAll", mock.Anything, mock.Anything).Return(products, int64(1), nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/products?page=1&page_size=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		assert.False(t, res["error"].(bool))
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService.On("GetAll", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("query error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestProductHandler_GetByID(t *testing.T) {
	mockService := new(mocks.MockProductService)
	router := setupTestProductRouter(mockService)
	validID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		prod := &responses.ProductResponse{ID: uuid.MustParse(validID), Name: "Widget A"}
		mockService.On("GetByID", mock.Anything, validID).Return(prod, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/products/"+validID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/products/abc-not-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, validID).Return(nil, apperrors.ErrNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/products/"+validID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
