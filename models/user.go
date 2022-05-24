package models

import "time"

// 用户结构体
type User struct {
	ID            int64     `json:"id" xorm:"pk autoincr"`
	Name          string    `json:"name" xorm:"varchar(32)"`
	Pwd           string    `json:"pwd" `                       //加密后的密码
	Address       string    `json:"address" xorm:"varchar(32)"` //地址
	AddressDetail string    `json:"address_detail"`             //地址详情
	Phone         string    `json:"phone"`                      //联系人手机号
	CreateTime    time.Time `json:"create_time" xorm:"created"`
	Img           string    `json:"img" xorm:"varchar(255)"`                      //用户的头像
	Flag          int       `json:"flag" form:"flag" xorm:"tinyint(1) default 0"` //用户状态 0为有效，1为删除， 2为被封禁
}
