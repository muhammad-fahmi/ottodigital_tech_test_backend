package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Brand represents a brand entity
type Brand struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	LogoURL     string    `json:"logo_url" gorm:"size:500"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Vouchers    []Voucher `json:"vouchers,omitempty" gorm:"foreignKey:BrandID"`
}

// Voucher represents a voucher entity
type Voucher struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	BrandID     uuid.UUID `json:"brand_id" gorm:"type:char(36);not null"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CostInPoint int       `json:"cost_in_point" gorm:"not null"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     time.Time `json:"valid_to"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Brand       Brand     `json:"brand,omitempty" gorm:"foreignKey:BrandID"`
}

// Customer represents a customer entity
type Customer struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Email     string    `json:"email" gorm:"size:255;unique;not null"`
	Phone     string    `json:"phone" gorm:"size:20"`
	Points    int       `json:"points" gorm:"default:0"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Transaction represents a redemption transaction
type Transaction struct {
	ID          uuid.UUID         `json:"id" gorm:"type:char(36);primary_key"`
	CustomerID  uuid.UUID         `json:"customer_id" gorm:"type:char(36);not null"`
	TotalPoints int               `json:"total_points" gorm:"not null"`
	Status      string            `json:"status" gorm:"size:50;default:'pending'"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Customer    Customer          `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Items       []TransactionItem `json:"items,omitempty" gorm:"foreignKey:TransactionID"`
}

// TransactionItem represents individual voucher items in a transaction
type TransactionItem struct {
	ID            uuid.UUID   `json:"id" gorm:"type:char(36);primary_key"`
	TransactionID uuid.UUID   `json:"transaction_id" gorm:"type:char(36);not null"`
	VoucherID     uuid.UUID   `json:"voucher_id" gorm:"type:char(36);not null"`
	Quantity      int         `json:"quantity" gorm:"not null"`
	PointsPerUnit int         `json:"points_per_unit" gorm:"not null"`
	TotalPoints   int         `json:"total_points" gorm:"not null"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	Transaction   Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
	Voucher       Voucher     `json:"voucher,omitempty" gorm:"foreignKey:VoucherID"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (brand *Brand) BeforeCreate(tx *gorm.DB) error {
	if brand.ID == uuid.Nil {
		brand.ID = uuid.New()
	}
	return nil
}

func (voucher *Voucher) BeforeCreate(tx *gorm.DB) error {
	if voucher.ID == uuid.Nil {
		voucher.ID = uuid.New()
	}
	return nil
}

func (customer *Customer) BeforeCreate(tx *gorm.DB) error {
	if customer.ID == uuid.Nil {
		customer.ID = uuid.New()
	}
	return nil
}

func (transaction *Transaction) BeforeCreate(tx *gorm.DB) error {
	if transaction.ID == uuid.Nil {
		transaction.ID = uuid.New()
	}
	return nil
}

func (item *TransactionItem) BeforeCreate(tx *gorm.DB) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	return nil
}
