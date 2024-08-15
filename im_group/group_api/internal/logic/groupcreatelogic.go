package logic

import (
	"context"
	"errors"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type Group_createLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroup_createLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Group_createLogic {
	return &Group_createLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Group_createLogic) Group_create(req *types.GroupCreateRequest) (resp *types.GroupCreateResponse, err error) {
	var groupModel = group_models.GroupModel{
		Creator:      req.UserID, // 自己创建的群，自己就是群主
		IsSearch:     false,
		Verification: 2,
		Size:         50,
	}

	switch req.Mode {
	case 1: // 直接创建模式
		if req.Name == "" {
			return nil, errors.New("群名不可为空")
		}
		if req.Size >= 1000 {
			return nil, errors.New("群规模错误")
		}
		groupModel.Title = req.Name
		groupModel.Size = req.Size
		groupModel.IsSearch = req.IsSearch

	case 2:
		// 选人创建模式
		if len(req.UserIDList) == 0 {
			return nil, errors.New("没有要选择的好友")
		}
		// 去算选择的用户昵称，是不是超过最大长度
		// 群名是32
		// 调用户信息列表
		var userIDList = []uint32{uint32(req.UserID)} // 先把自己放进去
		for _, u := range req.UserIDList {
			userIDList = append(userIDList, uint32(u))
		}
		userListResponse, err1 := l.svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
			UserIdList: userIDList,
		})
		if err1 != nil {
			logx.Error(err1)
			return nil, errors.New("用户服务错误")
		}
		// 你选择的这些用户id，是不是都是你的好友？

		// 去算这个昵称的长度 算到第几个人的时候会大于32
		var nameList []string
		for _, info := range userListResponse.UserInfo {
			if len(strings.Join(nameList, "、")) >= 29 {
				break
			}
			nameList = append(nameList, info.NickName)
		}
		groupModel.Title = strings.Join(nameList, "、") + "的群聊"
	}

	// 群头像
	// 1.默认头像  2.文字头像
	groupModel.Avatar = string([]rune(groupModel.Title)[0])
	err = l.svcCtx.DB.Create(&groupModel).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("创建群组失败")
	}

	return
}
