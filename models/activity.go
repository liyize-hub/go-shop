package models

import "time"

// activity 秒杀活动结构体
type Activity struct {
	ID         int64     `json:"id" xorm:"pk autoincr"` //主键 自增
	ProductID  int64     `json:"product_id" xorm:"index"`
	Price      float64   `json:"price"` //参与秒杀活动的商品数量
	Num        int       `json:"num"`   //参与秒杀活动的商品价格
	Product    *Product  `xorm:"-"`     //对应的秒杀商品
	CreateTime time.Time `json:"create_time" xorm:"created"`
	StartTime  time.Time `json:"start_time"`                                   //秒杀活动开始时间
	EndTime    time.Time `json:"end_time"`                                     //秒杀活动结束时间
	Flag       int       `json:"flag" form:"flag" xorm:"tinyint(1) default 0"` //是否被删除的标志字段 软删除 0为有效，1为删除
}
