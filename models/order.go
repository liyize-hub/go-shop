package models

import "time"

// 订单结构体
type Order struct {
	ID          int64      `json:"id" xorm:"pk autoincr"`
	UserId      int64      `json:"user_id" xorm:"index"` //用户编号Id
	User        *User      `xorm:"-"`                    //订单对应的账户，并不进行结构体字段映射
	ProductId   int64      `json:"product_id" xorm:"index"`
	Product     []*Product `xorm:"-"` //商品结构体，不进行映射
	OrderStatus int        `json:"order_status"`
	AddressId   int64      `json:"address_id" xorm:"index"` //地址结构体的Id
	Address     *Address   `xorm:"-"`                       //地址结构体，不进行映射
	SumMoney    int64      `json:"sum_money" xorm:"default 0"`
	Time        time.Time  `json:"time" xorm:"updated"`       //时间
	DelFlag     int        `json:"del_flag" xorm:"default 0"` //删除标志 0为正常 1为已删除
}

const (
	OrderWait    = iota
	OrderSuccess //1
	OrderFailed  //2
)
