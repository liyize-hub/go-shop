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
	Ctx            iris.Context            // iris框架自动为每个请求都绑定上下文对象
	ProductService services.ProductService // Product功能实体
}

/**
 * 获取全部商品
 * 接口：/Product/all
 * 方法：get
 */
func (p *ProductController) GetAll() mvc.Result {
	p.Ctx.Application().Logger().Info(" get all Product start")
	utils.Logger.Info("开始获取全部商品")

	ProductListandCount, err := p.ProductService.SelectAllProduct()
	if err != nil {
		return utils.NewJSONResponse(p.Ctx, models.ErrorCode.ERROR, "", err)
	}

	return utils.NewJSONResponse(p.Ctx,
		models.ErrorCode.SUCCESS,
		"查询成功",
		ProductListandCount,
	)
}

/**
 * 查询商品
 * 接口：/Product/select
 * 方法：post
 */
func (p *ProductController) PostSelect() mvc.Result {
	p.Ctx.Application().Logger().Info(" search Product start")
	utils.Logger.Info("开始查询商品")

	var product models.Product

	err := p.Ctx.ReadJSON(&product)
	if err != nil {
		utils.Logger.Error("ReadJSON error", zap.Any("err", err))
		return utils.NewJSONResponse(p.Ctx, models.ErrorCode.ERROR, "", err)
	}

	fmt.Printf("%#v", product)
	productListandCount, err := p.ProductService.SelectProducts(&product)
	if err != nil {
		utils.Logger.Error("SelectProduct error", zap.Any("err", err))
		return utils.NewJSONResponse(p.Ctx, models.ErrorCode.ERROR, err.Error(), "")
	}

	return utils.NewJSONResponse(p.Ctx,
		models.ErrorCode.SUCCESS,
		"查询成功",
		productListandCount,
	)
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

	err = p.ProductService.DeleteProduct(id)
	if err != nil {
		utils.Logger.Error("DeleteProduct error", zap.Any("err", err))
	}
	//p.Ctx.Redirect("/Product/all")
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

	err := p.Ctx.ReadJSON(&product)
	if err != nil {
		utils.Logger.Error("ReadJSON error", zap.Any("err", err))
	}

	err = p.ProductService.InsertProduct(&product)
	if err != nil {
		utils.Logger.Error("InsertProduct error", zap.Any("err", err))
	}

	utils.Logger.Info("插入商品成功", zap.Any("Product", product))
	//p.Ctx.Redirect("/Product/all")
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

	err := p.Ctx.ReadJSON(&product)
	if err != nil {
		utils.Logger.Error("ReadJSON error", zap.Any("err", err))
	}

	if product.ID <= 0 {
		utils.Logger.Info("更新数据不正确,id为0", zap.Any("Product", product))
		return
	}

	err = p.ProductService.UpdateProduct(&product)
	if err != nil {
		utils.Logger.Error("UpdateProduct error", zap.Any("err", err))
	}
	utils.Logger.Info("更新成功", zap.Any("product", product))
	//p.Ctx.Redirect("/Product/all") //跳转页面
}

/**
 * 展示商品到前台界面
 * 接口：/Product/detail
 * 方法：get
 */
func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.SelectProductByID(1)
	if err != nil {
		utils.Logger.Error("GetDetail error", zap.Any("err", err))
		return mvc.View{}
	}

	return mvc.View{
		Layout: "fronted/web/product/productLayout.html",
		Name:   "fronted/web/product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}
