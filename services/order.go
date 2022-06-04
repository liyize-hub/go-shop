package services

import (
	"errors"
	"go-shop/dao"
	"go-shop/models"
	"go-shop/utils"

	"xorm.io/xorm"
)

type OrderService interface {
	InsertOrder(*models.Order) error                         //添加订单
	DeleteOrder(int64) error                                 //删除订单
	UpdateOrder(*models.Order) error                         //更新订单
	SelectAllOrder() (*utils.ListAndCount, error)            //查询所有订单
	SelectOrders(*models.Order) (*utils.ListAndCount, error) //查询多条订单数据
	SelectOrderByID(int64) (*models.Order, error)            //查询一条订单数据
}

// 订单服务实现结构体
type orderService struct {
	db *xorm.Engine
}

// 初始化函数
func NewOrderService(db *xorm.Engine) OrderService {
	return &orderService{db: db}
}

func (o *orderService) InsertOrder(order *models.Order) (err error) {
	if order.Num == 0 || order.TotalPrice == 0 {
		return errors.New("插入的订单数量或者价格为0")
	}

	err = dao.NewOrderDao(o.db).AddOrder(order)
	return err
}

func (o *orderService) DeleteOrder(orderID int64) (err error) {
	isOk, err := dao.NewOrderDao(o.db).DeleteOrderByID(orderID)
	if err != nil {
		return
	}

	if isOk {
		utils.SugarLogger.Infof("删除订单成功,ID为:%d", orderID)
	} else {
		utils.SugarLogger.Infof("删除订单失败,ID为:%d", orderID)
	}

	return
}

func (o *orderService) UpdateOrder(order *models.Order) (err error) {
	id := order.ID
	order.ID = 0 //清空主键
	err = dao.NewOrderDao(o.db).UpdateOrderByID(id, order)
	return
}

func (o *orderService) SelectAllOrder() (*utils.ListAndCount, error) {
	orderListandCount, err := dao.NewOrderDao(o.db).GetOrders(&models.Order{})
	return orderListandCount, err
}

func (o *orderService) SelectOrders(order *models.Order) (*utils.ListAndCount, error) {
	orderListandCount, err := dao.NewOrderDao(o.db).GetOrders(order)
	return orderListandCount, err
}

func (o *orderService) SelectOrderByID(orderID int64) (*models.Order, error) {
	Order, err := dao.NewOrderDao(o.db).GetOrder(&models.Order{ID: orderID})
	return Order, err
}
