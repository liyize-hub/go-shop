package models

import "time"

// 订单结构体
type Order struct {
	ID         int64      `json:"id" xorm:"pk autoincr"`
	UserID     int64      `json:"user_id" xorm:"index"` //用户编号Id
	User       *User      `xorm:"-"`
	Num        int        `json:"num" form:"num"`                 //订单中商品总数
	TotalPrice float64    `json:"total_price" form:"total_price"` //订单中商品总价格
	ProductID  int64      `json:"product_id" xorm:"index"`
	Product    *Product   `xorm:"-"`
	ShopID     int64      `json:"shop_id" form:"shop_id" xorm:"index"`
	Shop       *Admin     `xorm:"-"`                                            //对应的商铺
	CreateTime time.Time  `json:"time" xorm:"created"`                          //订单创建时间
	Flag       int        `json:"flag" form:"flag" xorm:"tinyint(1) default 0"` //订单状态 0:已支付 1:已发货 2:已完成 3:已删除
	Pages      `xorm:"-"` //分页情况
}

const (
	OrderPay     = iota //已支付
	OrderSend           //已发货
	OrderSuccess        //已完成
	OrderDelete         //已删除
)
