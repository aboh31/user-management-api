package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-management-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(username, password, email string) (*model.UserResponse, error) {
	args := m.Called(username, password, email)
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserUsecase) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) GetProfile(userID uint) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUsecase) GetUsers() (*[]model.UserResponse, error) {
	args := m.Called()
	return args.Get(0).(*[]model.UserResponse), args.Error(1)
}

func TestLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.POST("/login", handler.Login)

	mockUsecase.On("Login", "descamp35", "password").Return("valid-token", nil)

	body, _ := json.Marshal(gin.H{"username": "descamp35", "password": "password"})
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestLoginUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.POST("/login", handler.Login)

	mockUsecase.On("Login", "invalid", "wrong").Return("", errors.New("unauthorized"))

	body, _ := json.Marshal(gin.H{"username": "invalid", "password": "wrong"})
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRegisterSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.POST("/register", handler.Register)

	// Setup mock return value (UserResponse)
	mockResponse := &model.UserResponse{
		Id:        "f921de9d-50c4-41be-aeb0-99b7980fd36f",
		Username:  "newuser",
		Email:     "new@example.com",
		CreatedAt: "2025-05-04T12:25:57",
	}

	// Mocking service
	mockUsecase.
		On("Register", "newuser", "password123", "new@example.com").
		Return(mockResponse, nil)

	// Prepare request body
	body, _ := json.Marshal(gin.H{
		"username": "newuser",
		"password": "password123",
		"email":    "new@example.com",
	})

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response body
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, mockResponse.Id, resp["id"])
	assert.Equal(t, mockResponse.Username, resp["username"])
	assert.Equal(t, mockResponse.Email, resp["email"])
	assert.Equal(t, mockResponse.CreatedAt, resp["created_at"])
}

func TestRegisterError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.POST("/register", handler.Register)

	// Setup mock to return an error
	mockUsecase.
		On("Register", "newuser", "password123", "new@example.com").
		Return((*model.UserResponse)(nil), errors.New("fail"))

	// Prepare request body
	body, _ := json.Marshal(gin.H{
		"username": "newuser",
		"password": "password123",
		"email":    "new@example.com",
	})
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert HTTP status
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Parse response body and assert error message
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Contains(t, resp, "error")
	assert.Equal(t, "fail", resp["error"])
}

func TestProfileSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.GET("/profile", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		handler.Profile(c)
	})

	mockUsecase.On("GetProfile", uint(1)).Return(&model.User{ID: 1}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProfileNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.GET("/profile", func(c *gin.Context) {
		c.Set("user_id", uint(2))
		handler.Profile(c)
	})

	mockUsecase.On("GetProfile", uint(2)).Return(&model.User{}, errors.New("not found"))

	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUsersSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.GET("/users", handler.Users)

	mockUsers := []model.UserResponse{
		{
			Id:        "uuid-1",
			Username:  "user1",
			Email:     "user1@example.com",
			CreatedAt: "2025-05-04T12:00:00",
		},
		{
			Id:        "uuid-2",
			Username:  "user2",
			Email:     "user2@example.com",
			CreatedAt: "2025-05-04T12:01:00",
		},
	}

	mockUsecase.On("GetUsers").Return(&mockUsers, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response JSON to map["users"] -> []model.UserResponse
	var response map[string][]model.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	users, ok := response["users"]
	assert.True(t, ok)
	assert.Len(t, users, 2)

	assert.Equal(t, "uuid-1", users[0].Id)
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "user1@example.com", users[0].Email)

	assert.Equal(t, "uuid-2", users[1].Id)
	assert.Equal(t, "user2", users[1].Username)
	assert.Equal(t, "user2@example.com", users[1].Email)
}

func TestLoginInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.POST("/login", handler.Login)

	body := []byte(`{invalid_json}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.POST("/register", handler.Register)

	body := []byte(`{invalid_json}`)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUsersError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	router := gin.Default()
	router.GET("/users", handler.Users)

	mockUsecase.On("GetUsers").Return((*[]model.UserResponse)(nil), errors.New("db error"))

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "db error", resp["error"])
}
