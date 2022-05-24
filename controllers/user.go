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
func (u *UserController) PostLogin() mvc.Result {
	u.Ctx.Application().Logger().Info(" User login ")
	utils.Logger.Info("用户登录")

	var (
		user     models.User
	)

	err := u.Ctx.ReadJSON(&user)
	if err != nil {
		utils.Logger.Error("ReadJSON failed", zap.Any("error", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadJSON err", nil)
	}

	// 数据参数检验
	if user.Name == "" || user.Pwd == "" {
		return utils.NewJSONResponse(models.ErrorCode.LoginError, "用户名或密码为空,请重新填写后尝试登录", nil)
	}

	//根据用户名、密码到数据库中查询对应的管理信息
	usr, err := u.Service.GetByUserNameAndPassword(user.Name, user.Pwd)
	if err != nil {
		utils.Logger.Error("GetByUserNameAndPassword error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	// 用户存在 设置session
	u.Session.Set("userID", strconv.FormatInt(usr.ID, 10))

	//登录成功
	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "用户登录成功", nil)
}

/**
 * 用户注册功能
 * 接口：/User/register
 * 方法: post
 */
func (u *UserController) PostRegister() mvc.Result {
	u.Ctx.Application().Logger().Info(" User register ")
	utils.Logger.Info("用户注册")

	var (
		user     models.User
	)

	err := u.Ctx.ReadJSON(&user)
	if err != nil {
		utils.Logger.Error("ReadJSON failed", zap.Any("error", err))
		utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadJSON err", nil)
	}

	// 写入数据库
	err = u.Service.AddUser(&user)
	if err != nil {
		utils.Logger.Error("InsertUser error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "添加用户失败", nil)
	}

	//注册成功
	return utils.NewJSONResponse(models.ErrorCode.SUCCESS,"用户注册成功",nil)

}
