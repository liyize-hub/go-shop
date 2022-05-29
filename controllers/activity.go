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

type ActivityController struct {
	Ctx     iris.Context             // iris框架自动为每个请求都绑定上下文对象
	Service services.ActivityService // Activity功能实体
}

/**
 * 添加秒杀活动
 * 接口：/activity/add
 * 方法：post
 */
func (p *ActivityController) PostAdd() {
	p.Ctx.Application().Logger().Info(" add Activity start")
	utils.Logger.Info("开始添加秒杀活动")

	var activity models.Activity
	err := p.Ctx.ReadForm(&activity)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
	}
	//从cookie中获取具体商铺
	uid := p.Ctx.GetCookie("uid")
	if uid != "1" {
		activity.ShopID, _ = strconv.ParseInt(uid, 10, 64)
	}

	err = p.Service.InsertActivity(&activity)
	if err != nil {
		utils.Logger.Error("InsertActivity error", zap.Any("err", err))
	}

	utils.Logger.Info("添加秒杀活动成功", zap.Any("Activity", activity))
}

/**
 * 删除秒杀活动
 * 接口：/activity/manager
 * 方法：post
 */
func (p *ActivityController) GetDelete() {
	p.Ctx.Application().Logger().Info(" delete Activity start")
	utils.Logger.Info("开始删除秒杀活动")

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

	err = p.Service.DeleteActivity(id)
	if err != nil {
		utils.Logger.Error("DeleteActivity error", zap.Any("err", err))
	}
}

/**
 * 修改秒杀活动
 * 接口：/activity/update
 * 方法：post
 */
func (p *ActivityController) PostUpdate() {
	p.Ctx.Application().Logger().Info(" update Activity ")
	utils.Logger.Info("更新秒杀活动")
	var activity models.Activity

	err := p.Ctx.ReadForm(&activity)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
	}

	if activity.ID <= 0 {
		utils.Logger.Info("更新数据不正确,id为0", zap.Any("Activity", activity))
		return
	}

	err = p.Service.UpdateActivity(&activity)
	if err != nil {
		utils.Logger.Error("UpdateActivity error", zap.Any("err", err))
	}
	utils.Logger.Info("更新成功", zap.Any("Activity", activity))
}

/**
 * 获取全部秒杀活动
 * 接口：/activity/all
 * 方法：get
 */
func (p *ActivityController) GetAll() mvc.Result {
	p.Ctx.Application().Logger().Info(" get all Activity start")
	utils.Logger.Info("开始获取全部秒杀活动")

	activityListandCount, err := p.Service.SelectAllActivity()
	if err != nil {
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "SelectAllActivity err", nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", activityListandCount)
}

/**
 * 查询秒杀活动
 * 接口：/activity/select
 * 方法：post
 */
func (p *ActivityController) PostSelect() mvc.Result {
	p.Ctx.Application().Logger().Info(" search Activity start")
	utils.Logger.Info("开始查询秒杀活动")
	var (
		activity models.Activity
	)

	err := p.Ctx.ReadForm(&activity)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm error", nil)
	}

	//从cookie中获取具体商铺
	uid := p.Ctx.GetCookie("uid")
	if uid != "1" {
		activity.ShopID, _ = strconv.ParseInt(uid, 10, 64)
	}

	activityListandCount, err := p.Service.SelectActivitys(&activity)
	if err != nil {
		utils.Logger.Error("SelectActivity error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", activityListandCount)
}
