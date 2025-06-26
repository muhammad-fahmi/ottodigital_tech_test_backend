package handlers

import (
	"net/http"
	"time"

	"my-backend-app/database"
	"my-backend-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RedemptionItem represents a voucher item in redemption request
type RedemptionItem struct {
	VoucherID string `json:"voucher_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// RedemptionRequest represents the request body for redemption
type RedemptionRequest struct {
	CustomerID string           `json:"customer_id" binding:"required"`
	Items      []RedemptionItem `json:"items" binding:"required,min=1"`
}

// CreateRedemption creates a new redemption transaction
func CreateRedemption(c *gin.Context) {
	var req RedemptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse customer ID
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	// Check if customer exists
	var customer models.Customer
	if err := database.GetDB().First(&customer, "id = ?", customerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer not found"})
		return
	}

	// Calculate total points and validate vouchers
	totalPoints := 0
	var transactionItems []models.TransactionItem

	for _, item := range req.Items {
		// Parse voucher ID
		voucherID, err := uuid.Parse(item.VoucherID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voucher ID"})
			return
		}

		// Get voucher details
		var voucher models.Voucher
		if err := database.GetDB().First(&voucher, "id = ? AND is_active = ?", voucherID, true).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Voucher not found or inactive"})
			return
		}

		// Validate voucher validity period
		now := time.Now()
		if !voucher.ValidFrom.IsZero() && now.Before(voucher.ValidFrom) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Voucher is not yet valid"})
			return
		}
		if !voucher.ValidTo.IsZero() && now.After(voucher.ValidTo) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Voucher has expired"})
			return
		}

		// Calculate points for this item
		itemTotalPoints := voucher.CostInPoint * item.Quantity
		totalPoints += itemTotalPoints

		// Create transaction item
		transactionItem := models.TransactionItem{
			VoucherID:     voucherID,
			Quantity:      item.Quantity,
			PointsPerUnit: voucher.CostInPoint,
			TotalPoints:   itemTotalPoints,
		}
		transactionItems = append(transactionItems, transactionItem)
	}

	// Check if customer has enough points
	if customer.Points < totalPoints {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient points"})
		return
	}

	// Start transaction
	tx := database.GetDB().Begin()

	// Create transaction record
	transaction := models.Transaction{
		CustomerID:  customerID,
		TotalPoints: totalPoints,
		Status:      "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Create transaction items
	for i := range transactionItems {
		transactionItems[i].TransactionID = transaction.ID
		if err := tx.Create(&transactionItems[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction items"})
			return
		}
	}

	// Deduct points from customer
	if err := tx.Model(&customer).Update("points", customer.Points-totalPoints).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer points"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Load transaction with items and customer details
	var result models.Transaction
	database.GetDB().Preload("Items.Voucher.Brand").Preload("Customer").First(&result, transaction.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Redemption successful",
		"data":    result,
	})
}

// GetTransactionDetail gets transaction details by transaction ID
func GetTransactionDetail(c *gin.Context) {
	transactionID := c.Query("transactionId")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	parsedTransactionID, err := uuid.Parse(transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var transaction models.Transaction
	if err := database.GetDB().Preload("Items.Voucher.Brand").Preload("Customer").First(&transaction, "id = ?", parsedTransactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

// GetCustomerTransactions gets all transactions for a customer
func GetCustomerTransactions(c *gin.Context) {
	customerID := c.Query("customerId")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var transactions []models.Transaction
	if err := database.GetDB().Preload("Items.Voucher.Brand").Where("customer_id = ?", parsedCustomerID).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
