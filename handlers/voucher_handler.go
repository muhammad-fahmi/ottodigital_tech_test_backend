package handlers

import (
	"net/http"
	"strconv"
	"time"

	"my-backend-app/database"
	"my-backend-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateVoucherRequest represents the request body for creating a voucher
type CreateVoucherRequest struct {
	BrandID     string    `json:"brand_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	CostInPoint int       `json:"cost_in_point" binding:"required,min=1"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     time.Time `json:"valid_to"`
}

// CreateVoucher creates a new voucher
func CreateVoucher(c *gin.Context) {
	var req CreateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse brand ID
	brandID, err := uuid.Parse(req.BrandID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
		return
	}

	// Check if brand exists
	var brand models.Brand
	if err := database.GetDB().First(&brand, "id = ?", brandID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand not found"})
		return
	}

	// Validate cost in point
	if req.CostInPoint <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cost in point must be greater than 0"})
		return
	}

	// Validate date range
	if !req.ValidFrom.IsZero() && !req.ValidTo.IsZero() && req.ValidFrom.After(req.ValidTo) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid from date must be before valid to date"})
		return
	}

	voucher := models.Voucher{
		BrandID:     brandID,
		Name:        req.Name,
		Description: req.Description,
		CostInPoint: req.CostInPoint,
		ValidFrom:   req.ValidFrom,
		ValidTo:     req.ValidTo,
		IsActive:    true,
	}

	if err := database.GetDB().Create(&voucher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create voucher"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Voucher created successfully",
		"data":    voucher,
	})
}

// GetVoucher gets a single voucher by ID
func GetVoucher(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Voucher ID is required"})
		return
	}

	voucherID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voucher ID"})
		return
	}

	var voucher models.Voucher
	if err := database.GetDB().Preload("Brand").First(&voucher, "id = ?", voucherID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voucher not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": voucher})
}

// GetVouchersByBrand gets all vouchers for a specific brand
func GetVouchersByBrand(c *gin.Context) {
	brandID := c.Query("id")
	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID is required"})
		return
	}

	parsedBrandID, err := uuid.Parse(brandID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var vouchers []models.Voucher
	var total int64

	database.GetDB().Model(&models.Voucher{}).Where("brand_id = ?", parsedBrandID).Count(&total)
	if err := database.GetDB().Where("brand_id = ?", parsedBrandID).Offset(offset).Limit(limit).Find(&vouchers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vouchers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": vouchers,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetVouchers gets all vouchers with pagination
func GetVouchers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var vouchers []models.Voucher
	var total int64

	database.GetDB().Model(&models.Voucher{}).Count(&total)
	if err := database.GetDB().Preload("Brand").Offset(offset).Limit(limit).Find(&vouchers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vouchers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": vouchers,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
