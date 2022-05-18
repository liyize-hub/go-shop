package models

import "time"

type User struct {
	ID         int64     `json:"id" xorm:"pk autoincr"`
	Name       string    `json:"name" xorm:"varchar(32)"`
	NickName   string    `json:"nick_name"` //别名
	Pwd        string    `json:"pwd"`       //加密后的密码
	CreateTime time.Time `json:"create_time" xorm:"created"`
	AddressId  int64     `json:"address_id" xorm:"index"`    //地址结构体的Id
	Address    *Address  `xorm:"-"`                          //地址结构体，不进行映射
	Img        string    `json:"avatar" xorm:"varchar(255)"` //用户的头像
	DelFlag    int       `json:"del_flag" xorm:"default 0"`  //是否被删除的标志字段 软删除
}
