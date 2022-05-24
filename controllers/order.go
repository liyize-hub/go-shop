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
	Ctx          iris.Context          // iris框架自动为每个请求都绑定上下文对象
	Service services.OrderService // Order功能实体
}

/**
 * 获取全部订单
 * 接口：/Order/all
 * 方法：get
 */
func (o *OrderController) GetAll() mvc.Result {
	o.Ctx.Application().Logger().Info(" get all Order start")
	utils.Logger.Info("开始获取全部订单")

	orderListandCount, err := o.Service.SelectAllOrder()
	if err != nil {
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "SelectAllOrder err", nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", orderListandCount)

}

/**
 * 查询订单
 * 接口：/Order/select
 * 方法：post
 */
func (o *OrderController) PostSelect() mvc.Result {
	o.Ctx.Application().Logger().Info(" search Order start")
	utils.Logger.Info("开始查询订单")

	var (
		order models.Order
	)

	err := o.Ctx.ReadJSON(&order)
	if err != nil {
		utils.Logger.Error("ReadJSON error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadJSON error", nil)
	}

	orderListandCount, err := o.Service.SelectOrders(&order)
	if err != nil {
		utils.Logger.Error("SelectOrder error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", orderListandCount)
}

/**
 * 删除订单
 * 接口：/Order/manager
 * 方法：post
 */
func (o *OrderController) GetDelete() {
	o.Ctx.Application().Logger().Info(" delete Order start")
	utils.Logger.Info("开始删除订单")

	idString := o.Ctx.URLParam("id")
	if idString == "" {
		utils.Logger.Info("传入id为空")
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		utils.Logger.Error("parseInt error", zap.Any("err", err))
		return
	}

	err = o.Service.DeleteOrder(id)
	if err != nil {
		utils.Logger.Error("DeleteOrder error", zap.Any("err", err))
	}
}

/**
 * 添加订单
 * 接口：/Order/add
 * 方法：post
 */
func (o *OrderController) PostAdd() {
	o.Ctx.Application().Logger().Info(" add Order start")
	utils.Logger.Info("开始插入订单")

	var order models.Order

	err := o.Ctx.ReadJSON(&order)
	if err != nil {
		utils.Logger.Error("ReadJSON error", zap.Any("err", err))
	}

	err = o.Service.InsertOrder(&order)
	if err != nil {
		utils.Logger.Error("InsertOrder error", zap.Any("err", err))
	}

	utils.Logger.Info("插入订单成功", zap.Any("order", order))
}

/**
 * 添加订单
 * 接口：/Order/update
 * 方法：post
 */
func (o *OrderController) PostUpdate() {
	o.Ctx.Application().Logger().Info(" update Order ")
	utils.Logger.Info("更新订单")
	var Order models.Order

	err := o.Ctx.ReadJSON(&Order)
	if err != nil {
		utils.Logger.Error("ReadJSON error", zap.Any("err", err))
	}

	if Order.ID <= 0 {
		utils.Logger.Info("更新数据不正确,id为0", zap.Any("Order", Order))
		return
	}

	err = o.Service.UpdateOrder(&Order)
	if err != nil {
		utils.Logger.Error("UpdateOrder error", zap.Any("err", err))
	}
	utils.Logger.Info("更新成功", zap.Any("Order", Order))
}

/*func (p *ProductController) GetOrder() mvc.View {
	productString := p.Ctx.URLParam("ID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.Atoi(productString) //string -> int
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	var orderID int64
	showMessage := "抢购失败！"
	//判断商品数量是否满足需求
	if product.ProductNum > 0 {
		//扣除商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		//创建订单
		userID, err := strconv.Atoi(userString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}

		order := &datamodels.Order{
			UserId:      int64(userID),
			ProductId:   int64(productID),
			OrderStatus: datamodels.OrderSuccess,
		}
		//新建订单
		orderID, err = p.Service.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功！"
		}
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": showMessage,
		},
	}
}*/
