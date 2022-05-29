package models

import "time"

//管理员结构体
type Admin struct {
	ID          int64      `json:"id" form:"id" xorm:"pk autoincr"` //主键 自增
	Name        string     `json:"name" form:"name" xorm:"varchar(32)"`
	Pwd         string     `json:"pwd" form:"pwd" xorm:"varchar(255)"`                  //管理员密码
	Phone       string     `json:"phone" form:"phone"`                                  //管理员手机号
	ShopImg     string     `json:"shop_img" form:"shop_img" xorm:"varchar(255)"`        //商铺展示图片
	ShopAddress string     `json:"shop_address" form:"shop_address" xorm:"varchar(32)"` //商铺地址
	ShopName    string     `json:"shop_name" form:"shop_name" xorm:"varchar(32)"`       //商铺名称
	ShopDetail  string     `json:"shop_detail" form:"shop_detail" xorm:"varchar(255)"`  //商铺简介
	ShopType    int        `json:"shop_type" form:"shop_type"`                          //商铺经营类型
	CreateTime  time.Time  `json:"create_time" form:"create_time" xorm:"created"`
	TimeRange   []int64    `json:"time_range" form:"time_range" xorm:"-"`        //时间范围
	Flag        int        `json:"flag" form:"flag" xorm:"tinyint(1) default 0"` //管理员状态 0为未激活，1为有效 2为被封禁 3 为被删除
	Pages       `xorm:"-"` //分页情况
}
