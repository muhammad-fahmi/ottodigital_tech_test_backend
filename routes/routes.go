package routes

import (
	"my-backend-app/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Brand routes
		brands := v1.Group("/brand")
		{
			brands.POST("", handlers.CreateBrand)
			brands.GET("", handlers.GetBrands)
			brands.GET("/:id", handlers.GetBrand)
		}

		// Voucher routes
		vouchers := v1.Group("/voucher")
		{
			vouchers.POST("", handlers.CreateVoucher)
			vouchers.GET("", handlers.GetVoucher)
			vouchers.GET("/brand", handlers.GetVouchersByBrand)
			vouchers.GET("/all", handlers.GetVouchers)
		}

		// Customer routes
		customers := v1.Group("/customer")
		{
			customers.POST("", handlers.CreateCustomer)
			customers.GET("", handlers.GetCustomers)
			customers.GET("/:id", handlers.GetCustomer)
			customers.PUT("/:id/points", handlers.UpdateCustomerPoints)
		}

		// Transaction routes
		transactions := v1.Group("/transaction")
		{
			transactions.POST("/redemption", handlers.CreateRedemption)
			transactions.GET("/redemption", handlers.GetTransactionDetail)
			transactions.GET("/customer", handlers.GetCustomerTransactions)
		}
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Voucher System API is running",
		})
	})
}
