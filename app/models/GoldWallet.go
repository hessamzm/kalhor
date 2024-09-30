package models

import (
	"time"
)

type WalletGold struct {
	UserID        uint64    `ch:"user_id"`      // برای ذخیره شناسه کاربر
	MelliNumber   string    `ch:"melli_number"` // برای ذخیره شناسه ملی کاربر
	BalanceIn     float64   `ch:"balance_in"`   // موجودی ریالی کاربر
	BalanceOut    float64   `ch:"balance_out"`  // موجودی ریالی کاربر
	FeebalanceIn  float64   `ch:"feebalance_in"`
	FeebalanceOut float64   `ch:"feebalance_out"`
	FreezBlIn     float64   `ch:"freez_bl_in"`
	FreezBlOut    float64   `ch:"freez_bl_out"`
	BanBlIn       float64   `ch:"ban_bl_in"`
	BanBlOut      float64   `ch:"ban_bl_out"`
	EventTime     time.Time `ch:"event_time"` // ثبت زمان تراکنش
}
