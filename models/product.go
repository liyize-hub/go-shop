package models

import "time"

// 商品结构体
type Product struct {
	ID          int64      `json:"id" form:"id" xorm:"pk autoincr"`
	Name        string     `json:"name" form:"name" xorm:"varchar(32)"`
	Num         int        `json:"num" form:"num"`
	ActivityNum int        `json:"activity_num" form:"activity_num" xorm:"default 0"` //参与秒杀活动的商品数量
	Img         string     `json:"img" form:"img"`
	Last        int        `json:"last" form:"last"` //秒杀活动持续时间
	Price       float64    `json:"price" form:"price"`
	LowPrice    float64    `json:"low_price" form:"low_price"` //参与秒杀活动的商品优惠价格
	ShopID      int64      `json:"shop_id" form:"shop_id" xorm:"index"`
	Shop        *Admin     `xorm:"-"`                              //对应的商铺
	CategoryID  int64      `json:"category_id" form:"category_id"` //商品类型
	Category    *Category  `xorm:"-"`
	CreateTime  time.Time  `json:"create_time" form:"create_time" xorm:"created"`
	Status      int        `json:"status" form:"status" xorm:"tinyint(1) default 0"` //是否参加了秒杀活动 0：未参加，1:参加了 请求：1：添加秒杀活动 2：删除秒杀活动
	Detail      string     `json:"detail" form:"detail"`                             //商品简介
	Flag        int        `json:"flag" form:"flag" xorm:"tinyint(1) default 0"`     //是否被删除的标志字段 软删除 0为有效，1为被删除
	Pages       `xorm:"-"` //分页情况
}
