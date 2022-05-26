package models

import "time"

// 商品结构体
type Product struct {
	ID         int64      `json:"id" form:"id" xorm:"pk autoincr"`
	Name       string     `json:"name" form:"name" xorm:"varchar(32)"`
	Num        int        `json:"num" form:"num"`
	Img        string     `json:"img" form:"img"`
	Price      float64    `json:"price" form:"price"`
	ShopID     int64      `json:"shop_id" form:"shop_id" xorm:"index"`
	Shop       *Admin     `xorm:"-"`                        //对应的商铺
	Category   int        `json:"category" form:"category"` //商品类型
	CreateTime time.Time  `json:"create_time" form:"create_time" xorm:"created"`
	Detail     string     `json:"detail" form:"detail"`                         //商品描述
	Flag       int        `json:"flag" form:"flag" xorm:"tinyint(1) default 0"` //是否被删除的标志字段 软删除 0为有效，1为被删除
	Pages      `xorm:"-"` //分页情况
}
