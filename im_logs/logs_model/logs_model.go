package logs_model

import "im_server/common/models"

type LogModel struct {
	models.Model
	LogType      int8   `json:"logType"` // 日志类型  2 操作日志 3 运行日志
	IP           string `gorm:"size:32" json:"ip"`
	Addr         string `gorm:"size:32" json:"addr"`
	UserID       uint   `json:"userID"`
	UserNickname string `gorm:"size:64" json:"userNickname"`
	UserAvatar   string `gorm:"size:256" json:"userAvatar"`
	Level        string `gorm:"size:12" json:"level"`
	Title        string `gorm:"size:32" json:"title"`
	Content      string `json:"content"`                // 日志详情
	Service      string `gorm:"size:32" json:"service"` // 服务  记录微服务的名称
}
