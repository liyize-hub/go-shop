package models

/**
 * 地址结构实体
 */
type Address struct {
	ID            int64  `xorm:"pk autoincr" json:"id"`
	Address       string `json:"address"`                   //地址
	Phone         string `json:"phone"`                     //联系人手机号
	AddressDetail string `json:"address_detail"`            //地址详情
	DelFlag       int    `json:"del_flag" xorm:"default 0"` //是否被删除的标志字段 软删除
}
