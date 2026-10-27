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

func setupTestOrderRouter(mockService *mocks.MockOrderService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewOrderHandler(mockService)

	r.GET("/orders", handler.GetAll)
	r.GET("/orders/:id", handler.GetByID)
	return r
}

func TestOrderHandler_GetAll(t *testing.T) {
	mockService := new(mocks.MockOrderService)
	router := setupTestOrderRouter(mockService)

	t.Run("Success", func(t *testing.T) {
		orders := []responses.OrderListResponse{
			{ID: uuid.New(), OrderNumber: "ORD-999", TotalAmount: 150.00},
		}
		mockService.On("GetAll", mock.Anything, mock.Anything).Return(orders, int64(1), nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/orders?page=1&page_size=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		assert.False(t, res["error"].(bool))
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService.On("GetAll", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("db failure")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestOrderHandler_GetByID(t *testing.T) {
	mockService := new(mocks.MockOrderService)
	router := setupTestOrderRouter(mockService)
	validID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		ord := &responses.OrderResponse{ID: uuid.MustParse(validID), OrderNumber: "ORD-999"}
		mockService.On("GetByID", mock.Anything, validID).Return(ord, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/orders/"+validID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/orders/non-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, validID).Return(nil, apperrors.ErrNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/orders/"+validID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
