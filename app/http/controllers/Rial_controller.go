package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"kalhor/app/models"
	"kalhor/services"
	"strconv"
)

type RialWalletController struct {
	walletrial *services.WalletServiceRial
}

func (g *RialWalletController) SharjHesab(ctx http.Context) http.Response {
	// اعتبارسنجی JWT
	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}

	// دریافت مبلغ از پارامترهای درخواست
	amountStr := ctx.Request().Input("amount")
	amount, e := strconv.ParseInt(amountStr, 10, 64)
	if e != nil || amount <= 0 {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "مبلغ نامعتبر است"})
	}

	// دریافت اطلاعات کاربر
	var user models.User
	facades.Auth(ctx).User(&user)

	//
	//client := services.MellatClient()
	//
	//// ایجاد درخواست پرداخت
	//
	////paymentService, error := client.PaymentRequest(1, amount, time.Now(), time.Now(), "test", "app.kalhorgold.ir/form")
	//
	//if error != nil || amount <= 0 {
	//	return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "مبلغ نامعتبر است"})
	//}

	//paymentUrl, err := paymentService.CreatePayment(user.ID, amount)
	//if err != nil {
	//	return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "خطا در ایجاد درخواست پرداخت"})
	//}
	//
	//// ذخیره اطلاعات تراکنش در دیتابیس
	//err = g.walletrial.SaveTransaction(user.ID, amount, "pending")
	//if err != nil {
	//	return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "خطا در ذخیره اطلاعات تراکنش"})
	//}

	// ارسال لینک پرداخت به کاربر
	return ctx.Response().Json(http.StatusOK, map[string]string{"payment_url": "paymentUrl"})
}

func (g *GoldController) DeSharjHesab(ctx http.Context) http.Response {
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
