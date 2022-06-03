package models

import "time"

// 订单结构体
type Order struct {
	ID         int64     `json:"id" xorm:"pk autoincr"`
	UserID     int64     `json:"user_id" xorm:"index"` //用户编号Id
	User       *User     `xorm:"-"`                    //订单对应的账户，并不进行结构体字段映射
	ProductID  int64     `json:"product_id" xorm:"index"`
	Product    *Product  `xorm:"-"`
	CreateTime time.Time `json:"time" xorm:"created"`                          //订单创建时间
	Flag       int       `json:"flag" form:"flag" xorm:"tinyint(1) default 0"` //是否被删除的标志字段 软删除 0为有效，1为删除
}

const (
	OrderWait    = iota //待支付
	OrderPay            //已支付
	OrderSend           //已发货
	OrderSuccess        //已完成
)
