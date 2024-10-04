package models

import (
	"time"
)

type WalletRial struct {
	MelliNumber string    `ch:"melli_number"`
	BalanceIn   float64   `ch:"balance_in"`  // موجودی ریالی کاربر
	BalanceOut  float64   `ch:"balance_out"` // موجودی ریالی کاربر
	FreezBlIn   float64   `ch:"freez_bl_in"`
	FreezBlOut  float64   `ch:"freez_bl_out"`
	BanBlIn     float64   `ch:"ban_bl_in"`
	BanBlOut    float64   `ch:"ban_bl_out"`
	EventTime   time.Time `ch:"event_time"` // ثبت زمان تراکنش
	TrakoneshId string    `ch:"trakonesh_id"`
}
