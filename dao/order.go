package dao

import (
	"errors"
	"go-shop/datasource"
	"go-shop/models"
	"go-shop/utils"
	"time"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

type OrderDao struct {
	*xorm.Engine
}

func NewOrderDao(db *xorm.Engine) *OrderDao {
	//判断数据库连接是否存在
	if db == nil {
		db, err := datasource.NewMysqlConn()
		if err != nil {
			utils.Logger.Info("订单数据重新建立数据库连接失败", zap.Any("error", err))
		}
		return &OrderDao{db}
	}

	return &OrderDao{db}
}

func (d *OrderDao) AddOrder(order *models.Order) (err error) {
	count, err := d.Insert(order)
	if err != nil {
		utils.Logger.Error("插入订单失败", zap.Any("Order", order))
		return
	}
	utils.SugarLogger.Infof("订单成功插入%d条数据,数据id为%d", count, order.ID)
	return
}

func (d *OrderDao) DeleteOrderByID(orderID int64) (bool, error) {
	count, err := d.ID(orderID).Update(&models.Order{Flag: 3})
	if err != nil {
		utils.Logger.Error("删除订单失败", zap.Int64("delete id", orderID))
		return false, err
	}

	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (d *OrderDao) UpdateOrderByID(OrderID int64, Order *models.Order) (err error) {
	count, err := d.ID(OrderID).MustCols("flag").Update(Order)
	if err != nil {
		utils.Logger.Error("更新订单失败", zap.Int64("OrderID", OrderID), zap.Any("Order", Order))
		return
	}
	utils.SugarLogger.Infof("订单成功更新%d条数据,数据id为%d", count, OrderID)

	return
}

// 获取多个数据
func (d *OrderDao) GetOrders(order *models.Order) (*utils.ListAndCount, error) {
	orders := []*models.Order{}

	sess := d.MustCols("flag").Asc("id")

	//如果Page项不为空
	if order.Size != 0 && order.No != 0 {
		sess = sess.Limit(order.Size, (order.No-1)*order.Size)
	}

	//时间范围
	if len(order.TimeRange) != 0 {
		first := time.Unix(order.TimeRange[0], 0).Format("2006-01-02 15:04:05")
		last := time.Unix(order.TimeRange[1], 0).Format("2006-01-02 15:04:05")
		sess = sess.Where("create_time between ? and ?", first, last)
	}

	err := sess.Find(&orders, order) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询订单失败", zap.Any("Order", order))
		return nil, err
	}
	if len(orders) == 0 {
		utils.Logger.Info("订单没有查到相关数据", zap.Any("Order", order))
		return nil, errors.New("订单没有查到相关数据")
	}

	//搜索总数
	count, _ := d.MustCols("flag").Count(order)

	return utils.Lists(orders, uint64(count)), nil
}

//获取单个数据
func (d *OrderDao) GetOrder(order *models.Order) (*models.Order, error) {
	exist, err := d.MustCols("flag").Get(order)
	if err != nil {
		utils.Logger.Error("查询订单失败", zap.Any("Order", order))
		return nil, err
	}
	if !exist {
		utils.Logger.Info("订单没有查到相关数据", zap.Any("Order", order))
		return nil, errors.New("订单没有查到相关数据")
	}

	return order, nil
}
