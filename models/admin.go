package models

import "time"

//管理员结构体
type Admin struct {
	ID         int64     `json:"id" xorm:"pk autoincr"` //主键 自增
	Name       string    `json:"name" xorm:"varchar(32)"`
	Pwd        string    `json:"pwd" xorm:"varchar(255)"` //管理员密码
	CreateTime time.Time `json:"create_time" xorm:"created"`
	DelFlag    int     `json:"del_flag" xorm:"default 0"` //是否被删除的标志字段 软删除
}
