package main

import (
	"go-shop/config"
	"go-shop/router"
	"go-shop/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	app := newApp()

	//应用App设置
	configation(app)

	//5. 注册控制器
	router.Router(app)

	//6. 启动服务
	app.Run(
		iris.Addr(":"+config.ServerConfig.Port),
		iris.WithoutServerError(iris.ErrServerClosed), //忽略iris框架错误
		iris.WithOptimizations,                        //优化
	)

	defer utils.Logger.Sync() //刷新所有缓冲日志条目
}

//构建App
func newApp() *iris.Application {
	//1. 创建iris 实例
	app := iris.New()
	utils.Logger.Info("创建iris 实例成功")
	iris.New().Logger().Info("创建iris 实例成功")

	//2. 设置日志级别，开发阶段为debug
	app.Logger().SetLevel(config.ServerConfig.LogLevel)
	//app.Logger().SetOutput(io.MultiWriter(utils.InfoWriter))

	app.Use(recover.New()) //can recover from any http-relative panics
	app.Use(logger.New())  //log the requests to the terminal. 	eg. app.Logger().Info(ctx.Path())

	//3. 注册视图文件
	template := iris.HTML("./background", ".html") //.Layout("login.html").Reload(true)
	app.RegisterView(template)

	//4. 注册静态资源
	app.HandleDir("/static", "./background/static")

	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问页面出错！"))
		ctx.ViewLayout("")
		ctx.View("error.html")
	})

	return app
}

// 项目设置
func configation(app *iris.Application) {
	app.Configure(iris.WithConfiguration(iris.TOML("./config/iris.toml")))
	/*/错误配置
	//未发现错误
	app.OnErrorCode(iris.StatusNotFound, func(context iris.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusNotFound,
			"msg":    " not found ",
			"data":   iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(context iris.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusInternalServerError,
			"msg":    " interal error ",
			"data":   iris.Map{},
		})
	})*/
}
