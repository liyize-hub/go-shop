package router

import (
	"go-shop/config"
	"go-shop/controllers"
	"go-shop/datasource"
	"go-shop/middleware"
	"go-shop/services"
	"go-shop/utils"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"go.uber.org/zap"
)

func Router(app *iris.Application) {
	//连接数据库
	db, err := datasource.NewMysqlConn()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}

	sess := sessions.New(sessions.Config{
		Cookie:  config.ServerConfig.SessionID,
		Expires: 1 * time.Hour, //1小时后过期
	})

	redis := datasource.NewRedis()
	sess.UseDatabase(redis)

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

	// 商品管理模块功能
	product := mvc.New(app.Party("/product"))
	product.Register(
		services.NewProductService(db),
	)
	product.Handle(new(controllers.ProductController))

	// 秒杀活动管理模块功能
	activity := mvc.New(app.Party("/activity"))
	activity.Register(
		services.NewActivityService(db),
	)
	activity.Handle(new(controllers.ActivityController))

	// 订单管理模块功能
	order := mvc.New(app.Party("/order"))
	order.Register(
		services.NewOrderService(db),
	)
	order.Handle(new(controllers.OrderController))

	// 管理员管理模块功能
	admin := mvc.New(app.Party("/admin"))
	admin.Register(
		services.NewAdminService(db),
		sess.Start,
	)
	admin.Handle(new(controllers.AdminController))

	// 用户管理模块功能
	user := mvc.New(app.Party("/user"))
	user.Register(
		services.NewUserService(db),
		sess.Start,
	)
	user.Handle(new(controllers.UserController))
}
