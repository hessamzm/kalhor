package routes

import (
	"github.com/goravel/framework/facades"
	"kalhor/app/http/controllers"
	"kalhor/app/http/middleware"
)

func Api() {

	priceController := controllers.NewPriceController()
	userController := controllers.UserController{}
	NotificationController := controllers.NotificationController{}
	TicketController := controllers.TicketController{}
	GoldWallet := controllers.GoldController{}
	RialController := controllers.RialWalletController{}

	//afterlogin := facades.Route().Middleware(middleware.Auth())

	//login
	facades.Route().Post("/login/isregister", userController.IsRegister)
	facades.Route().Post("/login", userController.Login)
	facades.Route().Post("/login/register", userController.Register)
	facades.Route().Post("/login/verifycode", userController.VerifyCode)

	//otp route
	facades.Route().Post("/askverifycode", userController.AskVerifyCode)
	facades.Route().Post("/verifycode", userController.VerifyCode)

	//after auth

	//info
	facades.Route().Middleware(middleware.Auth()).Get("/userinfo", userController.UserInfo)
	facades.Route().Middleware(middleware.Auth()).Get("/userwallet/:time-frame/:order/:type", userController.WalletInfo)
	facades.Route().Post("/send-message", NotificationController.SendMessage)
	facades.Route().Middleware(middleware.Auth()).Get("/rial", RialController.SharjHesab)

	//notif

	facades.Route().Post("/mark-as-seen", NotificationController.MarkAsSeen)
	facades.Route().Post("/store-admin-message", NotificationController.StoreAdminMessage)

	//ticket

	facades.Route().Put("/tickets", TicketController.Create)
	facades.Route().Get("/tickets", TicketController.Get)
	facades.Route().Put("/tickets/update", TicketController.Update)

	//goldwallet
	facades.Route().Post("/buygold", GoldWallet.ByGold)
	facades.Route().Post("/sellgold", GoldWallet.SellGold)

	//admin
	facades.Route().Get("/v1/price", priceController.Get)
	facades.Route().Put("/admin/price", priceController.Put)
	facades.Route().Put("/admin/price/status", priceController.ChStatus)
	facades.Route().Put("/admin/price/ograt", priceController.Putograt)
	facades.Route().Get("/auotmation/putauotmation", priceController.PutAuotmation)
}
