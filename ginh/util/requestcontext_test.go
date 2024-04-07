package util_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestGetUserIdAndUsername is a unit test for GetUserIdAndUsername function.
func TestGetUserIdAndUsername(t *testing.T) {
	// Create a new gin context with mock claims
	mockClaims := &util.CustomClaims{
		UserId:   123,
		Username: "testuser",
	}
	mockContext := &gin.Context{
		Request: &http.Request{},
		Params:  []gin.Param{},
		Data:    map[string]interface{}{},
	}
	mockContext.Set("claims", mockClaims)

	// Call the GetUserIdAndUsername function
	userId, username := GetUserIdAndUsername(mockContext)

	// Assert the returned values
	if userId != mockClaims.UserId {
		t.Errorf("Expected userId to be %d, got %d", mockClaims.UserId, userId)
	}
	if username != mockClaims.Username {
		t.Errorf("Expected username to be %s, got %s", mockClaims.Username, username)
	}
}
