package mqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"im_server/im_logs/logs_api/internal/svc"
	"im_server/im_logs/logs_model"
	"im_server/im_user/user_rpc/types/user_rpc"
	"sync"
)

type LogEvent struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *LogEvent {
	return &LogEvent{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Request struct {
	LogType int8   `json:"log_type"` // 日志类型 2 操作日志 3 运行日志
	IP      string `json:"ip"`
	UserID  uint   `json:"user_id"`
	Level   string `json:"level"`
	Title   string `json:"title"`
	Content string `json:"content"` // 日志详情
	Service string `json:"service"` // 服务 记录微服务的名称
}

func (l *LogEvent) Consume(ctx context.Context, key, val string) error {
	var req Request
	err := json.Unmarshal([]byte(val), &req)
	if err != nil {
		logx.Errorf("json 解析错误  %s  %s", err.Error(), val)
		return nil
	}
	// logx.Infof("PaymentSuccess key :%s , val :%s", key, val)
	// 查ip对应的地址
	// 调用户基础方法，获取用户昵称
	var info = logs_model.LogModel{
		LogType: req.LogType,
		IP:      req.IP,
		Addr:    "内网地址",
		UserID:  req.UserID,
		Level:   req.Level,
		Title:   req.Title,
		Content: req.Content,
		Service: req.Service,
	}
	if req.UserID != 0 {
		baseInfo, err1 := l.svcCtx.UserRpc.UserBaseInfo(l.ctx, &user_rpc.UserBaseInfoRequest{
			UserId: uint32(req.UserID),
		})
		if err1 == nil {
			info.UserAvatar = baseInfo.Avatar
			info.UserNickname = baseInfo.NickName
		}
	}

	mutex := sync.Mutex{}
	mutex.Lock()
	err = l.svcCtx.DB.Create(&info).Error
	mutex.Unlock()
	if err != nil {
		logx.Error(err)
		return err
	}

	logx.Infof("日志 %s 保存成功", req.Title)
	return nil
}
