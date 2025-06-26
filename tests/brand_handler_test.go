package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"my-backend-app/database"
	"my-backend-app/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BrandHandlerTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *BrandHandlerTestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set test mode environment variable
	os.Setenv("TEST_MODE", "true")

	// Load test environment variables
	godotenv.Load("config.env")

	// Initialize test database
	database.InitDB()

	// Setup router
	suite.router = gin.New()
	suite.router.POST("/brand", handlers.CreateBrand)
	suite.router.GET("/brand", handlers.GetBrands)
	suite.router.GET("/brand/:id", handlers.GetBrand)
}

func (suite *BrandHandlerTestSuite) TearDownSuite() {
	// Clean up test database if needed
	if database.DB != nil {
		sqlDB, err := database.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

func (suite *BrandHandlerTestSuite) TestCreateBrand_Success() {
	// Test data
	brandData := handlers.CreateBrandRequest{
		Name:        "Test Brand",
		Description: "Test Description",
		LogoURL:     "https://example.com/logo.png",
		IsActive:    true,
	}

	jsonData, _ := json.Marshal(brandData)

	// Create request
	req, _ := http.NewRequest("POST", "/brand", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "Brand created successfully", response["message"])
	assert.NotNil(suite.T(), response["data"])
}

func (suite *BrandHandlerTestSuite) TestCreateBrand_InvalidName() {
	// Test data with invalid name
	brandData := handlers.CreateBrandRequest{
		Name:        "A", // Too short
		Description: "Test Description",
		IsActive:    true,
	}

	jsonData, _ := json.Marshal(brandData)

	// Create request
	req, _ := http.NewRequest("POST", "/brand", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Contains(suite.T(), response["error"], "Brand name must be between 2 and 255 characters")
}

func (suite *BrandHandlerTestSuite) TestCreateBrand_MissingName() {
	// Test data without required name
	brandData := handlers.CreateBrandRequest{
		Description: "Test Description",
		IsActive:    true,
	}

	jsonData, _ := json.Marshal(brandData)

	// Create request
	req, _ := http.NewRequest("POST", "/brand", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *BrandHandlerTestSuite) TestGetBrands_Success() {
	// Create request
	req, _ := http.NewRequest("GET", "/brand", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.NotNil(suite.T(), response["data"])
	assert.NotNil(suite.T(), response["pagination"])
}

func TestBrandHandlerSuite(t *testing.T) {
	suite.Run(t, new(BrandHandlerTestSuite))
}
