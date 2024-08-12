package chat_models

import (
	"im_server/common/models"
	"im_server/common/models/ctype"
)

type ChatModel struct {
	models.Model
	SendUserID uint             `json:"sendUserID"`
	RevUserID  uint             `json:"revUserID"`
	MsgType    ctype.MsgType    `json:"msgType"`                                      // 消息类型 1 文本类型 2 图片消息 3 视频消息 4 文件消息 5 语音消息 6 语音通话 7 视频通话 8 撤回消息 9 回复消息 10 引用消息
	MsgPreView string           `gorm:"column:msg_preview;size:64" json:"msgPreview"` //消息预览
	Msg        *ctype.Msg       `json:"msg"`
	SystemMsg  *ctype.SystemMsg `json:"systemMsg"` //系统提示
}

func (chat ChatModel) MsgPreviewMethod() string {
	if chat.SystemMsg != nil {
		switch chat.SystemMsg.Type {
		case 1:
			return "[系统消息]- 该消息涉黄，已被系统拦截"
		case 2:
			return "[系统消息]- 该消息涉恐，已被系统拦截"
		case 3:
			return "[系统消息]- 该消息涉政，已被系统拦截"
		case 4:
			return "[系统消息]- 该消息存在不正当言论，已被系统拦截"
		}
		return "[系统消息]"
	}
	switch chat.MsgType {
	case 1:
		return chat.Msg.TextMsg.Content
	case 2:
		return "[图片消息] - " + chat.Msg.ImageMsg.Title
	case 3:
		return "[视频消息] - " + chat.Msg.VideoMsg.Title
	case 4:
		return "[文件消息] - " + chat.Msg.FileMsg.Title
	case 5:
		return "[语音消息] - "
	case 6:
		return "[语音通话] - "
	case 7:
		return "[视频通话] - "
	case 8:
		return "[撤回消息] - " + chat.Msg.WithdrawMsg.Content
	case 9:
		return "[回复消息] - " + chat.Msg.ReplyMsg.Content
	case 10:
		return "[引用消息] - " + chat.Msg.QuoteMsg.Content
	case 11:
		return "[@消息] - " + chat.Msg.AtMsg.Content
	}
	return "[未知消息]"
}
