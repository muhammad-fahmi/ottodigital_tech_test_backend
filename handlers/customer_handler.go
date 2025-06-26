package handlers

import (
	"net/http"
	"strconv"

	"my-backend-app/database"
	"my-backend-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateCustomerRequest represents the request body for creating a customer
type CreateCustomerRequest struct {
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	Phone  string `json:"phone"`
	Points int    `json:"points"`
}

// CreateCustomer creates a new customer
func CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate name length
	if len(req.Name) < 2 || len(req.Name) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer name must be between 2 and 255 characters"})
		return
	}

	// Check if email already exists
	var existingCustomer models.Customer
	if err := database.GetDB().Where("email = ?", req.Email).First(&existingCustomer).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Set default points if not provided
	if req.Points < 0 {
		req.Points = 0
	}

	customer := models.Customer{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Points:   req.Points,
		IsActive: true,
	}

	if err := database.GetDB().Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer created successfully",
		"data":    customer,
	})
}

// GetCustomer gets a single customer by ID
func GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customerID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var customer models.Customer
	if err := database.GetDB().First(&customer, "id = ?", customerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// GetCustomers gets all customers with pagination
func GetCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var customers []models.Customer
	var total int64

	database.GetDB().Model(&models.Customer{}).Count(&total)
	if err := database.GetDB().Offset(offset).Limit(limit).Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": customers,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateCustomerPoints updates customer points
func UpdateCustomerPoints(c *gin.Context) {
	id := c.Param("id")
	customerID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var req struct {
		Points int `json:"points" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Points < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Points cannot be negative"})
		return
	}

	var customer models.Customer
	if err := database.GetDB().First(&customer, "id = ?", customerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if err := database.GetDB().Model(&customer).Update("points", req.Points).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer points"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer points updated successfully",
		"data":    customer,
	})
}
