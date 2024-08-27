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
	items   []string
	ctx     context.Context
}

func (p *Pusher) Save(ctx context.Context) {

	if p.ctx == nil {
		// 如果没有ctx则使用老的ctx
		p.ctx = ctx
	}

	if p.client == nil {
		return
	}

	if len(p.items) > 0 {
		for _, item := range p.items {
			p.Content += item
		}
		p.items = []string{}
	}

	userIDs := p.ctx.Value("userID")
	var userID uint
	if userIDs != nil {
		userIntID, _ := strconv.Atoi(userIDs.(string))
		userID = uint(userIntID)
	}

	clientIP := p.ctx.Value("clientIP").(string)
	p.IP = clientIP
	p.UserID = userID
	if p.Level == "" {
		p.Level = "info"
	}

	byteData, err := json.Marshal(p)
	if err != nil {
		logx.Error(err)
	}
	p.client.Push(p.ctx, string(byteData))
}

// SetItem 这个函数是为了兼容之前的版本
func (p *Pusher) SetItem(label string, val any) {
	p.setItem("info", label, val)
}

func (p *Pusher) SetItemInfo(label string, val any) {
	p.setItem("info", label, val)
}

func (p *Pusher) SetItemWarn(label string, val any) {
	p.setItem("warn", label, val)
}

func (p *Pusher) SetItemErr(label string, val any) {
	p.setItem("err", label, val)
}

func (p *Pusher) setItem(level string, label string, val any) {
	var str string
	switch value := val.(type) {
	case string:
		str = fmt.Sprintf("<div class=\"log_item_label\">%s</div> <div class=\"log_item_content\">%s</div>", label, value)
	case int, uint, uint32, uint64, int32, int8:
		str = fmt.Sprintf("<div class=\"log_item_label\">%s</div> <div class=\"log_item_content\">%d</div>", label, value)
	default:
		byteData, _ := json.Marshal(val)
		str = fmt.Sprintf("<div class=\"log_item_label\">%s</div> <div class=\"log_item_content\">%s</div>", label, string(byteData))
	}
	logItem := fmt.Sprintf("<div class=\"log_item %s\">%s</div>", level, str)
	p.items = append(p.items, logItem)
}

// Info 为什么是指针 因为要改值
func (p *Pusher) Info(title string) {
	p.Title = title
	p.Level = "info"
}

func (p *Pusher) Warning(title string) {
	p.Title = title
	p.Level = "warning"
}

func (p *Pusher) Err(title string) {
	p.Title = title
	p.Level = "err"
}

func (p *Pusher) SetCtx(ctx context.Context) {
	p.ctx = ctx
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
