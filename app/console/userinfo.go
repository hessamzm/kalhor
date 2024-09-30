package console

import (
	"fmt"
	"github.com/goravel/framework/facades"
	"kalhor/app/models"
	"kalhor/services"
	"kalhor/utils"
)

var n models.Notification

func UserInfo(u models.User) map[string]any {
	var byprice float64 = 1

	s, e := services.NewWalletService()

	if e != nil {
		return map[string]any{"error": e.Error()}
	}

	totalfeeby, e := s.GetTotalFeeBy(u.MelliNumber)

	totalgold, e := s.GetBalanceDifference(u.MelliNumber)

	totalrial, e := s.GetBalanceDifferenceRial(u.MelliNumber)

	feenow := totalgold * byprice

	howprofitrial := feenow - totalfeeby

	totlamony := (totalgold * byprice) + totalrial

	howprofitinhondred := howprofitrial * 100 / totalfeeby

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
