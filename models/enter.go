package models

import "time"

type MODEL struct {
	ID        uint      `json:"id,select($any)" gorm:"primarykey" structs:"-"` // 主键ID
	CreatedAt time.Time `json:"created_at,select($any)" structs:"-"`           // 创建时间
	UpdatedAt time.Time `json:"-" structs:"-"`                                 // 更新时间
}

// PageInfo 图片展示页面（有分页功能）---- 因为和表结构强相关，所以放这里
type PageInfo struct {
	Page  int    `form:"page"`  // 当前页，当前在第几页
	Key   string `form:"key"`   // 模糊匹配的关键字
	Limit int    `form:"limit"` // 限制数量,如果不设置，默认为0，就是不分页
	Sort  string `form:"sort"`  // 排序,谁在前面谁在后面的排序
}

// RemoveRequest 批量删除图片，通过要删除的图片的id列表来实现批量删除
type RemoveRequest struct {
	IDList []uint `json:"id_list"`
}

// ESIDRequest 这个是存在es中的文章id
type ESIDRequest struct {
	ID string `json:"id" form:"id" uri:"id"`
}

// ESIDListRequest 这个是存在es中的文章id列表，用于批量删除文章
type ESIDListRequest struct {
	IDList []string `json:"id_list" binding:"required"`
}
