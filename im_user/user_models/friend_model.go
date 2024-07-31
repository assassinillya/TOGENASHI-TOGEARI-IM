package user_models

import (
	"gorm.io/gorm"
	"im_server/common/models"
)

// FriendModel 好友表
type FriendModel struct {
	models.Model
	SendUserID     uint      `json:"sendUserID"`                     //发起验证方
	SendUserModel  UserModel `gorm:"foreignKey:SendUserID" json:"-"` //发起验证方
	RevUserID      uint      `json:"revUserID"`                      //接受验证方
	RevUserModel   UserModel `gorm:"foreignKey:RevUserID" json:"-"`  //接受验证方
	SendUserNotice string    `gorm:"size:128" json:"sendUserNotice"` //A发送方备注 A对B的备注
	RevUserNotice  string    `gorm:"size:128" json:"revUserNotice"`  //B接收方备注 B对A的备注
}

func (f *FriendModel) IsFriend(db *gorm.DB, A, B uint) bool {
	err := db.Take(&f, "(send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?)", A, B, B, A).Error
	if err != nil {
		return false
	}
	return true
}

func (f *FriendModel) GetUserNotice(userID uint) string {
	if userID == f.SendUserID {
		// 如果我是发起方
		return f.SendUserNotice
	}
	if userID == f.RevUserID {
		// 如果我是接收方
		return f.RevUserNotice
	}
	return ""
}
