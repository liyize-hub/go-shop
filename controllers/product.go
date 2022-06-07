package controllers

import (
	"fmt"
	"go-shop/models"
	"go-shop/services"
	"go-shop/utils"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
)

type ProductController struct {
	Ctx     iris.Context            // iris框架自动为每个请求都绑定上下文对象
	Service services.ProductService // Product功能实体
}

/**
 * 添加商品
 * 接口：/Product/add
 * 方法：post
 */
func (p *ProductController) PostAdd() {
	p.Ctx.Application().Logger().Info(" add Product start")
	utils.Logger.Info("开始插入商品")

	var product models.Product
	err := p.Ctx.ReadForm(&product)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return
	}

	//从cookie中获取具体商铺
	uid := p.Ctx.GetCookie("uid")
	if uid != "1" {
		product.ShopID, _ = strconv.ParseInt(uid, 10, 64)
	}

	err = p.Service.InsertProduct(&product)
	if err != nil {
		utils.Logger.Error("InsertProduct error", zap.Any("err", err))
		return
	}

	utils.Logger.Info("插入商品成功", zap.Any("Product", product))
	//p.Ctx.Redirect("/Product/all")
}

/**
 * 删除商品
 * 接口：/Product/manager
 * 方法：post
 */
func (p *ProductController) GetDelete() {
	p.Ctx.Application().Logger().Info(" delete Product start")
	utils.Logger.Info("开始删除商品")

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

	err = p.Service.DeleteProduct(id)
	if err != nil {
		utils.Logger.Error("DeleteProduct error", zap.Any("err", err))
	}
}

/**
 * 修改商品
 * 接口：/Product/update
 * 方法：post
 */
func (p *ProductController) PostUpdate() {
	p.Ctx.Application().Logger().Info(" update Product ")
	utils.Logger.Info("更新商品")
	var product models.Product

	err := p.Ctx.ReadForm(&product)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
	}

	fmt.Println("**********", product.ID)
	if product.ID <= 0 {
		utils.Logger.Info("更新数据不正确,id为0", zap.Any("Product", product))
		return
	}

	err = p.Service.UpdateProduct(&product)
	if err != nil {
		utils.Logger.Error("UpdateProduct error", zap.Any("err", err))
	}
	utils.Logger.Info("更新成功", zap.Any("product", product))
	//p.Ctx.Redirect("/Product/all") //跳转页面
}

/**
 * 获取全部商品
 * 接口：/Product/all
 * 方法：get
 */
func (p *ProductController) GetAll() mvc.Result {
	p.Ctx.Application().Logger().Info(" get all Product start")
	utils.Logger.Info("开始获取全部商品")

	ProductListandCount, err := p.Service.SelectAllProduct()
	if err != nil {
		utils.Logger.Error("SelectAllProduct error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "SelectAllProduct err", nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", ProductListandCount)
}

/**
 * 查询商品
 * 接口：/Product/select
 * 方法：post
 */
func (p *ProductController) PostSelect() mvc.Result {
	p.Ctx.Application().Logger().Info(" search Product start")
	utils.Logger.Info("开始查询商品")
	var (
		product models.Product
	)

	err := p.Ctx.ReadForm(&product)
	if err != nil {
		utils.Logger.Error("ReadForm error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, "ReadForm error", nil)
	}

	//从cookie中获取具体商铺
	uid := p.Ctx.GetCookie("uid")
	if uid != "1" && uid != "" {
		product.ShopID, _ = strconv.ParseInt(uid, 10, 64)
	}

	productListandCount, err := p.Service.SelectProducts(&product)
	if err != nil {
		utils.Logger.Error("SelectProduct error", zap.Any("err", err))
		return utils.NewJSONResponse(models.ErrorCode.ERROR, err.Error(), nil)
	}

	return utils.NewJSONResponse(models.ErrorCode.SUCCESS, "查询成功", productListandCount)
}
