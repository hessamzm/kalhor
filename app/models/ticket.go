package models

import (
	"github.com/goravel/framework/database/orm"
	"time"
)

type Ticket struct {
	orm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      uint   `json:"user_id"`
	Department  string `json:"department"`
	Priority    string `json:"priority"`
	Product     string `json:"product"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
