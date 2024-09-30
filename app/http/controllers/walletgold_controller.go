package controllers

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"kalhor/app/models"
	"kalhor/services"
	"kalhor/utils"
)

var inputgold struct {
}

type GoldwalletController struct {
	walletService *services.WalletService
}

func NewGoldwallet() *GoldwalletController {
	walletService, err := services.NewWalletService()
	if err != nil {
		return nil
	}
	return &GoldwalletController{
		walletService: walletService,
	}
}

func (r *GoldwalletController) Income(ctx http.Context) http.Response {
	var gw models.WalletGold
	user = models.User{}
	istrue, strerr := services.MyJwt(ctx)

	if istrue == false {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": strerr})
	}

	facades.Auth(ctx).User(&user)

	if utils.KlDebug {
		fmt.Println("user info :", user)
	}
	facades.Orm().Query().Where("phone = ?", user.Phone).First(&user)

	if err := ctx.Request().Bind(&inputgold); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// استفاده از سرویس برای درج اطلاعات
	err := r.walletService.InsertWalletGold(&gw)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "insert error",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"msg": "Wallet record inserted successfully",
	})
}
