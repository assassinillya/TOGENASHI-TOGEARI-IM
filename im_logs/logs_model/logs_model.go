package logs_model

import "im_server/common/models"

type LogModel struct {
	models.Model
	LogType      int8   `json:"logType"` // 日志类型  2 操作日志 3 运行日志
	IP           string `json:"ip"`
	Addr         string `json:"addr"`
	UserID       uint   `json:"userID"`
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
	Level        string `json:"level"`
	Title        string `json:"title"`
	Content      string `json:"content"` // 日志详情
	Service      string `json:"service"` // 服务  记录微服务的名称
}
