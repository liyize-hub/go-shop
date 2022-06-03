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
	Shop       *Admin     `xorm:"-"`                              //对应的商铺
	CategoryID int64      `json:"category_id" form:"category_id"` //商品类型
	Category   *Category  `xorm:"-"`
	CreateTime time.Time  `json:"create_time" form:"create_time" xorm:"created"`
	TimeRange  []int64    `json:"time_range" form:"time_range" xorm:"-"`        //时间范围
	Detail     string     `json:"detail" form:"detail"`                         //商品简介
	Flag       int        `json:"flag" form:"flag" xorm:"tinyint(1) default 0"` //是否被删除的标志字段 软删除 0为有效，1为被删除
	Pages      `xorm:"-"` //分页情况
}
