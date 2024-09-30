package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"kalhor/app/console"
	"kalhor/app/models"
	"kalhor/services"
	"time"
)

type NotificationController struct{}

var model models.Notification

func (c *NotificationController) SendMessage(ctx http.Context) http.Response {

	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}

	// دریافت اطلاعات کاربر
	facades.Auth(ctx).User(&user)

	model.ToWho = user.Phone
	mapstring := console.UserNotif(model.ToWho)
	//if mapstring == nil {
	//	return ctx.Response().Json(http.StatusNotFound, map[string]string{
	//		"error": "no massages fund",
	//	})
	//}

	return ctx.Response().Json(http.StatusOK, mapstring)
}
func (c *NotificationController) MarkAsSeen(ctx http.Context) http.Response {

	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}

	id := ctx.Request().Input("id")

	var notification models.Notification
	if err := facades.Orm().Query().Find(&notification, id); err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Notification not found",
		})
	}

	notification.IsSee = true
	if err := facades.Orm().Query().Save(&notification); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update notification",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": "Notification marked as seen",
	})
}

func (c *NotificationController) StoreAdminMessage(ctx http.Context) http.Response {

	to := ctx.Request().Input("to")
	subject := ctx.Request().Input("subject")
	message := ctx.Request().Input("message")

	notification := models.Notification{
		ToWho:     to,
		Subject:   subject,
		Messages:  message,
		IsSee:     false,
		CreatedAt: time.Now(),
	}

	if err := facades.Orm().Query().Create(&notification); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to store admin message",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": "Admin message stored successfully",
	})
}
