package controllers

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"kalhor/app/console"
	"kalhor/app/models"
	"kalhor/services"
	"kalhor/utils"
	"math/rand"
	"time"
)

type UserController struct{}

var user models.User
var Otp models.Otp
var input models.Input
var wg models.WalletGold
var wr models.WalletRial

func generateUniqueUserNum() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%07d", rand.Intn(9000)+1000) // تولید عدد 7 رقمی
}

func (c *UserController) Register(ctx http.Context) http.Response {
	user = models.User{}
	Otp = models.Otp{}

	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	facades.Orm().Query().Where("phone = ?", input.Phone).Get(&Otp)
	facades.Orm().Query().Where("phone = ?", input.Phone).First(&user)

	if user.Phone == input.Phone {
		return ctx.Response().Json(http.StatusConflict, map[string]string{"error": "Phone number already exists"})
	}

	facades.Orm().Query().Where("melli_number = ?", input.MelliNumber).First(&user)

	if user.MelliNumber == input.MelliNumber {
		return ctx.Response().Json(http.StatusConflict, map[string]string{"error": "Melli number already exists"})
	}
	if len(input.MelliNumber) != 10 {
		return ctx.Response().Json(http.StatusConflict, map[string]string{"error": "Melli number must have 10 numeric"})
	}

	if Otp.Status == false {
		return ctx.Response().Json(http.StatusConflict, map[string]string{"error": "Otp code not valid"})
	}

	// تولید یک UserNum یکتا
	var newUserNum string
	for {
		newUserNum = generateUniqueUserNum()
		var existingUser models.User
		facades.Orm().Query().Where("user_num = ?", newUserNum).First(&existingUser)
		if existingUser.ID == 0 {
			break // اگر این UserNum در دیتابیس وجود ندارد، آن را انتخاب کنید
		}
	}

	datetavalod, err := time.Parse("2006/01/02", input.TarikhTavalod)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "Invalid date format"})
	}

	// ایجاد کاربر جدید
	newUser := models.User{
		Name:          input.Name,
		MelliNumber:   input.MelliNumber,
		Phone:         input.Phone,
		TarikhTavalod: datetavalod,
		UserNum:       newUserNum,
	}
	user.OtpCode = true
	if err := facades.Orm().Query().Create(&newUser); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "Error creating user"})
	}

	token, err := facades.Auth(ctx).LoginUsingID(newUser.ID)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "Error logging in"})
	}

	Otp.OtpCode = ""
	Otp.Step = 0
	Otp.Phone = input.Phone
	Otp.Status = false
	facades.Orm().Query().Save(&Otp)
	return ctx.Response().Header("Authorization", token).Json(http.StatusCreated, map[string]string{"user": newUserNum})

}

func (c *UserController) Login(ctx http.Context) http.Response {
	user = models.User{}
	Otp = models.Otp{}
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if utils.KlDebug {
		fmt.Println("Phone input from login page:", input.Phone)
		fmt.Println("code input from login page:", input.Code)
	}

	// پیدا کردن کاربر بر اساس شماره تلفن

	facades.Orm().Query().Where("phone = ?", input.Phone).Get(&user)

	if user.ID == 0 {
		return ctx.Response().Json(http.StatusFailedDependency, map[string]string{"error": "Invalid credentials"})
	} else if user.Phone != input.Phone {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{"error": "user not found"})
	}
	facades.Orm().Query().Where("phone = ?", input.Phone).Get(&Otp)
	if utils.KlDebug {
		fmt.Println("opt code", user.OtpCode)
	}
	if Otp.OtpCode != input.Code {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": "Invalid verification code"})
	}

	if user.Freez == true {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": "user is freez"})
	}

	Otp.OtpCode = ""
	Otp.Step = 0
	user.OtpCode = true
	facades.Orm().Query().Save(&user)
	facades.Orm().Query().Save(&Otp)
	token, err := facades.Auth(ctx).LoginUsingID(user.ID)
	if utils.KlDebug {
		fmt.Println("jwt eeror", err)
	}
	return ctx.Response().Header("Authorization", token).Json(http.StatusOK, map[string]any{
		"status": "success",
	})
}

func (c *UserController) IsRegister(ctx http.Context) http.Response {
	user = models.User{}
	Otp = models.Otp{}
	input = models.Input{}

	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"status": "badrequest",
			"error":  "Invalid input"})
	}
	if input.Phone == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"status": "badrequest",
			"error":  "Phone number required"})
	}

	if utils.KlDebug {
		fmt.Println("Phone input from login page:", input.Phone)
	}

	// پیدا کردن کاربر بر اساس شماره تلفن

	facades.Orm().Query().Where("phone = ?", input.Phone).Get(&user)

	otpCode := generateOtp()

	// ارسال پیامک
	response, err := services.SendByBaseNumber("9216318161", "Y@N!0", otpCode, input.Phone, 246010)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"status": "server-error-checkphonenumber",
			"error":  "SMS sending failed"})
	}

	if utils.KlDebug {
		fmt.Println("sms response: ", response)
	}

	if user.Phone != input.Phone {
		Otp.Step = 1
		Otp.Phone = input.Phone
		Otp.OtpCode = otpCode
		Otp.UpdatedAt = time.Now()

		facades.Orm().Query().UpdateOrCreate(&Otp, Otp, otpCode)
		return ctx.Response().Json(http.StatusNotFound, map[string]string{"status": "user not found",
			"sms":   "send",
			"go-to": "regestry",
		})
	}

	// به‌روزرسانی Otp
	Otp.Step = 1
	Otp.Phone = input.Phone
	Otp.OtpCode = otpCode
	Otp.UpdatedAt = time.Now()

	// ذخیره Otp در دیتابیس
	facades.Orm().Query().Update(&Otp)

	if utils.KlDebug {
		fmt.Println("Otp response: ", Otp)
	}

	return ctx.Response().Json(http.StatusOK, map[string]string{"status": "user-finde",
		"goto": "login",
		"sms":  "send"})
}

func (c *UserController) VerifyCode(ctx http.Context) http.Response {

	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	Otp = models.Otp{}

	// جستجو در دیتابیس بر اساس phone
	facades.Orm().Query().Where("phone = ?", input.Phone).Get(&Otp)

	// جستجو در دیتابیس بر اساس optcode

	if utils.KlDebug {
		fmt.Println("user otpcode :", Otp.OtpCode)
	}

	// بررسی اینکه کاربری پیدا شده یا نه
	//if user.ID == 0 {
	// 09214295835
	//	return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": "Invalid verification code"})
	//}

	// بررسی کد تایید
	if input.Code == Otp.OtpCode {
		Otp.Step = 0
		Otp.Status = true
		Otp.OtpCode = ""
		facades.Orm().Query().Save(&Otp)
		facades.Orm().Query().Update(&user)
		return ctx.Response().Json(http.StatusOK, map[string]string{"message": "Verification successful"})
	} else {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": "Invalid verification code"})
	}
}

func (c *UserController) AskVerifyCode(ctx http.Context) http.Response {
	// دریافت ورودی
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// جستجو در جدول کاربران
	facades.Orm().Query().Where("phone = ?", input.Phone).First(&user)

	// ریست کردن Otp قبل از هر جستجو
	//Otp.UpdatedAt = time.Now().Add(120)
	Otp = models.Otp{}
	// جستجو در جدول otps
	facades.Orm().Query().Where("phone = ?", input.Phone).First(&Otp)

	// اگر رکوردی برای شماره تلفن پیدا نشد، یک Otp جدید ایجاد کن
	// محاسبه اختلاف زمان
	timeDifference := time.Now().Sub(Otp.UpdatedAt)

	if utils.KlDebug {
		fmt.Println("user models :", user)
		fmt.Println("Otp models :", Otp)
		fmt.Println("timeDifference  :", timeDifference)
	}

	// تعیین گام Otp
	if timeDifference.Seconds() < 120 {
		Otp.Step = 1
	} else if timeDifference.Seconds() > 120 {
		Otp.Step = 0
	} else if user.Freez == true {
		Otp.Step = 100
	}

	if utils.KlDebug {
		fmt.Println("Otp.step :", Otp.Step)
	}

	// بررسی مراحل Otp
	switch Otp.Step {

	case 1:
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"try later": time.Now().Sub(Otp.UpdatedAt).String()})

	case 100:
		return ctx.Response().Json(http.StatusLocked, map[string]string{"you are block ": "wait 1 day"})

	case 0:
		// تولید کد Otp
		otpCode := generateOtp()

		// ارسال پیامک
		response, err := services.SendByBaseNumber("9216318161", "Y@N!0", otpCode, input.Phone, 246010)
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "SMS sending failed"})
		}

		if utils.KlDebug {
			fmt.Println("sms response: ", response)
		}

		// به‌روزرسانی Otp
		Otp.Step = 1
		Otp.Phone = input.Phone
		Otp.OtpCode = otpCode
		Otp.UpdatedAt = time.Now()

		// ذخیره Otp در دیتابیس
		facades.Orm().Query().Save(&Otp)

		if utils.KlDebug {
			fmt.Println("Otp response: ", Otp)
		}

		return ctx.Response().Json(http.StatusOK, map[string]string{"sms send to": input.Phone})
	}

	return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": "SMS failed to send"})
}

func (c *UserController) UserInfo(ctx http.Context) http.Response {
	user = models.User{}
	// اعتبارسنجی JWT
	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}

	// دریافت اطلاعات کاربر
	facades.Auth(ctx).User(&user)

	// جستجوی اطلاعات کاربر در دیتابیس
	facades.Orm().Query().Where("phone = ?", user.Phone).Get(&user)

	// دریافت اطلاعات WalletGold
	mapstring := console.UserInfo(user)

	return ctx.Response().Json(http.StatusOK, mapstring)
}

func (c *UserController) WalletInfo(ctx http.Context) http.Response {
	user := models.User{}
	// اعتبارسنجی JWT
	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}

	// دریافت اطلاعات کاربر
	facades.Auth(ctx).User(&user)
	facades.Orm().Query().Where("phone = ?", user.Phone).Get(&user)

	// دریافت پارامترهای مسیر
	timeFrame := ctx.Request().Route("time-frame")
	order := ctx.Request().Route("order")
	walletType := ctx.Request().Route("type")

	var query string

	// تنظیم شرط زمانی بر اساس پارامتر time-frame
	timeCondition := ""
	switch timeFrame {
	case "week":
		timeCondition = "AND event_time >= now() - INTERVAL 1 WEEK"
	case "month":
		timeCondition = "AND event_time >= now() - INTERVAL 1 MONTH"
	case "year":
		timeCondition = "AND event_time >= now() - INTERVAL 1 YEAR"
	default:
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"status": "badrequest",
			"error":  "Invalid time frame",
		})
	}

	// تنظیم order بر اساس خرید (BalanceIn) یا فروش (BalanceOut)
	orderCondition := ""
	switch order {
	case "sell":
		orderCondition = "AND balance_out > 0"
	case "buy":
		orderCondition = "AND balance_in > 0"
	default:
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"status": "badrequest",
			"error":  "Invalid order type",
		})
	}
	kal := "kal" + user.MelliNumber
	// تنظیم جدول بر اساس پارامتر type
	switch walletType {
	case "gold":
		//var wallets models.WalletGold
		query = fmt.Sprintf("SELECT * FROM wallet_gold WHERE melli_number = '%s' %s %s ORDER BY event_time DESC", kal, timeCondition, orderCondition)
		s, e := services.NewWalletService()

		wallets, e := s.Queryforgold(query)
		fmt.Printf("e", e)
		fmt.Printf("wallets 6", wallets)
		fmt.Println(wallets)
		if e != nil {
			return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
				"status": "badrequest",
			})
		}
		return ctx.Response().Json(http.StatusOK, wallets)

	case "rial":

		query = fmt.Sprintf("SELECT * FROM wallet_rial WHERE melli_number = '%s' %s %s ORDER BY event_time DESC", kal, timeCondition, orderCondition)

		s, e := services.NewWalletServiceRial()

		wallets, e := s.Queryforrial(query)
		fmt.Printf("e", e)
		fmt.Printf("wallets 6", wallets)
		if e != nil {
			return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
				"status": "badrequest",
			})
		}
		return ctx.Response().Json(http.StatusOK, wallets)
	default:
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"status": "badrequest",
			"error":  "Invalid wallet type",
		})
	}

}

func generateOtp() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // تولید یک عدد شش رقمی
}
