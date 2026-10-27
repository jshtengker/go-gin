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

func setupTestUserRouter(mockService *mocks.MockUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewUserHandler(mockService)

	r.GET("/users", handler.GetAll)
	r.GET("/users/:id", handler.GetByID)
	return r
}

func TestUserHandler_GetAll(t *testing.T) {
	mockService := new(mocks.MockUserService)
	router := setupTestUserRouter(mockService)

	t.Run("Success", func(t *testing.T) {
		users := []responses.UserResponse{
			{ID: uuid.New(), FullName: "John", Email: "john@example.com"},
		}
		mockService.On("GetAll", mock.Anything, mock.Anything).Return(users, int64(1), nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users?page=1&page_size=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		assert.False(t, res["error"].(bool))
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService.On("GetAll", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("service failure")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHandler_GetByID(t *testing.T) {
	mockService := new(mocks.MockUserService)
	router := setupTestUserRouter(mockService)
	validID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		user := &responses.UserResponse{ID: uuid.MustParse(validID), FullName: "John"}
		mockService.On("GetByID", mock.Anything, validID).Return(user, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users/"+validID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, validID).Return(nil, apperrors.ErrNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users/"+validID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
