package logic

import (
	"context"
	"errors"
	"im_server/common/models/ctype"
	"im_server/im_user/user_api/internal/svc"
	"im_server/im_user/user_api/internal/types"
	"im_server/im_user/user_models"
	"im_server/utils/maps"
	"log"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoUpdateLogic {
	return &UserInfoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoUpdateLogic) UserInfoUpdate(req *types.UserInfoUpdateRequest) (resp *types.UserInfoUpdateResponse, err error) {
	log.Println(req.UserID)
	userMaps := maps.RefToMap(*req, "user")
	if len(userMaps) != 0 {
		var user user_models.UserModel
		err = l.svcCtx.DB.Take(&user, req.UserID).Error
		if err != nil {
			return nil, errors.New("用户不存在")
		}

		err = l.svcCtx.DB.Model(&user).Updates(userMaps).Error
		if err != nil {
			logx.Error(userMaps)
			logx.Error(err)
			return nil, errors.New("用户信息更新失败")
		}
	}

	userConfMaps := maps.RefToMap(*req, "user_conf")
	if len(userConfMaps) != 0 {
		var userConf user_models.UserConfModel
		err = l.svcCtx.DB.Take(&userConf, "user_id = ?", req.UserID).Error
		if err != nil {
			return nil, errors.New("用户配置不存在")
		}

		verificationQuestion, ok := userConfMaps["verification_question"]
		if ok {
			delete(userConfMaps, "verification_question")

			data := ctype.VerificationQuestion{}
			maps.MapToStruct(verificationQuestion.(map[string]any), &data)

			if val, ok := verificationQuestion.(map[string]any)["problem1"]; ok {
				s := val.(string)
				data.Problem1 = &s
			}
			if val, ok := verificationQuestion.(map[string]any)["problem2"]; ok {
				s := val.(string)
				data.Problem2 = &s
			}
			if val, ok := verificationQuestion.(map[string]any)["problem3"]; ok {
				s := val.(string)
				data.Problem3 = &s
			}
			if val, ok := verificationQuestion.(map[string]any)["answer1"]; ok {
				s := val.(string)
				data.Answer1 = &s
			}
			if val, ok := verificationQuestion.(map[string]any)["answer2"]; ok {
				s := val.(string)
				data.Answer2 = &s
			}
			if val, ok := verificationQuestion.(map[string]any)["answer3"]; ok {
				s := val.(string)
				data.Answer3 = &s
			}
			err = l.svcCtx.DB.Model(&userConf).Updates(&user_models.UserConfModel{
				VerificationQuestion: &data,
			}).Error
			if err != nil {
				logx.Error("更新用户配置Q&A信息失败")
				return nil, errors.New("更新用户配置Q&A信息失败")
			}
		}

		err = l.svcCtx.DB.Model(&userConf).Updates(userConfMaps).Error
		if err != nil {
			logx.Error("用户配置更新失败")
			return nil, errors.New("用户配置更新失败")
		}
	}

	return
}
