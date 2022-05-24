package models

// amis分页处理

type Pages struct {
	No   int `form:"page,default=1"`     // 页码
	Size int `form:"perPage,default=10"` // 分页大小
}
