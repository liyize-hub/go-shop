package router

import (
	"go-shop/controllers"
	"go-shop/datasource"
	"go-shop/services"
	"go-shop/utils"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

func Router(app *iris.Application) {
	//连接数据库
	db, err := datasource.NewMysqlConn()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}

	sess := sessions.New(sessions.Config{
		Cookie:  "Cookie",
		Expires: 1 * time.Hour, //1小时后过期
	})

	redis := datasource.NewRedis()
	sess.UseDatabase(redis)

	// 商品管理模块功能
	product := mvc.New(app.Party("/product"))
	product.Register(
		services.NewProductService(db),
	)
	product.Handle(new(controllers.ProductController))

	// 订单管理模块功能
	order := mvc.New(app.Party("/order"))
	order.Register(
		services.NewOrderService(db),
	)
	order.Handle(new(controllers.OrderController))

	// 管理员模块功能
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
