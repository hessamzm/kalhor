package services

import (
	"errors"
	"fmt"
	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func MyJwt(ctx http.Context) (bool, string) {

	token := ctx.Request().Header("Authorization")
	// بررسی معتبر بودن توکن
	payload, err := facades.Auth(ctx).Parse(token)
	if errors.Is(err, auth.ErrorTokenExpired) {
		return false, "Token expired"
	}

	if payload == nil {
		fmt.Println(payload)
		//fmt.Println(payload.Key)
		fmt.Println(token)
	}

	if err != nil {

		return false, "Invalid token"

	}
	return true, "true"
}
