package services

import (
	"encoding/json"
	"go-shop/datasource"
	"go-shop/middleware"
	"go-shop/models"
	"go-shop/utils"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

func Seckill(ctx iris.Context) {
	rdb, err := datasource.NewRedisConn()
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	defer rdb.Close()

	token := ctx.URLParam("token")
	productKey := "pro:" + ctx.URLParam("pid")

	if token == "" || productKey == "" {
		utils.Logger.Error("请求数据为空")
		utils.SendJSON(ctx, models.ErrorCode.ERROR, "请求数据为空", nil)
		return
	}

	//查找用户id
	cmd := rdb.Get(token)
	if cmd.Err() != nil {
		utils.Logger.Error("Get Redis error", zap.String("token", token), zap.Any("error", cmd.Err()))
		utils.SendJSON(ctx, models.ErrorCode.ERROR, "用户未登录", nil)
		return
	}

	uid, err := cmd.Int64()
	if err != nil {
		utils.Logger.Error("ParseInt error", zap.Int64("uid", uid), zap.Any("error", err))
		utils.SendJSON(ctx, models.ErrorCode.ERROR, "转换错误", nil)
		return
	}

	//执行抢购脚本
	cmd1 := rdb.EvalSha(datasource.Sha1.Val(), []string{productKey})
	result, err := cmd1.Int()
	if err != nil {
		utils.Logger.Error("EvalSha error", zap.Any("err", err))
		utils.SendJSON(ctx, models.ErrorCode.ERROR, "抢购出错", nil)
		return
	}

	if result != 1 {
		ctx.WriteString("抢购失败！！！")
		return
	}

	//创建消息体
	message := models.NewMessage(uid, productKey)
	//类型转化
	byteMessage, err := json.Marshal(message)
	if err != nil {
		utils.Logger.Error("json.Marshal error", zap.Any("err", err))
		utils.SendJSON(ctx, models.ErrorCode.ERROR, "转换错误", nil)
		return
	}

	//1. 先返回消息
	ctx.WriteString("抢购成功！！！")
	//新建消息对象，队列名为“seckill”
	mq := middleware.NewRabbitMQSimple("seckill")
	//2. 后异步处理传入数据库
	mq.PublishSimple(string(byteMessage))
	return
}
