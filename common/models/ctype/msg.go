package ctype

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type MsgType = int8

const (
	TextMsgType MsgType = iota + 1
	ImageMsgType
	VideoMsgType
	FileMsgType
	VoiceMsgType
	VoiceCallMsgType
	VideoCallMsgType
	WithdrawMsgType
	ReplyMsgType
	QuoteMsgType
	AtMsgType
	TipMsgType
)

type Msg struct {
	Type         int8          `json:"type"`         // 消息类型 和 MsgType 一样
	Content      *string       `json:"content"`      // 为1时使用
	TextMsg      *TextMsg      `json:"textMsg"`      // 文本消息
	ImageMsg     *ImageMsg     `json:"imageMsg"`     // 图片消息
	VideoMsg     *VideoMsg     `json:"videoMsg"`     // 视频消息
	FileMsg      *FileMsg      `json:"fileMsg"`      // 文件消息
	VoiceMsg     *VoiceMsg     `json:"voiceMsg"`     // 语音消息
	VoiceCallMsg *VoiceCallMsg `json:"voiceCallMsg"` // 语音通话
	VideoCallMsg *VideoCallMsg `json:"videoCallMsg"` // 视频通话
	WithdrawMsg  *WithdrawMsg  `json:"withdrawMsg"`  // 撤回消息
	ReplyMsg     *ReplyMsg     `json:"replyMsg"`     // 回复消息
	QuoteMsg     *QuoteMsg     `json:"quoteMsg"`     // 视频信息
	AtMsg        *AtMsg        `json:"atMsg"`        // @用户的消息 群聊才有
}

// Scan 入库的数据
func (c *Msg) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}

// Value 入库的数据
func (c *Msg) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

type TextMsg struct {
	Content string `json:"content"`
}

type ImageMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
}

type VideoMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
	Time  int    `json:"time"` //单位秒
}

type FileMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
	Size  int64  `json:"size"` //文件大小
	Type  string `json:"type"` //文件类型

}

type VoiceMsg struct {
	Src  string `json:"src"`
	Time int    `json:"time"` //单位秒
}

type VoiceCallMsg struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	EndReason int8      `json:"endReason"` //结束原因 0 发起方挂断 1 接收方挂断 2 网络原因挂断 3 未打通
}

type VideoCallMsg struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	EndReason int8      `json:"endReason"` //结束原因 0 发起方挂断 1 接收方挂断 2 网络原因挂断 3 未打通

}

// WithdrawMsg 撤回消息
type WithdrawMsg struct {
	Content   string `json:"content"` // 撤回的提示词
	OriginMsg *Msg   `json:"-"`       //原消息

}

type ReplyMsg struct {
	MsgID   uint   `json:"msgID"`   //信息id
	Content string `json:"content"` //回复的文本消息，目前只能限制回复文本
	Msg     *Msg   `json:"msg"`
}

type QuoteMsg struct {
	MsgID   uint   `json:"msgID"`   //信息id
	Content string `json:"content"` //回复的文本消息，目前只能限制回复文本
	Msg     *Msg   `json:"msg"`
}

// AtMsg @消息
type AtMsg struct {
	UserID  uint   `json:"userID"`
	Content string `json:"content"` //回复的文本消息
	Msg     *Msg   `json:"msg"`
}
