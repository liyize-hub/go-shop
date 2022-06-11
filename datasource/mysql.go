package datasource

import (
	"go-shop/config"
	"go-shop/models"
	"go-shop/utils"

	_ "github.com/go-sql-driver/mysql" //不能忘记导入
	"go.uber.org/zap"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var DB *xorm.Engine

// 数据库连接
func NewMysqlConn() {
	//如果数据库已连接
	if DB != nil {
		return
	}

	//1.创建数据库引擎对象
	engine, err := xorm.NewEngine(config.DataBaseConfig.Drive, config.DataBaseConfig.URL)
	if err != nil {
		utils.Logger.Error("创建数据库引擎对象失败", zap.Any("error", err))
		return
	}

	//2.数据库引擎关闭
	//defer engine.Close()

	engine.ShowSQL(config.DataBaseConfig.SQLLog)               //设置显示SQL语句
	engine.Logger().SetLevel(log.LOG_DEBUG)                    //设置日志级别
	engine.SetMaxOpenConns(config.DataBaseConfig.MaxOpenConns) //设置最大连接数
	engine.SetMapper(core.GonicMapper{})                       //设置名称映射规则
	err = engine.Sync2(
		new(models.Product),
		new(models.Category),
		new(models.Admin),
		new(models.User),
		new(models.Order),
	) //将自定义的结构体同步到数据库中
	if err != nil {
		utils.Logger.Error("结构体同步数据库失败", zap.Any("error", err))
		return
	}

	DB = engine
}

func init() {
	NewMysqlConn()
}
