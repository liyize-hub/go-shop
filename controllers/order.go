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

type OrderController struct {
	Ctx     iris.Context            // iris框架自动为每个请求都绑定上下文对象
	Service services.OrderService // Order功能实体
}

/**
 * 添加订单
 * 接口：/Order/add
 * 方法：post
 */
func (p *OrderController) PostAdd() {
	p.Ctx.Application().Logger().Info(" add Order start")
	utils.Logger.Info("开始插入订单")

	var Order models.Order
	err := p.Ctx.ReadForm(&Order)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return
	}

	//从cookie中获取具体商铺
	uid := p.Ctx.GetCookie("uid")
	if uid != "1" {
		Order.ShopID, _ = strconv.ParseInt(uid, 10, 64)
	}

	err = p.Service.InsertOrder(&Order)
	if err != nil {
		utils.Logger.Error("InsertOrder error", zap.Any("err", err))
		return
	}

	utils.Logger.Info("插入订单成功", zap.Any("Order", Order))
	//p.Ctx.Redirect("/Order/all")
}

/**
 * 删除订单
 * 接口：/Order/manager
 * 方法：post
 */
func (p *OrderController) GetDelete() {
	p.Ctx.Application().Logger().Info(" delete Order start")
	utils.Logger.Info("开始删除订单")

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

	err = p.Service.DeleteOrder(id)
	if err != nil {
		utils.Logger.Error("DeleteOrder error", zap.Any("err", err))
	}
}

/**
 * 修改订单
 * 接口：/Order/update
 * 方法：post
 */
func (p *OrderController) PostUpdate() {
	p.Ctx.Application().Logger().Info(" update Order ")
	utils.Logger.Info("更新订单")
	var Order models.Order

	err := p.Ctx.ReadForm(&Order)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
	}

	if Order.ID <= 0 {
		utils.Logger.Info("更新数据不正确,id为0", zap.Any("Order", Order))
		return
	}

	err = p.Service.UpdateOrder(&Order)
	if err != nil {
		utils.Logger.Error("UpdateOrder error", zap.Any("err", err))
	}
	utils.Logger.Info("更新成功", zap.Any("Order", Order))
	//p.Ctx.Redirect("/Order/all") //跳转页面
}

/**
 * 获取全部订单
 * 接口：/Order/all
 * 方法：get
 */
func (p *OrderController) GetAll() mvc.Result {
	p.Ctx.Application().Logger().Info(" get all Order start")
	utils.Logger.Info("开始获取全部订单")

	OrderListandCount, err := p.Service.SelectAllOrder()
	if err != nil {
		utils.Logger.Error("SelectAllOrder error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "SelectAllOrder err", nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", OrderListandCount)
}

/**
 * 查询订单
 * 接口：/Order/select
 * 方法：post
 */
func (p *OrderController) PostSelect() mvc.Result {
	p.Ctx.Application().Logger().Info(" search Order start")
	utils.Logger.Info("开始查询订单")
	var (
		Order models.Order
	)

	err := p.Ctx.ReadForm(&Order)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm error", nil)
	}

	//从cookie中获取具体商铺
	uid := p.Ctx.GetCookie("uid")
	if uid != "1" {
		Order.ShopID, _ = strconv.ParseInt(uid, 10, 64)
	}

	OrderListandCount, err := p.Service.SelectOrders(&Order)
	if err != nil {
		utils.Logger.Error("SelectOrder error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", OrderListandCount)
}
