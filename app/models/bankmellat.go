package models

import (
	"gorm.io/gorm"
	"time"
)

type MellaForm struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	RefID         string    `json:"ref_id"`
	EncPan        string    `json:"enc_pan"`
	Enc           string    `json:"enc"`
	PhoneNumber   string    `json:"phone_number"`
	Body          string    `json:"body"`
	SaleOrderID   string    `json:"sale_order_id"`  // Based on Mellat's SaleOrderID
	SaleReference string    `json:"sale_reference"` // Sale Reference from Mellat
	ResCode       string    `json:"res_code"`       // Response Code from Mellat
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"` // Added for automatic Gorm timestamping
}

// Gorm methods (Optional, depending on Goravel configuration)
func (m *MellaForm) BeforeCreate(tx *gorm.DB) (err error) {
	// Add any pre-insert logic here, like setting CreatedAt
	m.CreatedAt = time.Now()
	return
}

func (m *MellaForm) BeforeUpdate(tx *gorm.DB) (err error) {
	// Add any pre-update logic here, like setting UpdatedAt
	m.UpdatedAt = time.Now()
	return
}
