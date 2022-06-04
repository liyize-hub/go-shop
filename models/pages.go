package models

// amis分页处理

type Pages struct {
	No        int     `form:"page,default=1"`                        // 页码
	Size      int     `form:"perPage,default=10"`                    // 分页大小
	TimeRange []int64 `json:"time_range" form:"time_range" xorm:"-"` //时间范围
}
