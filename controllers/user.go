package controllers

import (
	"encoding/json"
	"go-shop/config"
	"go-shop/models"
	"go-shop/services"
	"go-shop/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
)

type UserController struct {
	Ctx     iris.Context         // iris框架自动为每个请求都绑定上下文对象
	Service services.UserService // User功能实体
}

/**
 * 检查用户服务端登录状态
 * 接口：/user/check
 * 方法：get
 */
func (u *UserController) GetCheck() {
	token := u.Ctx.URLParam("token")
	userID := u.Ctx.URLParam("id")

	id, err := u.Service.GetToken(token)
	if err != nil {
		utils.Logger.Info("用户服务端登录过期", zap.Any("UserID", userID), zap.Any("error", err))
		utils.SendJSON(u.Ctx, models.ErrorCode.NotFound, "用户服务端登录过期", nil)
	}
	if id == userID {
		utils.SendJSON(u.Ctx, models.ErrorCode.SUCCESS, "", nil)
	} else {
		utils.SendJSON(u.Ctx, models.ErrorCode.NotFound, "用户服务端登录过期", nil)
	}
}

/**
 * 检查用户是否注册、用户登录
 * 接口：/user/login
 * 方法：get
 */
func (u *UserController) GetLogin() {
	client := &http.Client{}
	var result map[string]string

	// 获取登录状态
	code := u.Ctx.URLParam("code")
	req, err := http.NewRequest("GET", "https://api.weixin.qq.com/sns/jscode2session", nil)
	if err != nil {
		utils.Logger.Error("NewRequest error", zap.Any("err", err))
		return
	}
	query := req.URL.Query()
	query.Add("appid", config.ServerConfig.WeChatAppID)
	query.Add("secret", config.ServerConfig.WeChatSecret)
	query.Add("js_code", code)
	query.Add("grant_type", "authorization_code")
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		utils.Logger.Error("client.Do failed", zap.Any("err", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.Error("ReadAll failed", zap.Any("err", err))
		return
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		utils.Logger.Error("Unmarshal failed", zap.Any("err", err))
		return
	}
	openID := result["openid"]
	if err != nil {
		utils.Logger.Error("ParseInt failed", zap.Any("err", err))
		return
	}

	user, exist, err := u.Service.SelectUser(&models.User{OpenID: openID, Flag: 0})
	if err != nil {
		utils.Logger.Error("SelectUserByID failed", zap.Any("err", err))
		return
	}

	if !exist {
		token := utils.RandToken()
		//用户注册保留10分钟的注册时间 token 保存openid值
		err := u.Service.SetToken(token, openID, time.Minute*10)
		if err != nil {
			utils.Logger.Error("SetToken failed", zap.Any("err", err))
		}
		utils.SendJSON(u.Ctx, models.ErrorCode.NotFound, "用户不存在", token)
		return
	}

	//openID对用户不可见
	user.OpenID = ""
	user.ID = 0
	//token为用户服务端登录态标识
	user.Token = utils.RandToken()
	err = u.Service.SetToken(user.Token, user.ID, time.Hour*2) //设置2个小时的用户登录状态 key:token value:用户ID
	if err != nil {
		utils.Logger.Error("SetToken failed", zap.Any("err", err))
		return
	}
	utils.SendJSON(u.Ctx, models.ErrorCode.SUCCESS, "用户登录成功", user)
}

/**
 * 用户注册功能
 * 接口：/user/register
 * 方法: post
 */
func (u *UserController) PostRegister() {
	u.Ctx.Application().Logger().Info(" User register ")
	utils.Logger.Info("用户注册")
	var (
		user models.User
	)

	err := u.Ctx.ReadJSON(&user)
	if err != nil {
		utils.Logger.Error("ReadForm failed", zap.Any("error", err))
		utils.SendJSON(u.Ctx, models.ErrorCode.ERROR, "", nil)
	}
	if user.Token == "" {
		utils.Logger.Info("无token", zap.Any("User", user))
		utils.SendJSON(u.Ctx, models.ErrorCode.NotFound, "微信没有登录", nil)
	}

	openid, err := u.Service.GetToken(user.Token)
	if err != nil {
		utils.Logger.Info("用户注册过期", zap.Any("User", user), zap.Any("error", err))
		utils.SendJSON(u.Ctx, models.ErrorCode.NotFound, "用户注册过期", nil)
	}

	//从redis获取openID
	user.OpenID = openid

	// 写入数据库
	err = u.Service.AddUser(&user)
	if err != nil {
		utils.Logger.Error("InsertUser error", zap.Any("err", err))
		utils.SendJSON(u.Ctx, models.ErrorCode.ERROR, "添加用户失败", nil)
	}

	utils.SendJSON(u.Ctx, models.ErrorCode.SUCCESS, "用户登录成功", nil)
}

/**
 * 删除用户
 * 接口：/user/manager
 * 方法：post
 */
func (u *UserController) GetDelete() {
	u.Ctx.Application().Logger().Info(" delete User start")
	utils.Logger.Info("开始删除商User铺")

	idString := u.Ctx.URLParam("id")
	if idString == "" {
		utils.Logger.Info("传入id为空")
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.Logger.Error("parseInt error", zap.Any("err", err))
		return
	}

	err = u.Service.DeleteUser(id)
	if err != nil {
		utils.Logger.Error("DeleteUser error", zap.Any("err", err))
	}
}

/**
 * 修改用户
 * 接口：/user/update
 * 方法：post
 */
func (u *UserController) PostUpdate() {
	u.Ctx.Application().Logger().Info(" update User ")
	utils.Logger.Info("更新用户")
	var User models.User

	err := u.Ctx.ReadForm(&User)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
	}

	if User.ID <= 0 {
		utils.Logger.Info("更新数据不正确,id为0", zap.Any("User", User))
		return
	}

	err = u.Service.UpdateUser(&User)
	if err != nil {
		utils.Logger.Error("UpdateUser error", zap.Any("err", err))
	}
	utils.Logger.Info("更新成功", zap.Any("User", User))
	//p.Ctx.Redirect("/User/all") //跳转页面
}

/**
 * 查询用户信息
 * 接口：/user/select
 * 方法：post
 */
func (u *UserController) PostSelect() mvc.Result {
	u.Ctx.Application().Logger().Info(" search User start")
	utils.Logger.Info("开始查询用户")

	var (
		User models.User
	)

	err := u.Ctx.ReadForm(&User)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm error", nil)
	}

	UserListandCount, err := u.Service.SelectUsers(&User)
	if err != nil {
		utils.Logger.Error("SelectUser error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", UserListandCount)
}
