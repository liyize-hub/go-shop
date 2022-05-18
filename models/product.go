package models

import "time"

// 商品结构体
type Product struct {
	ID            int64     `json:"id" xorm:"pk autoincr"`
	Name          string    `json:"name" xorm:"varchar(32)"`
	Num           int64     `json:"num"`
	Img           string    `json:"img"`
	Url           string    `json:"url"`
	Price         float64	`json:"price"`
	OriginalPrice float64   `json:"original_price"`
	Time          time.Time `json:"time" xorm:"updated"`       //时间
	DelFlag       int       `json:"del_flag" xorm:"default 0"` //是否被删除的标志字段 软删除
}
