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

type RialWalletController struct {
	walletrial *services.WalletServiceRial
}

const (
	merchantID      = "86fc011d-6b71-4e39-9c9b-bf8547689e3b"
	callbackURL     = "http://localhost:3000/verifyrial"
	apiURL          = "https://www.zarinpal.com/pg/services/WebGate/wsdl"
	zarinpalGateURL = "https://www.zarinpal.com/pg/StartPay/"
)

func (g *RialWalletController) AskSharjHesab(ctx http.Context) http.Response {

	input := models.Input{}
	user := models.User{}
	rw := &models.WalletRial{}

	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// اعتبارسنجی JWT
	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}

	// دریافت اطلاعات کاربر
	facades.Auth(ctx).User(&user)

	clinet := services.NewPaymentGatewayImplementationServicePortType("", false, nil)

	amountInt, error := strconv.Atoi(input.Amount)
	if error != nil {
		// هندل کردن خطا
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "Amount Error"})
	}

	// تبدیل amountInt از int به int32
	amount := int32(amountInt)

	if utils.KlDebug {

		fmt.Println("input.Amount :", input.Amount)
		fmt.Println("amountInt:", amountInt)
		fmt.Println("amount:", amount)

	}

	// Create a new payment request to Zarinpal
	resp, error := clinet.PaymentRequest(&services.PaymentRequest{
		MerchantID:  merchantID,
		Amount:      amount,
		Description: "شارژ حساب ریالی",
		Email:       "customer@domain.ir",
		Mobile:      user.Phone,
		//CardPan:     user.bank,
		CallbackURL: callbackURL,
	})

	// Check if response is error free
	if error != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"مشکلی در ارتباط رخ داد": error.Error(),
		})
	}

	if resp.Status == 100 {
		// redirect user to zarinpal
		// http.Redirect(w, r, zarinpalGateURL+resp.Authority, http.StatusFound)

		s, e := services.NewMellatService()

		rw.MelliNumber = user.MelliNumber
		rw.FreezBlIn, _ = strconv.ParseFloat(input.Amount, 64)
		rw.Authority = resp.Authority

		err := s.InsertPaymentGatewayLevl1(rw)

		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		if e != nil {
			return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
				"status": "server error",
			})
		}

		return ctx.Response().Json(http.StatusFound, map[string]string{"url": zarinpalGateURL + resp.Authority}) //"rw.Authority":   rw.Authority,
		//"rw.MelliNumber": rw.MelliNumber,
		//"rw.FreezBlIn":   strconv.FormatFloat(rw.FreezBlIn, 'f', -1, 64),

	}

	//http.Error(w, fmt.Sprintln("خطایی رخ داد:", resp.Status), http.StatusInternalServerError)
	return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": string(resp.Status)})
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

func (g *RialWalletController) VerifySharjHessab(ctx http.Context) http.Response {
	status := ctx.Request().Query("Status")
	authority := ctx.Request().Query("Authority")

	// بررسی وضعیت تراکنش
	if status != "OK" {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Transaction not approved",
		})
	}

	s, e := services.NewMellatService()
	if e != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": e.Error(),
		})
	}

	err := s.QueryPaymentGatewayLevl1()

	clinet := services.NewPaymentGatewayImplementationServicePortType("", false, nil)

	// Create a new payment request to Zarinpal
	resp, err := clinet.PaymentVerification(&services.PaymentVerification{
		MerchantID: merchantID,
		Amount:     1000,
		Authority:  authority,
	})

	// Check if response is error free
	if err != nil {
		//	http.Error(w, fmt.Sprintln("مشکلی در ارتباط رخ داد", err), http.StatusInternalServerError)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "resperror"})
	}
	if utils.KlDebug {
		fmt.Println(resp)
		fmt.Println(resp.RefID)
		fmt.Println(resp.Status)
	}

	switch resp.Status {
	case 100:

		// تراکنش با موفقیت تایید شد
		return ctx.Response().Json(http.StatusOK, map[string]string{
			"status": "OK",
			"RefID":  strconv.FormatInt(resp.RefID, 10),
		})

	case 101:
		// تراکنش قبلا تایید شده است
		return ctx.Response().Json(http.StatusOK, map[string]string{
			"status": "Transaction already verified",
		})

	case -1:
		// اطلاعات ارسال شده ناقص است
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Incomplete information sent",
		})

	case -2:
		// آی پی یا مرچنت کد اشتباه است
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Invalid IP or merchant code",
		})

	case -3:
		// مقدار قابل قبول نمی‌باشد
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Amount is invalid",
		})

	case -4:
		// سطح دسترسی تاجر اشتباه است
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Merchant access level is invalid",
		})

	case -11:
		// درخواست یافت نشد
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Request not found",
		})

	case -21:
		// هیچ نوع عملیات مالی برای این تراکنش یافت نشد
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "No financial operation for this transaction",
		})

	case -22:
		// تراکنش ناموفق است
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Transaction failed",
		})

	case -33:
		// مبلغ تراکنش با مبلغ پرداخت شده مطابقت ندارد
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Transaction amount mismatch",
		})

	case -34:
		// سقف تقسیم تراکنش از حد مجاز عبور نموده است
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Transaction split limit exceeded",
		})

	case -40:
		// دسترسی به متد مورد نظر امکان پذیر نمی‌باشد
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"error": "Access to the requested method is not allowed",
		})

	case -54:
		// درخواست با توجه به محدودیت زمانی قادر به پردازش نمی‌باشد
		return ctx.Response().Json(http.StatusRequestTimeout, map[string]string{
			"error": "Transaction request is expired",
		})

	case -101:
		// تراکنش لغو شده است
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Transaction canceled",
		})

	default:
		// خطای ناشناخته
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Unknown error occurred",
		})
	}

}
