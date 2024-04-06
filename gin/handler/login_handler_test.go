package handler_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"go-networking/gin/common"
// 	"go-networking/gin/handler"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// func TestLogin(t *testing.T) {

// 	gin.SetMode(gin.TestMode)

// 	tests := []struct {
// 		name     string
// 		input    handler.Loginrequest
// 		expected string
// 	}{
// 		{"valid credentials", handler.Loginrequest{"hello", "world"}, "OK"},
// 		{"invalid credentials", handler.Loginrequest{"wrong", "credentials"}, "Failed"},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			router := gin.Default()
// 			router.POST("/login", handler.Login)

// 			body, _ := json.Marshal(tt.input)
// 			req, err := http.NewRequest("POST", "/login", bytes.NewReader(body))
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			resp := httptest.NewRecorder()
// 			router.ServeHTTP(resp, req)

// 			var got common.SuccessResponse
// 			if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
// 				t.Fatal(err)
// 			}

// 			assert.Equal(t, tt.expected, got.Data)
// 		})
// 	}
// }
