package models

import "time"

// 用户结构体
type User struct {
	ID         int64      `json:"uid" form:"uid" xorm:"pk autoincr"`
	OpenID     string     `json:"open_id" xorm:"index"`
	Name       string     `json:"name" xorm:"varchar(32)"`
	Address    string     `json:"address" xorm:"varchar(32)"` //地址
	Phone      string     `json:"phone"`                      //联系人手机号
	CreateTime time.Time  `json:"create_time" xorm:"created"`
	TimeRange  []int64    `json:"time_range" form:"time_range" xorm:"-"`        //时间范围
	Img        string     `json:"img" xorm:"varchar(255)"`                      //用户的头像
	Flag       int        `json:"flag" form:"flag" xorm:"tinyint(1) default 0"` //用户状态 0为有效，1为删除， 2为被封禁
	Pages      `xorm:"-"` //分页情况
	Token      string     `json:"token" form:"token" xorm:"-"` //存放用户的临时token
}
