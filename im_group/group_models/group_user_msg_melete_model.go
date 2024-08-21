package group_models

import "im_server/common/models"

type GroupUserMsgDeleteModel struct {
	models.Model
	UserID  uint `json:"userId"`  // 用户id
	MsgID   uint `json:"msgId"`   // 群聊天记录的id
	GroupID uint `json:"groupId"` // 群id
}
