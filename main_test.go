package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestGetIndexes(t *testing.T) {
	// r := SetUpRouter()
	// mockResponse := `{"message":"Welcome to the Tech Company listing API with Golang"}`
	req, _ := http.NewRequest("GET", ":8080/v1/products", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	fmt.Println("fffffffffffffff", string(responseData))
	// assert.Equal(t, mockResponse, string(responseData))
	// assert.Equal(t, http.StatusOK, w.Code)
}
