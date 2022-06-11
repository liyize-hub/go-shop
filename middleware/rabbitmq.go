package middleware

import (
	"encoding/json"
	"errors"
	"go-shop/datasource"
	"go-shop/models"
	"go-shop/utils"
	"strconv"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

//连接信息
const MQURL = "amqp://root:liyize01@127.0.0.1:5672/test"

//rabbitMQ结构体
type RabbitMQ struct {
	//要保存的连接
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
}

//创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

//断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		utils.Logger.Error(message, zap.Any("error", err))
	}
}

//创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabb"+
		"itmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

//直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		r.failOnErr(err, "创建队列失败")
	}
	//调用channel 发送消息到队列中
	err = r.channel.Publish(
		r.Exchange, //调用默认的default 交换机
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		r.failOnErr(err, "发送消息失败")
	}
}

//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		r.failOnErr(err, "创建队列失败")
	}

	//消费者流控
	r.channel.Qos(
		1,     //当前消费者一次能接受的最大消息数量
		0,     //服务器传递的最大容量（以八位字节为单位）
		false, //如果设置为true 对channel可用
	)

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		false, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞 设置为阻塞：消费完一个下一个再消费
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		r.failOnErr(err, "接收消息失败")
	}

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			msg := models.Message{}
			err = json.Unmarshal(d.Body, &msg)
			if err != nil {
				utils.Logger.Error("json.Unmarshal error", zap.Any("error", err))
			}
			//插入订单
			err = InsertOrderByMessage(msg)
			if err != nil {
				utils.Logger.Error("InsertOrderByMessage error", zap.Any("error", err))
			}
			//为false表示确认当前消息
			d.Ack(false)
		}
	}()
	<-forever
}

func InsertOrderByMessage(msg models.Message) (err error) {
	if msg.ProductKey == "" || msg.UserID == 0 {
		utils.Logger.Error("msg is empty", zap.Any("msg", msg))
		return
	}

	order   := &models.Order{}
	product := &models.Product{}

	product.ID, err = strconv.ParseInt(msg.ProductKey[4:], 10, 64)
	if err != nil {
		utils.Logger.Error("ParseInt error", zap.String("productID", msg.ProductKey[4:]))
		return
	}

	// 查询商品信息,获取shopID
	exist, err := datasource.DB.MustCols("flag").Get(product)
	if err != nil {
		utils.Logger.Error("查询商品失败", zap.Any("product", product))
		return err
	}
	if !exist {
		utils.Logger.Info("商品没有查到相关数据", zap.Any("product", product))
		return errors.New("商品没有查到相关数据")
	}

	order.ProductID = product.ID
	order.ShopID = product.ShopID
	order.UserID = msg.UserID
	order.Num = 1
	order.TotalPrice = product.LowPrice
	_, err = datasource.DB.Insert(order)
	if err != nil {
		utils.Logger.Error("插入订单失败", zap.Any("Order", order))
		return
	}

	return nil
}
