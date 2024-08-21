package logic

import (
	"context"
	"strconv"

	"im_server/im_user/user_rpc/internal/svc"
	"im_server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOnlineListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserOnlineListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOnlineListLogic {
	return &UserOnlineListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserOnlineListLogic) UserOnlineList(in *user_rpc.UserOnlineListRequest) (resp *user_rpc.UserOnlineListResponse, err error) {

	resp = new(user_rpc.UserOnlineListResponse)
	// 查哪些用户在线
	onlineMap := l.svcCtx.RedisConf.HGetAll("online").Val()
	for key, _ := range onlineMap {
		val, err1 := strconv.Atoi(key)
		if err1 != nil {
			logx.Error(err1)
			continue
		}
		resp.UserIdList = append(resp.UserIdList, uint32(val))
	}
	return
}
