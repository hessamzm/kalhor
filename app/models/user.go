package models

import (
	"github.com/goravel/framework/database/orm"
	"time"
)

type User struct {
	orm.Model
	Name          string
	MelliNumber   string    `json:"mellinumber" gorm:"unique"`
	Phone         string    `json:"username"`
	Email         string    `json:"email"`
	KartInfo      string    `json:"kartinfo"`
	TarikhTavalod time.Time `json:"tarikh_tavalod"`
	StepLvl       int64
	UserNum       string
	OtpCode       bool
	Freez         bool
	UpdatedAt     time.Time
	orm.SoftDeletes
}
