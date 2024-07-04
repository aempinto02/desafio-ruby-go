package controller_test

import (
	"encoding/json"
	"ms-go/app/controller"
	"ms-go/app/models"
	"ms-go/app/services/kafka/producers"
	"ms-go/app/services/products"
	mockdb "ms-go/app/test/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize mock database and set UseMock to true
	products.SetUseMock(true)
	mockdb.MockData = []models.Product{
		{ID: 1, Name: "Product 1", Brand: "Brand A", Price: 100.0, Description: "Desc 1", Stock: 10},
		{ID: 2, Name: "Product 2", Brand: "Brand B", Price: 200.0, Description: "Desc 2", Stock: 20},
	}

	// Disable Kafka producer for tests
	producers.EnableKafkaProducer = false
}

func TestIndexProducts(t *testing.T) {

	// Setup a mock Gin context
	router := gin.Default()
	router.GET("/products", controller.IndexProducts) // Import and use IndexProducts from the controller package

	// Create a mock HTTP request to /products endpoint
	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the HTTP request to the response recorder
	router.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Ensure the expected data is present in the response
	data, ok := response["data"].([]interface{})
	assert.True(t, ok, "Expected 'data' field in response, got none")
	assert.Len(t, data, 2, "Expected 2 products, got %d", len(data))
}

func TestShowProducts(t *testing.T) {
	// Set up mock Gin context with a mock product ID
	router := gin.Default()
	router.GET("/products/:id", controller.ShowProducts)

	// Mock HTTP request to /products/:id endpoint with product ID 1
	req, err := http.NewRequest("GET", "/products/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the HTTP request to the response recorder
	router.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Ensure the expected data is present in the response
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Expected 'data' field in response")
	assert.Equal(t, float64(1), data["id"].(float64))   // Example check for ID
	assert.Equal(t, "Product 1", data["name"].(string)) // Example check for Name
	// Add more assertions as needed for other fields
}

func TestCreateProducts(t *testing.T) {
	// Setup a mock Gin router
	router := gin.Default()
	router.POST("/products", func(c *gin.Context) {
		controller.CreateProducts(c)
	})

	// Mock JSON request body for product creation
	reqBody := `{"name":"New Product","brand":"Brand X","price":123.45,"description":"New Description","stock":5}`
	req, err := http.NewRequest("POST", "/products", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the HTTP request to the response recorder
	router.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse the response body into a map
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Ensure the response contains the 'data' field
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "Expected 'data' field in response")

	// Assert specific fields of the newly created product
	assert.Equal(t, "New Product", data["name"].(string))
	assert.Equal(t, "Brand X", data["brand"].(string))
	assert.Equal(t, 123.45, data["price"].(float64))
	assert.Equal(t, "New Description", data["description"].(string))
	assert.Equal(t, 5.0, data["stock"].(float64))

	// Optionally, assert other fields as needed

	// Ensure no errors occurred during the test
	assert.Empty(t, response["message"], "Unexpected error message in response")
}
