package dao

import (
	"errors"
	"go-shop/datasource"
	"go-shop/models"
	"go-shop/utils"

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
			utils.Logger.Info("重新建立数据库连接失败", zap.Any("error", err))
		}
		return &OrderDao{db}
	}

	return &OrderDao{db}
}

func (o *OrderDao) AddOrder(order *models.Order) (err error) {
	count, err := o.Insert(order)
	if err != nil {
		utils.Logger.Error("插入失败", zap.Any("order", order))
		return
	}
	utils.SugarLogger.Infof("成功插入%d条数据,数据id为%d", count, order.ID)
	return
}

func (o *OrderDao) DeleteOrderByID(orderID int64) (bool, error) {
	count, err := o.Delete(&models.Order{ID: orderID})
	if err != nil {
		utils.Logger.Error("删除失败", zap.Int64("delete id", orderID))
		return false, err
	}

	if count == 1 {
		return true, nil
	}
	return false, nil
}

func (o *OrderDao) UpdateOrderByID(orderID int64, order *models.Order) (err error) {
	count, err := o.ID(orderID).MustCols("flag").Update(order)
	if err != nil {
		utils.Logger.Error("更新失败", zap.Int64("orderID", orderID), zap.Any("order", order))
		return
	}
	utils.SugarLogger.Infof("成功更新%d条数据,数据id为%d", count, orderID)

	return
}

//获取多个数据
func (o *OrderDao) GetOrders(order *models.Order) (*utils.ListAndCount, error) {
	orders := []*models.Order{}

	err := o.MustCols("flag").Asc("id").Find(&orders, order) // 返回值，条件
	if err != nil {
		utils.Logger.Error("查询失败", zap.Any("order", order))
		return nil, err
	}
	if len(orders) == 0 {
		utils.Logger.Info("没有查到相关数据", zap.Any("order", order))
		return nil, errors.New("没有查到相关数据")
	}

	return utils.Lists(orders, uint64(len(orders))), nil
}

//获取单个数据
func (o *OrderDao) GetOrder(order *models.Order) (*models.Order, error) {
	exist, err := o.Get(order)
	if err != nil {
		utils.Logger.Error("查询失败", zap.Any("order", order))
		return nil, err
	}
	if !exist {
		utils.Logger.Info("没有查到相关数据", zap.Any("order", order))
		return nil, errors.New("没有查到相关数据")
	}

	return order, nil
}

