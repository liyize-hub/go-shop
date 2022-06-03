package models

// 商品种类结构体
type Category struct {
	ID   int64  `json:"id" form:"id" xorm:"pk autoincr"` //主键 自增
	Name string `json:"name" form:"name" xorm:"varchar(32)"`
	Flag int    `json:"flag" form:"flag" xorm:"tinyint(1) default 0"` //是否被删除的标志字段 软删除 0为有效，1为删除
}
