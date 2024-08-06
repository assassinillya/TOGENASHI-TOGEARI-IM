package logic

import (
	"context"
	"errors"
	"im_server/common/models/ctype"
	"im_server/im_user/user_models"

	"im_server/im_user/user_api/internal/svc"
	"im_server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFriendLogic) AddFriend(req *types.AddFriendRequest) (resp *types.AddFriendResponse, err error) {
	// 判断是否已经是好友了
	var friend user_models.FriendModel
	if friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		return nil, errors.New("已经是好友了")
	}

	var userConf user_models.UserConfModel
	err = l.svcCtx.DB.Take(&userConf, "user_id = ?", req.FriendID).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	resp = new(types.AddFriendResponse)
	var VerifyModel = user_models.FriendVerifyModel{
		SendUserID:         req.UserID,
		RevUserID:          req.FriendID,
		AdditionalMessages: req.Verify,
	}
	switch userConf.Verification {
	case 0: // 不允许任何人添加
		return nil, errors.New("该用户不允许任何人添加")
	case 1: // 允许任何人添加
		// 直接添加好友
		// 先往验证表中添加一条记录, 然后通过
		VerifyModel.Status = 1
	case 2: // 需要验证问题
	case 3: // 需要回答问题, 不需要正确回答问题
		if req.VerificationQuestion != nil {
			VerifyModel.VerificationQuestion = &ctype.VerificationQuestion{
				Problem1: req.VerificationQuestion.Problem1,
				Problem2: req.VerificationQuestion.Problem2,
				Problem3: req.VerificationQuestion.Problem3,
				Answer1:  req.VerificationQuestion.Answer1,
				Answer2:  req.VerificationQuestion.Answer2,
				Answer3:  req.VerificationQuestion.Answer3,
			}
		}

	case 4: // 需要回答问题, 需要正确回答问题
		// 判断问题是否回答正确
		if req.VerificationQuestion != nil && userConf.VerificationQuestion != nil {
			if req.VerificationQuestion.Answer1 != userConf.VerificationQuestion.Answer1 ||
				req.VerificationQuestion.Answer2 != userConf.VerificationQuestion.Answer2 ||
				req.VerificationQuestion.Answer3 != userConf.VerificationQuestion.Answer3 {
				return nil, errors.New("答案错误")
			}
			// 直接添加好友
			VerifyModel.Status = 1
			VerifyModel.VerificationQuestion = userConf.VerificationQuestion
			// 添加好友
			var userFriend = user_models.FriendModel{
				SendUserID: req.UserID,
				RevUserID:  req.FriendID,
			}
			l.svcCtx.DB.Create(&userFriend)
		}
	default:

	}
	err = l.svcCtx.DB.Create(&VerifyModel).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("添加好友失败")
	}

	return
}
