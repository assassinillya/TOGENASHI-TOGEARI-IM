package models

import (
	"im_server/common/models"
	"im_server/common/models/ctype"
)

// GroupMsgModel 群消息表
type GroupMsgModel struct {
	models.Model
	GroupID    uint             `json:"groupID"`                     //群id
	GroupModel GroupModel       `gorm:"foreignKey:GroupID" json:"-"` //群
	SendUserID uint             `json:"sendUserID"`                  //发送者id
	MsgType    int8             `json:"msgType"`                     // 消息类型 1 文本类型 2 图片消息 3 视频消息 4 文件消息 5 语音消息 6 语音通话 7 视频通话 8 撤回消息 9 回复消息 10 引用消息 11 at消息
	MsgPreView string           `gorm:"size:64" json:"msgPreView"`   //消息预览
	Msg        ctype.Msg        `json:"msg"`                         //消息内容
	SystemMsg  *ctype.SystemMsg `json:"systemMsg"`                   //系统提示
}
