package controllers

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"kalhor/app/models"
	"kalhor/utils"
	"time"
)

var price models.Price

var statusPrice struct {
	Status string `json:"status"`
}

var inputPrice struct {
	SellPrice float64 `json:"sellprice"` // تگ JSON به درستی تعریف شده
	ByPrice   float64 `json:"byprice"`   // حذف تگ اضافی JSON
	Base_18   float64 `json:"base_18"`
	Base_24   float64 `json:"base_24"`
}

var ojratPrice struct {
	Ojrat  float64 `json:"ojrat"`
	Maliat float64 `json:"maliat"`
	Sood   float64 `json:"sood"`
}

type PriceController struct {
	//Dependent services
}

func NewPriceController() *PriceController {
	return &PriceController{
		//Inject services
	}
}

func (p *PriceController) Get(ctx http.Context) http.Response {

	err := facades.Orm().Query().Order("updated_at desc").First(&price)
	if err != nil {
		// اگر خطا رخ دهد
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to retrieve price",
		})
	}

	// چاپ داده‌ها برای اشکال‌زدایی
	fmt.Printf("Retrieved price data: %+v\n", price)

	// اگر داده‌ای پیدا شد
	return ctx.Response().Json(http.StatusOK, http.Json{
		"قیمت خرید":         price.ByPrice,
		"قیمت فروش":         price.SellPrice,
		"پایه قیمت 18 عیار": price.Base_18,
		"پایه قیمت 24 عیار": price.Base_24,
	})
}

func (p *PriceController) Put(ctx http.Context) http.Response {
	// Bind کردن داده‌های JSON به مدل
	if err := ctx.Request().Bind(&inputPrice); err != nil {
		// اگر خطایی در Bind داده‌ها وجود داشته باشد

		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": err.Error(),
		})
	}

	// تنظیم مقادیر Status و UpdatedAt
	price.UpdatedAt = time.Now()

	if err := facades.Orm().Query().Order("updated_at desc").First(&price); err != nil {
		// اگر رکوردی پیدا نشود
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"error": "No recent price found",
		})
	}

	// به‌روزرسانی یا ایجاد رکورد
	price.SellPrice = inputPrice.SellPrice
	price.ByPrice = inputPrice.ByPrice
	price.Base_18 = inputPrice.Base_18
	price.Base_24 = inputPrice.Base_24
	price.UpdatedAt = time.Now()

	// به‌روزرسانی رکورد
	if err := facades.Orm().Query().Model(&price).Save(&price); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to update price",
		})
	}

	// پاسخ موفقیت
	return ctx.Response().Json(http.StatusOK, http.Json{
		"msg":    "Status updated successfully",
		"status": inputPrice, // مقدار جدید Status را برمی‌گردانیم
	})

}

func (p *PriceController) ChStatus(ctx http.Context) http.Response {

	// Bind کردن داده‌های ورودی به statusPrice (فقط فیلد status را بایند می‌کنیم)
	if err := ctx.Request().Bind(&statusPrice); err != nil {
		// اگر خطایی در Bind داده‌ها وجود داشته باشد
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": err.Error(),
		})
	}

	// دریافت جدیدترین رکورد بر اساس UpdatedAt
	if err := facades.Orm().Query().Order("updated_at desc").First(&price); err != nil {
		// اگر رکوردی پیدا نشود
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"error": "No recent price found",
		})
	}

	// به‌روزرسانی فیلد Status با مقدار جدید از statusPrice
	if _, err := facades.Orm().Query().Model(&price).Update("status", statusPrice.Status); err != nil {
		// اگر خطایی در به‌روزرسانی وجود داشته باشد
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to update status",
		})
	}

	// پاسخ موفقیت
	return ctx.Response().Json(http.StatusOK, http.Json{
		"msg":    "Status updated successfully",
		"status": statusPrice.Status, // مقدار جدید Status را برمی‌گردانیم
	})
}

func (p *PriceController) Putograt(ctx http.Context) http.Response {

	// Bind کردن داده‌های JSON به مدل
	if err := ctx.Request().Bind(&ojratPrice); err != nil {
		// اگر خطایی در Bind داده‌ها وجود داشته باشد

		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": err.Error(),
		})
	}

	// تنظیم مقادیر Status و UpdatedAt
	price.UpdatedAt = time.Now()

	if err := facades.Orm().Query().Order("updated_at desc").First(&price); err != nil {
		if utils.KlDebug {
			fmt.Printf("row is %+v\n", err)
		}
		// اگر رکوردی پیدا نشود
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"error": "No recent price found",
		})
	}

	// به‌روزرسانی یا ایجاد رکورد
	price.Ojrat = ojratPrice.Ojrat
	price.Sood = ojratPrice.Sood
	price.Maliat = ojratPrice.Maliat

	// به‌روزرسانی رکورد
	if err := facades.Orm().Query().Model(&price).Save(&price); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to update price",
		})
	}

	// پاسخ موفقیت
	return ctx.Response().Json(http.StatusOK, http.Json{
		"msg":    "Status updated successfully",
		"status": ojratPrice, // مقدار جدید Status را برمی‌گردانیم
	})

}

///atomation

func (p *PriceController) PutAuotmation(ctx http.Context) http.Response {

	data, err := utils.GetDatamethodget("http://0.0.0.0:3000/v1/price", "error")
	if err != nil {
		println(err)
	}

	if utils.KlDebug {
		fmt.Printf("data: %+v\n", data)
	}

	ByPrice, ok := data["قیمت خرید"].(float64)
	if !ok {
		println("Error converting ByPrice")

	}

	SellPrice, ok := data["قیمت فروش"].(float64)
	if !ok {
		println("Error converting SellPrice")
	}

	Base_18, ok := data["پایه قیمت 18 عیار"].(float64)
	if !ok {
		println("Error converting Base_18")
	}

	Base_24, ok := data["پایه قیمت 24 عیار"].(float64)
	if !ok {
		println("Error converting Base_24")
	}

	if err := facades.Orm().Query().Order("updated_at desc").First(&price); err != nil {
		// اگر رکوردی پیدا نشود
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"error": "No recent price found",
		})
	}

	if utils.KlDebug {
		// لاگ کردن داده‌های قیمت قبل از ذخیره‌سازی برای بررسی مقادیر
		// use flag --kldebug
		fmt.Printf("SellPrice: %f, ByPrice: %f, Base_18: %f, Base_24: %f\n", SellPrice, ByPrice, Base_18, Base_24)
	}
	// لاگ کردن داده‌های قیمت قبل از ذخیره‌سازی برای بررسی مقادیر

	// به‌روزرسانی یا ایجاد رکورد
	price.SellPrice = SellPrice
	price.ByPrice = ByPrice
	price.Base_18 = Base_18
	price.Base_24 = Base_24
	price.UpdatedAt = time.Now()

	// به‌روزرسانی رکورد
	if err := facades.Orm().Query().Model(&price).Save(&price); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to update price",
		})
	}

	// پاسخ موفقیت
	return ctx.Response().Json(http.StatusOK, http.Json{
		"msg":       "Status updated successfully",
		"sellprice": price.SellPrice,
		"Base_18":   price.Base_18,
		"Base_24":   price.Base_24,
		"ByPrice":   price.ByPrice,
	})
}
