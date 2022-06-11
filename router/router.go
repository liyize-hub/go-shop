package router

import (
	"go-shop/controllers"
	"go-shop/datasource"
	"go-shop/middleware"
	"go-shop/services"
	"go-shop/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
)

func Router(app *iris.Application) {
	//管理员登录页面
	app.Get("/", func(ctx iris.Context) {
		ctx.View("login.html")
	})

	//后台管理页面
	app.Get("/background", func(ctx iris.Context) {
		uid := ctx.GetCookie("uid")

		if uid == "1" {
			utils.Logger.Info("超级管理员登录！", zap.String("uid", uid))
			ctx.View("index_root.html") //超级管理员后台管理页面
			return
		} else {
			utils.Logger.Info("商铺管理员登录！", zap.String("uid", uid))
			ctx.View("index.html") //商铺管理员后台管理页面
			return
		}
	}).Use(middleware.AuthConProduct)

	app.Get("/test", func(ctx iris.Context) {
		ctx.WriteString("test")
	})

	// 获取图片并上传到腾讯云COS
	app.Post("/uploadimg", services.GetImg)

	// 商品种类管理模块功能
	app.Get("/category", services.GetAllCategory)

	// 商品秒杀
	app.Get("/seckill",services.Seckill)

	// 商品管理模块功能
	product := mvc.New(app.Party("/product"))
	product.Register(
		services.NewProductService(datasource.DB, datasource.Rdb),
	)
	product.Handle(new(controllers.ProductController))

	// 订单管理模块功能
	order := mvc.New(app.Party("/order"))
	order.Register(
		services.NewOrderService(datasource.DB),
	)
	order.Handle(new(controllers.OrderController))

	// 管理员管理模块功能
	admin := mvc.New(app.Party("/admin"))
	admin.Register(
		services.NewAdminService(datasource.DB),
	)
	admin.Handle(new(controllers.AdminController))

	// 用户管理模块功能
	user := mvc.New(app.Party("/user"))
	user.Register(
		services.NewUserService(datasource.DB, datasource.Rdb),
	)
	user.Handle(new(controllers.UserController))
}
