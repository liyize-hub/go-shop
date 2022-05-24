package middleware

import (
	"go-shop/utils"
	"strings"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

func AuthConProduct(ctx iris.Context) {

	uid := ctx.GetCookie("uid")
	sign := ctx.GetCookie("sign")
	if uid == "" || sign == "" {
		utils.Logger.Info("请先登录！")
		ctx.Redirect("/")
		return
	}

	uidByte, err := utils.DePwdCode(strings.Replace(sign, " ", "+", 1))
	if err != nil {
		utils.Logger.Error("EnPwdCode error", zap.Any("err", err))
		return
	}

	if uid == string(uidByte) {
		utils.Logger.Info("cookie验证成功,登录成功!")
		ctx.Next()
	}

	return
}
