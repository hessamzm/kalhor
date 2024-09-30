package models

import (
	"time"
)

type Price struct {
	ID        uint      `json:"id"`        // استفاده از uint برای شناسه
	SellPrice float64   `json:"sellprice"` // تگ JSON به درستی تعریف شده
	ByPrice   float64   `json:"byprice"`   // حذف تگ اضافی JSON
	Status    string    `json:"status"`
	Base_18   float64   `json:"base_18"`
	Base_24   float64   `json:"base_24"`
	Ojrat     float64   `json:"ojrat"`
	Maliat    float64   `json:"maliat"`
	Sood      float64   `json:"sood"`
	UpdatedAt time.Time `json:"updated_at"`
}
