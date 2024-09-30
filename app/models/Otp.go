package models

import (
	"github.com/goravel/framework/database/orm"
	"time"
)

type Otp struct {
	orm.Model
	Phone     string `json:"phone"`
	OtpCode   string `json:"otp_code"`
	Step      int
	Status    bool
	UpdatedAt time.Time
	orm.SoftDeletes
}
