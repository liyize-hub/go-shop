package controllers

import (
	"encoding/json"
	"go-shop/models"
	"go-shop/services"
	"go-shop/utils"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
)

type UserController struct {
	Ctx     iris.Context         // iris框架自动为每个请求都绑定上下文对象
	Service services.UserService // User功能实体
}

/**
 * 检查用户是否已注册
 * 接口：/user/
 * 方法：post
 */
func (u *UserController) GetCheck() {
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
	query.Add("appid", "wxd7cb4fcb6da2a7e3")
	query.Add("secret", "6884321740e25ea18d1a347d6ecb7c2b")
	query.Add("js_code", code)
	query.Add("grant_type", "authorization_code")
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		utils.Logger.Error("Unmarshal failed", zap.Any("err", err))
	}

	utils.SendJSON(u.Ctx, models.ErrorCode.SUCCESS, "检测成功", result["openid"])

}

/**
 * 用户登录功能
 * 接口：/user/login
 * 方法：post
 */
func (u *UserController) PostLogin() mvc.Result {
	u.Ctx.Application().Logger().Info(" User login ")
	utils.Logger.Info("用户登录")
	var (
		User models.User
	)

	err := u.Ctx.ReadForm(&User)
	if err != nil {
		utils.Logger.Error("ReadForm failed", zap.Any("error", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm failed", nil)
	}

	//根据用户名、密码到数据库中查询对应的管理信息
	ad, err := u.Service.GetByUserNameAndPassword(User.Name, User.Pwd)
	if err != nil {
		utils.Logger.Error("GetByUserNameAndPassword error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	//写入用户ID到cookie中
	utils.GlobalCookie(u.Ctx, "uid", strconv.FormatInt(ad.ID, 10))

	uidByte := []byte(strconv.FormatInt(ad.ID, 10))
	uidString, err := utils.EnPwdCode(uidByte)
	if err != nil {
		utils.Logger.Error("EnPwdCode error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}
	//写入加密后的用户id到cookie
	utils.GlobalCookie(u.Ctx, "sign", uidString)

	//登录成功
	return utils.NewJSONResponse(
		models.ErrorCode.SUCCESS,
		"用户登录成功",
		nil,
	)
}

/**
 * 用户注册功能
 * 接口：/user/register
 * 方法: post
 */
func (u *UserController) PostRegister() mvc.Result {
	u.Ctx.Application().Logger().Info(" User register ")
	utils.Logger.Info("用户注册")
	var (
		User models.User
	)

	err := u.Ctx.ReadForm(&User)
	if err != nil {
		utils.Logger.Error("ReadForm failed", zap.Any("error", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm err", nil)
	}

	// 写入数据库
	err = u.Service.AddUser(&User)
	if err != nil {
		utils.Logger.Error("InsertUser error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "添加用户用户失败", nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "用户注册成功", nil)
}

/**
 * 删除商铺
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
 * 修改商铺
 * 接口：/user/update
 * 方法：post
 */
func (u *UserController) PostUpdate() {
	u.Ctx.Application().Logger().Info(" update User ")
	utils.Logger.Info("更新商铺")
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
 * 查询商铺信息
 * 接口：/user/select
 * 方法：post
 */
func (u *UserController) PostSelect() mvc.Result {
	u.Ctx.Application().Logger().Info(" search User start")
	utils.Logger.Info("开始查询商铺")

	var (
		User models.User
	)

	err := u.Ctx.ReadForm(&User)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm error", nil)
	}

	//从cookie中获取具体商铺
	uid := u.Ctx.GetCookie("uid")
	if uid != "1" {
		User.ID, _ = strconv.ParseInt(uid, 10, 64)
	}

	UserListandCount, err := u.Service.SelectUser(&User)
	if err != nil {
		utils.Logger.Error("SelectUser error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", UserListandCount)
}
