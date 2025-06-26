package handlers

import (
	"net/http"
	"strconv"

	"my-backend-app/database"
	"my-backend-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateBrandRequest represents the request body for creating a brand
type CreateBrandRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	IsActive    bool   `json:"is_active"`
}

// CreateBrand creates a new brand
func CreateBrand(c *gin.Context) {
	var req CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate name length
	if len(req.Name) < 2 || len(req.Name) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand name must be between 2 and 255 characters"})
		return
	}

	brand := models.Brand{
		Name:        req.Name,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		IsActive:    req.IsActive,
	}

	if err := database.GetDB().Create(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create brand"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Brand created successfully",
		"data":    brand,
	})
}

// GetBrand gets a single brand by ID
func GetBrand(c *gin.Context) {
	id := c.Param("id")
	brandID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
		return
	}

	var brand models.Brand
	if err := database.GetDB().First(&brand, "id = ?", brandID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": brand})
}

// GetBrands gets all brands with pagination
func GetBrands(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var brands []models.Brand
	var total int64

	database.GetDB().Model(&models.Brand{}).Count(&total)
	if err := database.GetDB().Offset(offset).Limit(limit).Find(&brands).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch brands"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": brands,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
