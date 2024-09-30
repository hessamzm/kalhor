package models

import "time"

type UserWallet struct {
	ID               uint64 `json:"id" gorm:"primary_key"`
	MelliNumber      string `json:"mellinumber" gorm:"unique"`
	BallanceGold     string `json:"ballancegold"`
	BallanceRial     string `json:"ballancerial"`
	Income           string `json:"income"`
	Outcome          string `json:"outcome"`
	Shabanumber      string `json:"shabanumber"`
	ShomareTrakonesh string `json:"shomaretrakonesh"`
	createat         time.Time
}
