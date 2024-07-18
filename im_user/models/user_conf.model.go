package models

import (
	"im_server/common/models"
	"im_server/common/models/ctype"
)

// UserConfModel 用户配置表
type UserConfModel struct {
	models.Model
	UserID               uint                        `json:"userID"`
	UserModel            UserModel                   `gorm:"foreignKey:UserID" json:"-"`
	RecallMessage        *string                     `gorm:"size:32" json:"recallMessage"` //撤回消息的提示内容
	FriendOnline         bool                        `json:"friendOnline"`                 //好友上线提醒
	Sound                bool                        `json:"sound"`                        //剩余
	SecureLink           bool                        `json:"secureLink"`                   //安全链接
	SavePwd              bool                        `json:"savePwd"`                      //保存密码
	SearchUser           int8                        `json:"searchUser"`                   //别人查找到你的方式 0 不允许别人查找到我 1 通过用户号查找到我 2 可以通过昵称查找到我
	Verification         int8                        `json:"Verification"`                 //好友验证 0 不允许任何人添加 1允许任何人添加 2 需要验证信息 3 需要回答问题 4 需要正确回答问题
	VerificationQuestion *ctype.VerificationQuestion `json:"verificationQuestion"`         //验证问题 当好友验证为3和4的时候需要
	Online               bool                        `json:"online"`                       //是否在线
}

type VerificationQuestion struct {
	Problem1 *string `json:"problem1"`
	Problem2 *string `json:"problem2"`
	Problem3 *string `json:"problem3"`
	Answer1  *string `json:"answer1"`
	Answer2  *string `json:"answer2"`
	Answer3  *string `json:"answer3"`
}
