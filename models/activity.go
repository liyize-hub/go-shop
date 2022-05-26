package models

import "time"

// activity 秒杀活动结构体
type Activity struct {
	ID         int64      `json:"id" form:"id" xorm:"pk autoincr"` //主键 自增
	ProductID  int64      `json:"product_id" form:"product_id" xorm:"index"`
	ShopID     int64      `json:"shop_id" form:"shop_id" xorm:"index"`
	Price      float64    `json:"product_price" form:"product_price"` //参与秒杀活动的商品数量
	Num        int        `json:"product_num" form:"product_num"`     //参与秒杀活动的商品价格
	Name       string     `json:"name" form:"name" xorm:"varchar(32)"`//参与秒杀活动的商品名称
	Product    *Product   `xorm:"-"`                                  //对应的秒杀商品
	CreateTime time.Time  `json:"create_time" form:"create_time" xorm:"created"`
	Last       int        `json:"last" form:"last"`                             //秒杀活动持续时间
	Flag       int        `json:"flag" form:"flag" form:"flag" xorm:"tinyint(1) default 0"` // 0为未激活，1已失效, 2为秒杀活动中
	Pages      `xorm:"-"` //分页情况
}
