package models

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 商品结构体
type Product struct {
	ID          int64      `json:"id" form:"id" xorm:"pk autoincr"`
	Name        string     `json:"name" form:"name" xorm:"varchar(32)"`
	Num         int        `json:"num" form:"num"`
	Img         string     `json:"img" form:"img"`
	Price       float64    `json:"price" form:"price"`
	ShopID      int64      `json:"shop_id" form:"shop_id" xorm:"index"`
	Shop        *Admin     `xorm:"-"`                              //对应的商铺
	CategoryID  int64      `json:"category_id" form:"category_id"` //商品类型
	Category    *Category  `xorm:"-"`
	CreateTime  time.Time  `json:"create_time" form:"create_time" xorm:"created"`
	Detail      string     `json:"detail" form:"detail"`                              //商品简介
	Flag        int        `json:"flag" form:"flag" xorm:"tinyint(1) default 0"`      //是否被删除的标志字段 软删除 0为有效，1为被删除
	ActivityNum int        `json:"activity_num" form:"activity_num" xorm:"default 0"` //参与秒杀活动的商品数量exi
	Last        int        `json:"last" form:"last" xorm:"default 0"`                 //秒杀活动持续时间
	LowPrice    float64    `json:"low_price" form:"low_price"`                        //参与秒杀活动的商品优惠价格
	Status      int        `json:"status" form:"status" xorm:"tinyint(1) default 0"`  //是否参加了秒杀活动 0：未参加，1:参加了 请求：1：添加秒杀活动 2：删除秒杀活动
	Pages       `xorm:"-"` //分页情况
}

/**
 * 从Product数据库实体转换为前端请求的resp的json格式
 */
func (this *Product) ProductToRespDesc(rdb *redis.Client) interface{} {
	respDesc := map[string]interface{}{
		"id":      this.ID,
		"shop_id": this.ShopID,
		"name":    this.Name,
		"price":   this.Price,
		"detail":  this.Detail,
		"num":     this.Num,
		"img":     this.Img,
	}
	rdb.HMSet(this.Name+":pro:"+strconv.FormatInt(this.ID, 10), respDesc)
	return respDesc
}
