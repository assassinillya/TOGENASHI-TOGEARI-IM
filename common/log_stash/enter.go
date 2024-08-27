package log_stash

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type Pusher struct {
	LogType int8   `json:"logType"` // 日志类型 2 操作日志 3 运行日志
	IP      string `json:"ip"`
	UserID  uint   `json:"userID"`
	Level   string `json:"level"`
	Title   string `json:"title"`
	Content string `json:"content"` // 日志详情
	Service string `json:"service"` // 服务 记录微服务的名称
	client  *kq.Pusher
}

func (p *Pusher) Save(ctx context.Context) {
	fmt.Println(p.client)
	userIDs := ctx.Value("UserID")
	var userID uint
	if userIDs != nil {
		userIntID, _ := strconv.Atoi(userIDs.(string))
		userID = uint(userIntID)
	}

	clientIP := ctx.Value("clientIP").(string)
	p.IP = clientIP
	p.UserID = userID

	if p.client == nil {
		return
	}
	byteData, err := json.Marshal(p)
	if err != nil {
		logx.Error(err)
	}
	p.client.Push(context.Background(), string(byteData))
}

// Info 为什么是指针 因为要改值
func (p *Pusher) Info(title string, content string) {
	p.Title = title
	p.Content = content
}

func NewActionPusher(client *kq.Pusher, serviceName string) *Pusher {
	return NewPusher(client, 2, serviceName)

}
func NewRuntimePusher(client *kq.Pusher, serviceName string) *Pusher {
	return NewPusher(client, 3, serviceName)
}

func NewPusher(client *kq.Pusher, LogType int8, serviceName string) *Pusher {
	//userIDs := ctx.Value("UserID")
	//var userID uint
	//if userIDs != nil {
	//	userIntID, _ := strconv.Atoi(userIDs.(string))
	//	userID = uint(userIntID)
	//}
	//
	//clientIP := ctx.Value("clientIP").(string)

	return &Pusher{
		//IP:      clientIP,
		//UserID:  userID,
		LogType: LogType,
		Service: serviceName,
		client:  client,
	}

}
