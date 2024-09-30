package models

import (
	"github.com/goravel/framework/database/orm"
	"time"
)

type Notification struct {
	orm.Model
	ToWho           string    `json:"to" gorm:"index"` // ایجاد ایندکس برای سریع‌تر کردن جستجو
	Subject         string    `json:"subject" gorm:"unique"`
	Messages        string    `json:"messages"` // اصلاح نام فیلد
	IsSee           bool      `json:"is_see" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at"` // اصلاح نام فیلد
	orm.SoftDeletes           // برای پشتیبانی از حذف نرم
}
