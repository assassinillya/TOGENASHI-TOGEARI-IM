package logic

import (
	"context"
	"errors"
	"fmt"
	"im_server/im_group/group_api/internal/svc"
	"im_server/im_group/group_api/internal/types"
	"im_server/im_group/group_models"
	"im_server/im_user/user_rpc/types/user_rpc"
	"im_server/utils/set"
	"strings"
	"time"

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
		Abstract:     fmt.Sprintf("本群创建于%s: 群主很懒, 什么都没有留下", time.Now().Format(time.DateOnly)),
		IsSearch:     false,
		Verification: 2,
		Size:         50,
	}

	var groupUserList = []uint{req.UserID}

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
			groupUserList = append(groupUserList, u)
		}

		userFriendResponse, err2 := l.svcCtx.UserRpc.FriendList(l.ctx, &user_rpc.FriendListRequest{
			User: uint32(req.UserID),
		})
		if err2 != nil {
			logx.Error(err2)
			return nil, errors.New("用户服务错误")
		}

		var friendIDList []uint
		for _, i2 := range userFriendResponse.FriendList {
			friendIDList = append(friendIDList, uint(i2.UserId))
		}

		// 判断他们两个是不是一致的
		slice := set.Difference(req.UserIDList, friendIDList)
		if len(slice) != 0 {
			return nil, errors.New("你选择的好友中有人不是你的好友")
		}

		userListResponse, err1 := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
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

	default:
		return nil, errors.New("不支持的模式")
	}

	// 群头像
	// 1.默认头像  2.文字头像
	groupModel.Avatar = string([]rune(groupModel.Title)[0])
	err = l.svcCtx.DB.Create(&groupModel).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("创建群组失败")
	}

	var members []group_models.GroupMemberModel
	for i, u := range groupUserList {
		memberModel := group_models.GroupMemberModel{
			GroupID: groupModel.ID,
			UserID:  u,
			Role:    3,
		}
		if i == 0 {
			memberModel.Role = 1
		}
		members = append(members, memberModel)
	}

	err = l.svcCtx.DB.Create(&members).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("群成员添加失败")
	}

	return
}
