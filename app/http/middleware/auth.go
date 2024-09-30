package middleware

import (
	"fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"kalhor/utils"
)

func Auth() http.Middleware {
	return func(ctx http.Context) {
		// دریافت توکن از هدر Authorization
		token := ctx.Request().Header("Authorization")

		if token == "" {

			ctx.Response().Json(http.StatusUnauthorized, map[string]string{
				"error": "Token not provided",
			})
		}

		// بررسی معتبر بودن توکن
		_, err := facades.Auth(ctx).Parse(token)
		if err != nil {

			ctx.Response().Json(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token",
			})
			return
		}

		// فقط برای debug
		if utils.KlDebug {
			fmt.Println("token:", token)
		}

		// ادامه پردازش درخواست
		// از ctx.Request().Next() استفاده کنید
	}
}
