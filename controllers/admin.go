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
	var (
		admin models.Admin
	)

	err := ac.Ctx.ReadForm(&admin)
	if err != nil {
		utils.Logger.Error("ReadForm failed", zap.Any("error", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm failed", nil)
	}

	// 数据参数检验
	if admin.Name == "" || admin.Pwd == "" {
		return utils.NewJSONResponse(models.ErrorCode.LoginError, "用户名或密码为空,请重新填写后尝试登录", nil)
	}

	//根据用户名、密码到数据库中查询对应的管理信息
	ad, err := ac.Service.GetByAdminNameAndPassword(admin.Name, admin.Pwd)
	if err != nil {
		utils.Logger.Error("GetByAdminNameAndPassword error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	// 管理员存在 设置session
	ac.Session.Set("adminID", strconv.FormatInt(ad.ID, 10))

	//登录成功
	return utils.NewJSONResponse(
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
	var (
		admin models.Admin
	)

	err := ac.Ctx.ReadForm(&admin)
	if err != nil {
		utils.Logger.Error("ReadForm failed", zap.Any("error", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm err", nil)
	}

	// 写入数据库
	err = ac.Service.AddAdmin(&admin)
	if err != nil {
		utils.Logger.Error("InsertAdmin error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "添加管理员用户失败", nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "管理员注册成功", nil)

}

/**
 * 查询商铺信息
 * 接口：/admin/select
 * 方法：post
 */
func (p *AdminController) PostSelect() mvc.Result {
	p.Ctx.Application().Logger().Info(" search admin start")
	utils.Logger.Info("开始查询商Admin铺")

	var (
		admin models.Admin
	)

	err := p.Ctx.ReadForm(&admin)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm error", nil)
	}

	adminListandCount, err := p.Service.SelectShop(&admin)
	if err != nil {
		utils.Logger.Error("Selectadmin error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", adminListandCount)
}

/**
 * 删除商铺
 * 接口：/Admin/manager
 * 方法：post
 */
func (p *AdminController) GetDelete() {
	p.Ctx.Application().Logger().Info(" delete Admin start")
	utils.Logger.Info("开始删除商Admin铺")

	idString := p.Ctx.URLParam("id")
	if idString == "" {
		utils.Logger.Info("传入id为空")
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.Logger.Error("parseInt error", zap.Any("err", err))
		return
	}

	err = p.Service.DeleteShop(id)
	if err != nil {
		utils.Logger.Error("DeleteAdmin error", zap.Any("err", err))
	}
}

/**
 * 修改商铺
 * 接口：/Admin/update
 * 方法：post
 */
func (p *AdminController) PostUpdate() {
	p.Ctx.Application().Logger().Info(" update Admin ")
	utils.Logger.Info("更新商铺")
	var Admin models.Admin

	err := p.Ctx.ReadForm(&Admin)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
	}

	if Admin.ID <= 0 {
		utils.Logger.Info("更新数据不正确,id为0", zap.Any("Admin", Admin))
		return
	}

	err = p.Service.UpdateShop(&Admin)
	if err != nil {
		utils.Logger.Error("UpdateAdmin error", zap.Any("err", err))
	}
	utils.Logger.Info("更新成功", zap.Any("Admin", Admin))
	//p.Ctx.Redirect("/Admin/all") //跳转页面
}
