// Code generated by goctl. DO NOT EDIT.
package types

type LogInfoResponse struct {
	ID           uint   `json:"id"`
	CreatedAt    string `json:"created_at"`
	LogType      int8   `json:"logType"` // 日志类型 2 操作日志 3 运行日志
	IP           string `json:"ip"`
	Addr         string `json:"addr"`
	UserID       uint   `json:"userId"`
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
	Level        string `json:"level"`
	Title        string `json:"title"`
	Content      string `json:"content"` // 日志详情
	Service      string `json:"service"` // 服务 记录微服务的名称
	IsRead       bool   `json:"isRead"`
}

type LogListRequest struct {
	Page  int `form:"page,optional"`
	Limit int `form:"limit,optional"`
	Key   int `form:"key,optional"`
}

type LogListResponse struct {
	List  []LogInfoResponse `json:"list"`
	Count int               `json:"count"`
}
