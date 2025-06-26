package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"my-backend-app/database"
	"my-backend-app/handlers"
	"my-backend-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VoucherHandlerTestSuite struct {
	suite.Suite
	router  *gin.Engine
	brandID uuid.UUID
}

func (suite *VoucherHandlerTestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set test mode environment variable
	os.Setenv("TEST_MODE", "true")

	// Load test environment variables
	godotenv.Load("config.env")

	// Initialize test database
	database.InitDB()

	// Create a test brand
	brand := models.Brand{
		ID:          uuid.New(),
		Name:        "Test Brand",
		Description: "Test Brand Description",
		LogoURL:     "https://example.com/logo.png",
		IsActive:    true,
	}
	database.GetDB().Create(&brand)
	suite.brandID = brand.ID

	// Setup router
	suite.router = gin.New()
	suite.router.POST("/voucher", handlers.CreateVoucher)
	suite.router.GET("/voucher", handlers.GetVoucher)
	suite.router.GET("/voucher/brand", handlers.GetVouchersByBrand)
	suite.router.GET("/voucher/all", handlers.GetVouchers)
}

func (suite *VoucherHandlerTestSuite) TearDownSuite() {
	// Clean up test database if needed
	if database.DB != nil {
		sqlDB, err := database.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

func (suite *VoucherHandlerTestSuite) TestCreateVoucher_Success() {
	// Test data
	voucherData := handlers.CreateVoucherRequest{
		BrandID:     suite.brandID.String(),
		Name:        "Test Voucher",
		Description: "Test Description",
		CostInPoint: 1000,
		ValidFrom:   time.Now(),
		ValidTo:     time.Now().AddDate(0, 1, 0), // 1 month from now
	}

	jsonData, _ := json.Marshal(voucherData)

	// Create request
	req, _ := http.NewRequest("POST", "/voucher", bytes.NewBuffer(jsonData))
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

	assert.Equal(suite.T(), "Voucher created successfully", response["message"])
	assert.NotNil(suite.T(), response["data"])
}

func (suite *VoucherHandlerTestSuite) TestCreateVoucher_InvalidBrandID() {
	// Test data with invalid brand ID
	voucherData := handlers.CreateVoucherRequest{
		BrandID:     "invalid-uuid",
		Name:        "Test Voucher",
		CostInPoint: 1000,
	}

	jsonData, _ := json.Marshal(voucherData)

	// Create request
	req, _ := http.NewRequest("POST", "/voucher", bytes.NewBuffer(jsonData))
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

	assert.Contains(suite.T(), response["error"], "Invalid brand ID")
}

func (suite *VoucherHandlerTestSuite) TestCreateVoucher_InvalidCostInPoint() {
	// Test data with invalid cost in point
	voucherData := handlers.CreateVoucherRequest{
		BrandID:     suite.brandID.String(),
		Name:        "Test Voucher",
		CostInPoint: 0, // Invalid: must be greater than 0
	}

	jsonData, _ := json.Marshal(voucherData)

	// Create request
	req, _ := http.NewRequest("POST", "/voucher", bytes.NewBuffer(jsonData))
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

	// Check for validation error (either from binding or custom validation)
	assert.True(suite.T(),
		response["error"] == "Key: 'CreateVoucherRequest.CostInPoint' Error:Field validation for 'CostInPoint' failed on the 'required' tag" ||
			response["error"] == "Cost in point must be greater than 0",
		"Expected validation error for cost in point")
}

func (suite *VoucherHandlerTestSuite) TestCreateVoucher_InvalidDateRange() {
	// Test data with invalid date range
	voucherData := handlers.CreateVoucherRequest{
		BrandID:     suite.brandID.String(),
		Name:        "Test Voucher",
		CostInPoint: 1000,
		ValidFrom:   time.Now().AddDate(0, 1, 0), // 1 month from now
		ValidTo:     time.Now(),                  // Now (invalid: from > to)
	}

	jsonData, _ := json.Marshal(voucherData)

	// Create request
	req, _ := http.NewRequest("POST", "/voucher", bytes.NewBuffer(jsonData))
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

	assert.Contains(suite.T(), response["error"], "Valid from date must be before valid to date")
}

func (suite *VoucherHandlerTestSuite) TestGetVoucher_MissingID() {
	// Create request without ID
	req, _ := http.NewRequest("GET", "/voucher", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Contains(suite.T(), response["error"], "Voucher ID is required")
}

func (suite *VoucherHandlerTestSuite) TestGetVouchersByBrand_MissingBrandID() {
	// Create request without brand ID
	req, _ := http.NewRequest("GET", "/voucher/brand", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Contains(suite.T(), response["error"], "Brand ID is required")
}

func TestVoucherHandlerSuite(t *testing.T) {
	suite.Run(t, new(VoucherHandlerTestSuite))
}
