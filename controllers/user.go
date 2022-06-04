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

type UserController struct {
	Ctx     iris.Context         // iris框架自动为每个请求都绑定上下文对象
	Service services.UserService // User功能实体
	Session *sessions.Session    // session对象
}

/**
 * 用户登录功能
 * 接口：/User/login
 * 方法：post
 */
func (ac *UserController) PostLogin() mvc.Result {
	ac.Ctx.Application().Logger().Info(" User login ")
	utils.Logger.Info("用户登录")
	var (
		User models.User
	)

	err := ac.Ctx.ReadForm(&User)
	if err != nil {
		utils.Logger.Error("ReadForm failed", zap.Any("error", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm failed", nil)
	}

	// 数据参数检验
	if User.Name == "" || User.Pwd == "" {
		return utils.NewJSONResponse(models.ErrorCode.LoginError, "用户名或密码为空,请重新填写后尝试登录", nil)
	}

	//根据用户名、密码到数据库中查询对应的管理信息
	ad, err := ac.Service.GetByUserNameAndPassword(User.Name, User.Pwd)
	if err != nil {
		utils.Logger.Error("GetByUserNameAndPassword error", zap.Any("err", err))
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

	// 用户存在 设置session
	//ac.Session.Set("UserID", strconv.FormatInt(ad.ID, 10))

	//登录成功
	return utils.NewJSONResponse(
		models.ErrorCode.SUCCESS,
		"用户登录成功",
		nil,
	)
}

/**
 * 用户注册功能
 * 接口：/User/register
 * 方法: post
 */
func (ac *UserController) PostRegister() mvc.Result {
	ac.Ctx.Application().Logger().Info(" User register ")
	utils.Logger.Info("用户注册")
	var (
		User models.User
	)

	err := ac.Ctx.ReadForm(&User)
	if err != nil {
		utils.Logger.Error("ReadForm failed", zap.Any("error", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm err", nil)
	}

	// 写入数据库
	err = ac.Service.AddUser(&User)
	if err != nil {
		utils.Logger.Error("InsertUser error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "添加用户用户失败", nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "用户注册成功", nil)
}

/**
 * 删除商铺
 * 接口：/User/manager
 * 方法：post
 */
func (ac *UserController) GetDelete() {
	ac.Ctx.Application().Logger().Info(" delete User start")
	utils.Logger.Info("开始删除商User铺")

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

	err = ac.Service.DeleteUser(id)
	if err != nil {
		utils.Logger.Error("DeleteUser error", zap.Any("err", err))
	}
}

/**
 * 修改商铺
 * 接口：/User/update
 * 方法：post
 */
func (ac *UserController) PostUpdate() {
	ac.Ctx.Application().Logger().Info(" update User ")
	utils.Logger.Info("更新商铺")
	var User models.User

	err := ac.Ctx.ReadForm(&User)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
	}

	if User.ID <= 0 {
		utils.Logger.Info("更新数据不正确,id为0", zap.Any("User", User))
		return
	}

	err = ac.Service.UpdateUser(&User)
	if err != nil {
		utils.Logger.Error("UpdateUser error", zap.Any("err", err))
	}
	utils.Logger.Info("更新成功", zap.Any("User", User))
	//p.Ctx.Redirect("/User/all") //跳转页面
}

/**
 * 查询商铺信息
 * 接口：/User/select
 * 方法：post
 */
func (ac *UserController) PostSelect() mvc.Result {
	ac.Ctx.Application().Logger().Info(" search User start")
	utils.Logger.Info("开始查询商铺")

	var (
		User models.User
	)

	err := ac.Ctx.ReadForm(&User)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm error", nil)
	}

	//从cookie中获取具体商铺
	uid := ac.Ctx.GetCookie("uid")
	if uid != "1" {
		User.ID, _ = strconv.ParseInt(uid, 10, 64)
	}

	UserListandCount, err := ac.Service.SelectUser(&User)
	if err != nil {
		utils.Logger.Error("SelectUser error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", UserListandCount)
}
