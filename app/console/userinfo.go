package console

import (
	"fmt"
	"github.com/goravel/framework/facades"
	"kalhor/app/models"
	"kalhor/services"
	"kalhor/utils"
	"math"
)

var n models.Notification
var price models.Price

func UserInfo(u models.User) map[string]any {
	err := facades.Orm().Query().Order("updated_at desc").First(&price)
	if err != nil {
		// اگر خطا رخ دهد
		return map[string]any{
			"error": "Failed to retrieve price",
		}
	}
	byprice := price.ByPrice
	s, e := services.NewWalletService()

	// جمع کل هزینه‌ها
	totalfeeby, e := s.GetTotalFeeBy(u.MelliNumber)

	if math.IsNaN(totalfeeby) {
		totalfeeby = 0
	}

	// جمع کل طلا
	totalgold, e := s.GetBalanceDifference(u.MelliNumber)

	if math.IsNaN(totalgold) {
		totalgold = 0
	}

	// جمع کل ریال
	totalrial, e := s.GetBalanceDifferenceRial(u.MelliNumber)

	if math.IsNaN(totalrial) {
		totalrial = 0
	}

	// محاسبات سود
	feenow := totalgold * byprice
	howprofitrial := feenow - totalfeeby
	if math.IsNaN(howprofitrial) {
		howprofitrial = 0
	}

	// جمع کل پول
	totlamony := (totalgold * byprice) + totalrial
	if math.IsNaN(totlamony) {
		totlamony = 0
	}

	// محاسبه درصد سود
	howprofitinhondred := 0.0
	if totalfeeby != 0 {
		howprofitinhondred = howprofitrial * 100 / totalfeeby
	}

	if e != nil {
		return map[string]any{"error": e.Error()}
	}

	if utils.KlDebug {
		fmt.Println("WalletGold:", totalgold)
	}
	return map[string]any{
		"name":      u.Name,
		"gold":      totalgold,
		"rial":      totalrial,
		"userlevel": u.StepLvl,
		"usernum":   u.UserNum,
		"totalmony": totlamony,
		"sodrial":   howprofitrial,
		"soddarsad": howprofitinhondred,
	}
}

func UserNotif(to string) map[string]any {

	var notifications []models.Notification

	// جستجو در دیتابیس بر اساس فیلد "to_who"
	facades.Orm().Query().Where("to_who = ?", to).Get(&notifications)

	// اگر هیچ پیامی یافت نشد
	if len(notifications) == 0 {
		return map[string]any{
			"error": "no messages found",
		}
	}

	// برگرداندن لیست پیام‌ها
	return map[string]any{
		"notifications": notifications,
	}
}

func UserWallet(u models.User) map[string]any {
	price = models.Price{}

	s, e := services.NewWalletService()
	profit, e := s.GetBalanceDifference(u.MelliNumber)
	totalrial, e := s.GetBalanceDifferenceRial(u.MelliNumber)

	err := facades.Orm().Query().Order("updated_at desc").First(&price)
	if err != nil {
		// اگر خطا رخ دهد
		return map[string]any{
			"error": "Failed to retrieve price",
		}
	}

	if utils.KlDebug {
		// چاپ داده‌ها برای اشکال‌زدایی
		fmt.Printf("Retrieved price data: %+v\n", price)

	}

	if e != nil {
		return map[string]any{"error": e.Error()}
	}

	goldtorial := profit * price.ByPrice

	return map[string]any{
		"user-number": u.UserNum,
		"profit-gold": profit,
		"profit-rial": totalrial,
		"rial-gold":   goldtorial,
	}
}
