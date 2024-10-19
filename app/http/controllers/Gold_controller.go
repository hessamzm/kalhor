package controllers

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"kalhor/app/models"
	"kalhor/services"
	"kalhor/utils"
	"strconv"
)

type GoldController struct{}

var InsertWalletGold models.WalletGold
var OutWalletRial models.WalletRial
var formgold struct {
	Q string `json:"qaunty"`
	R string `json:"rial"`
}

func (g *GoldController) ByGold(ctx http.Context) http.Response {

	s, e := services.NewWalletService()
	sr, e := services.NewWalletServiceRial()

	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}

	error := facades.Orm().Query().Order("updated_at desc").First(&price)
	if error != nil {
		// اگر خطا رخ دهد
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to retrieve price",
		})
	}

	byprice := price.Base_18 + price.Base_18*1/100

	if utils.KlDebug {
		// چاپ داده‌ها برای اشکال‌زدایی
		fmt.Printf("Retrieved price data: %+v\n", price)

	}

	// دریافت اطلاعات کاربر
	facades.Auth(ctx).User(&user)

	if err := ctx.Request().Bind(&formgold); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": err.Error(),
		})
	}

	totalrial, e := s.GetBalanceDifferenceRial(user.MelliNumber)

	rial, e := strconv.ParseFloat(formgold.R, 64)
	tetad, e := strconv.ParseFloat(formgold.Q, 64)
	if rial > totalrial || tetad*byprice > totalrial {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": "your requst more than your ceredits"})
	}

	InsertWalletGold.MelliNumber = user.MelliNumber
	InsertWalletGold.BalanceIn = tetad
	InsertWalletGold.FeebalanceIn = tetad * byprice

	OutWalletRial.MelliNumber = user.MelliNumber
	OutWalletRial.BalanceOut = rial
	OutWalletRial.TrakoneshId = "change"

	//
	//fmt.Println(InsertWalletGold.BalanceIn)
	e = s.InsertWalletGold(&InsertWalletGold)
	e = sr.OutWalletGoldRial(&OutWalletRial)

	if e != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{"error": e.Error()})
	}

	return ctx.Response().Json(http.StatusOK, map[string]any{
		"wallet_gold": InsertWalletGold,
	})
}

func (g *GoldController) SellGold(ctx http.Context) http.Response {
	s, e := services.NewWalletService()
	sr, e := services.NewWalletServiceRial()

	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}
	// دریافت اطلاعات کاربر
	facades.Auth(ctx).User(&user)

	if err := ctx.Request().Bind(&formgold); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": err.Error(),
		})
	}

	error := facades.Orm().Query().Order("updated_at desc").First(&price)
	if error != nil {
		// اگر خطا رخ دهد
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to retrieve price",
		})
	}

	Sellprice := price.Base_18 + price.Base_18*1/100

	totalgold, e := s.GetBalanceDifference(user.MelliNumber)

	q, e := strconv.ParseFloat(formgold.Q, 64)
	r, e := strconv.ParseFloat(formgold.R, 64)

	if totalgold < q {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": "your gold sell want more than your ceredits"})
	}

	if r != q*Sellprice {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": "your gold sell to rial take diffrent from server price"})

	}
	if e != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{"error": e.Error()})
	}

	InsertWalletGold.MelliNumber = user.MelliNumber
	InsertWalletGold.BalanceOut = q
	InsertWalletGold.FeebalanceOut = q * Sellprice

	OutWalletRial.MelliNumber = user.MelliNumber
	OutWalletRial.BalanceIn = r
	OutWalletRial.TrakoneshId = "change"

	e = s.TakeOutWalletGold(&InsertWalletGold)
	e = sr.InsertWalletGoldRial(&OutWalletRial)

	secend, e := s.GetBalanceDifference(user.MelliNumber)
	return ctx.Response().Json(http.StatusOK, map[string]any{
		"wallet_gold":           formgold,
		"total_gold_after_sell": secend,
	})
}
