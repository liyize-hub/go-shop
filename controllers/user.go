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
	var User models.User
	err := u.Ctx.ReadJSON(&User)
	if err != nil {
		utils.Logger.Error("ReadJSON failed", zap.Any("error", err))
		return utils.NewJSONResponse(u.Ctx, models.ErrorCode.ERROR, "", err)
	}

	// 数据参数检验
	if User.Name == "" || User.Pwd == "" {
		return utils.NewJSONResponse(u.Ctx,
			models.ErrorCode.LoginError,
			"用户名或密码为空,请重新填写后尝试登录",
			nil,
		)
	}

	//根据用户名、密码到数据库中查询对应的管理信息
	usr, err := u.Service.GetByUserNameAndPassword(User.Name, User.Pwd)
	if err != nil {
		utils.Logger.Error("GetByUserNameAndPassword error", zap.Any("err", err))
		return utils.NewJSONResponse(u.Ctx, models.ErrorCode.ERROR, err.Error(), "")
	}

	// 用户存在 设置session
	u.Session.Set("userID", strconv.FormatInt(usr.ID, 10))
	return utils.NewJSONResponse(u.Ctx,
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
func (u *UserController) PostRegister() mvc.Result {
	u.Ctx.Application().Logger().Info(" User register ")
	utils.Logger.Info("用户注册")
	var User models.User
	err := u.Ctx.ReadJSON(&User)
	if err != nil {
		utils.NewJSONResponse(u.Ctx, models.ErrorCode.ERROR, "", err)
	}
	// 数据参数检验
	if User.Name == "" || User.Pwd == "" {
		return utils.NewJSONResponse(u.Ctx,
			models.ErrorCode.LoginError,
			"用户名或密码为空,请重新填写后尝试登录",
			nil,
		)
	}

	// 写入数据库
	err = u.Service.AddUser(&User)
	if err != nil {
		utils.Logger.Error("InsertUser error", zap.Any("err", err))
		return utils.NewJSONResponse(u.Ctx, models.ErrorCode.ERROR, "添加用户失败", err)
	}

	return utils.NewJSONResponse(u.Ctx,
		models.ErrorCode.SUCCESS,
		"用户注册成功",
		nil,
	)
}
