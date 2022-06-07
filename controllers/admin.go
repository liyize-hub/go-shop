package controllers

import (
	"go-shop/models"
	"go-shop/services"
	"go-shop/utils"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
)

type AdminController struct {
	Ctx     iris.Context          // iris框架自动为每个请求都绑定上下文对象
	Service services.AdminService // admin功能实体
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

	//写入用户ID到cookie中
	utils.GlobalCookie(ac.Ctx, "uid", strconv.FormatInt(ad.ID, 10))

	uidByte := []byte(strconv.FormatInt(ad.ID, 10))
	uidString, err := utils.EnPwdCode(uidByte)
	if err != nil {
		utils.Logger.Error("EnPwdCode error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}
	//写入加密后的用户id到cookie
	utils.GlobalCookie(ac.Ctx, "sign", uidString)

	// 管理员存在 设置session
	//ac.Session.Set("adminID", strconv.FormatInt(ad.ID, 10))

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
 * 删除商铺
 * 接口：/Admin/manager
 * 方法：post
 */
func (ac *AdminController) GetDelete() {
	ac.Ctx.Application().Logger().Info(" delete Admin start")
	utils.Logger.Info("开始删除商Admin铺")

	idString := ac.Ctx.URLParam("id")
	if idString == "" {
		utils.Logger.Info("传入id为空")
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.Logger.Error("parseInt error", zap.Any("err", err))
		return
	}

	err = ac.Service.DeleteShop(id)
	if err != nil {
		utils.Logger.Error("DeleteAdmin error", zap.Any("err", err))
	}
}

/**
 * 修改商铺
 * 接口：/Admin/update
 * 方法：post
 */
func (ac *AdminController) PostUpdate() {
	ac.Ctx.Application().Logger().Info(" update Admin ")
	utils.Logger.Info("更新商铺")
	var Admin models.Admin

	err := ac.Ctx.ReadForm(&Admin)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
	}

	if Admin.ID <= 0 {
		utils.Logger.Info("更新数据不正确,id为0", zap.Any("Admin", Admin))
		return
	}

	err = ac.Service.UpdateShop(&Admin)
	if err != nil {
		utils.Logger.Error("UpdateAdmin error", zap.Any("err", err))
	}
	utils.Logger.Info("更新成功", zap.Any("Admin", Admin))
	//p.Ctx.Redirect("/Admin/all") //跳转页面
}

/**
 * 查询商铺信息
 * 接口：/admin/select
 * 方法：post
 */
func (ac *AdminController) PostSelect() mvc.Result {
	ac.Ctx.Application().Logger().Info(" search admin start")
	utils.Logger.Info("开始查询商铺")

	var (
		admin models.Admin
	)

	err := ac.Ctx.ReadForm(&admin)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm error", nil)
	}

	//从cookie中获取具体商铺
	uid := ac.Ctx.GetCookie("uid")
	if uid != "1" {
		admin.ID, _ = strconv.ParseInt(uid, 10, 64)
	}

	adminListandCount, err := ac.Service.SelectShop(&admin)
	if err != nil {
		utils.Logger.Error("Selectadmin error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", adminListandCount)
}
