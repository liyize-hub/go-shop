package controllers

import (
	"go-shop/models"
	"go-shop/services"
	"go-shop/utils"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"go.uber.org/zap"
)

type AdminController struct {
	Ctx     iris.Context          // iris框架自动为每个请求都绑定上下文对象
	Service services.AdminService // admin功能实体
	Session *sessions.Session     // session对象
}

/**
 * 管理员登录功能
 * 接口：/admin/login
 * 方法：post
 */
func (ac *AdminController) PostLogin() mvc.Result {
	ac.Ctx.Application().Logger().Info(" admin login ")
	utils.Logger.Info("管理员登录")
	var admin models.Admin
	err := ac.Ctx.ReadJSON(&admin)
	if err != nil {
		utils.Logger.Error("ReadJSON failed", zap.Any("error", err))
		return utils.NewJSONResponse(ac.Ctx, models.ErrorCode.ERROR, "", err)
	}

	// 数据参数检验
	if admin.Name == "" || admin.Pwd == "" {
		return utils.NewJSONResponse(ac.Ctx,
			models.ErrorCode.LoginError,
			"用户名或密码为空,请重新填写后尝试登录",
			nil,
		)
	}

	//根据用户名、密码到数据库中查询对应的管理信息
	ad, err := ac.Service.GetByAdminNameAndPassword(admin.Name, admin.Pwd)
	if err != nil {
		utils.Logger.Error("GetByAdminNameAndPassword error", zap.Any("err", err))
		return utils.NewJSONResponse(ac.Ctx, models.ErrorCode.ERROR, err.Error(), "")
	}

	// 管理员存在 设置session
	ac.Session.Set("adminID", strconv.FormatInt(ad.ID, 10))
	return utils.NewJSONResponse(ac.Ctx,
		models.ErrorCode.SUCCESS,
		"管理员登录成功",
		nil,
	)
}

/**
 * 管理员注册功能
 * 接口：/admin/register
 * 方法: post
 */
func (ac *AdminController) PostRegister() mvc.Result {
	ac.Ctx.Application().Logger().Info(" admin register ")
	utils.Logger.Info("管理员注册")
	var admin models.Admin
	err := ac.Ctx.ReadJSON(&admin)
	if err != nil {
		utils.NewJSONResponse(ac.Ctx, models.ErrorCode.ERROR, "", err)
	}
	// 数据参数检验
	if admin.Name == "" || admin.Pwd == "" {
		return utils.NewJSONResponse(ac.Ctx,
			models.ErrorCode.LoginError,
			"用户名或密码为空,请重新填写后尝试登录",
			nil,
		)
	}

	// 写入数据库
	err = ac.Service.AddAdmin(&admin)
	if err != nil {
		utils.Logger.Error("InsertAdmin error", zap.Any("err", err))
		return utils.NewJSONResponse(ac.Ctx, models.ErrorCode.ERROR, "添加管理员用户失败", err)
	}

	return utils.NewJSONResponse(ac.Ctx,
		models.ErrorCode.SUCCESS,
		"管理员注册成功",
		nil,
	)
}
