package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"kalhor/app/models"
	"kalhor/services"
	"time"
)

type TicketController struct {
	// کنترلر ساختن
}

// ایجاد یک تیکت جدید
func (t *TicketController) Create(ctx http.Context) http.Response {
	var ticket models.Ticket

	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}

	// دریافت اطلاعات کاربر
	facades.Auth(ctx).User(&user)

	if err := ctx.Request().Bind(&ticket); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": err.Error(),
		})
	}

	ticket.Status = "open"
	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()

	ticket.UserID = user.ID

	if err := facades.Orm().Query().Create(&ticket); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to create ticket",
		})
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"ticket": ticket,
	})
}

// دریافت تمام تیکت‌ها
func (t *TicketController) Get(ctx http.Context) http.Response {
	var tickets []models.Ticket

	istrue, err := services.MyJwt(ctx)
	if !istrue {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{"error": err})
	}

	// دریافت اطلاعات کاربر
	facades.Auth(ctx).User(&user)

	facades.Orm().Query().Where("user_id", user.ID).Get(&tickets)
	if tickets == nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{"error": "Ticket not found"})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"tickets": tickets,
	})
}

// بروز رسانی یک تیکت
func (t *TicketController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	var ticket models.Ticket

	if err := facades.Orm().Query().Where("id", id).First(&ticket); err != nil {
		ctx.Response().Json(http.StatusNotFound, http.Json{
			"error": "Ticket not found",
		})
		return nil
	}

	if err := ctx.Request().Bind(&ticket); err != nil {
		ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": err.Error(),
		})
		return nil
	}

	ticket.UpdatedAt = time.Now()
	if err := facades.Orm().Query().Save(&ticket); err != nil {
		ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to update ticket",
		})
		return nil
	}

	ctx.Response().Json(http.StatusOK, http.Json{
		"ticket": ticket,
	})
	return nil
}
