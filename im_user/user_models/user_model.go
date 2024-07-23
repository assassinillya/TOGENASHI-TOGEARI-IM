package user_models

import "im_server/common/models"

// UserModel 用户表
type UserModel struct {
	models.Model
	//用户号就是models的ID
	Pwd      string `gorm:"size:64" json:"pwd"`
	Nickname string `gorm:"size:32" json:"nickname"`
	Abstract string `gorm:"size:128" json:"abstract"`
	Avatar   string `gorm:"size:256" json:"avatar"`
	IP       string `gorm:"size:32" json:"ip"`
	Addr     string `gorm:"size:64" json:"addr"`
	Role     int8   `json:"role"`                 // 角色 1 管理员 2 普通用户
	OpenID   string `gorm:"size:64" json:"token"` // 第三方平台登录凭证
}
