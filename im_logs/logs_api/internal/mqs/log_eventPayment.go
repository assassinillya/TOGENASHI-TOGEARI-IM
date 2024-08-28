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
	LogType int8   `json:"logType"` // 日志类型 2 操作日志 3 运行日志
	IP      string `json:"ip"`
	UserID  uint   `json:"userID"`
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
	// 判断是不是运行日志
	if info.LogType == 3 {
		// 运行日志
		// 先查一下今天这个服务有没有日志, 有的话就更新, 没有再创建
		mutex := sync.Mutex{}
		mutex.Lock()
		var logModel logs_model.LogModel
		err1 := l.svcCtx.DB.Take(&logModel, "log_type = ? and service = ? and to_days(created_at) =  to_days(now())", 3, info.Service).Error
		mutex.Unlock()
		if err1 == nil {
			// 找到了
			l.svcCtx.DB.Model(&logModel).Update("content", logModel.Content+"\n"+info.Content)
			logx.Infof("运行日志 %s 更新成功", req.Title)
			return nil
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
