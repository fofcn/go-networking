// user_test.go

package user

import (
	"bytes"
	"encoding/json"
	"go-networking/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	// Initialize test router
	router := gin.Default()
	router.POST("/auth/login", db.WithDB(Login))

	// Create mock request body
	mockRequestBody := LoginCmd{
		Username: "test",
		Password: "password",
	}
	mockRequestJSON, _ := json.Marshal(mockRequestBody)

	// Perform mock request
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(mockRequestJSON))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validate response
	if w.Code != http.StatusOK {
		t.Errorf("Expected response code %d but got %d", http.StatusOK, w.Code)
	}

	// TODO: Add more validation for response data
}

func TestRegister(t *testing.T) {
	// Initialize test router
	router := gin.Default()
	router.POST("/auth/register", db.WithDB(Register))

	// Create mock request body
	mockRequestBody := RegisterCmd{
		Username: "test",
		Password: "password",
	}
	mockRequestJSON, _ := json.Marshal(mockRequestBody)

	// Perform mock request
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(mockRequestJSON))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validate response
	if w.Code != http.StatusOK {
		t.Errorf("Expected response code %d but got %d", http.StatusOK, w.Code)
	}

	// TODO: Add more validation for response data
}

func TestGetUserInfo(t *testing.T) {
	// Initialize test router
	router := gin.Default()
	router.GET("/auth/userinfo", db.WithDB(GetUserInfo))

	// Create mock request with custom claims
	mockClaims := &CustomClaims{
		Username: "test",
		UserId:   1,
	}
	router.Use(func(c *gin.Context) {
		c.Set("claims", mockClaims)
	})

	// Perform mock request
	req, _ := http.NewRequest("GET", "/auth/userinfo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validate response
	if w.Code != http.StatusOK {
		t.Errorf("Expected response code %d but got %d", http.StatusOK, w.Code)
	}

	// TODO: Add more validation for response data
}
