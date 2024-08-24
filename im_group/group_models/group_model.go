package group_models

import (
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"im_server/common/models"
	"im_server/common/models/ctype"
	"time"
)

// GroupModel 群组
type GroupModel struct {
	models.Model
	Title                string                      `gorm:"size:32" json:"title"`        // 群名
	Abstract             string                      `gorm:"size:128" json:"abstract"`    // 简介
	Avatar               string                      `gorm:"size:256" json:"avatar"`      // 群头像
	Creator              uint                        `json:"creator"`                     // 群主
	IsSearch             bool                        `json:"isSearch"`                    // 是否可以被搜索
	Verification         int8                        `json:"verification"`                // 群验证 0 不允许任何人添加 1允许任何人添加 2 需要验证信息 3 需要回答问题 4 需要正确回答问题
	VerificationQuestion *ctype.VerificationQuestion `json:"verificationQuestion"`        // 验证问题 当好友验证为3和4的时候需要
	IsInvite             bool                        `json:"isInvite"`                    // 群成员是否可邀请好友
	IsTemporarySession   bool                        `json:"isTemporarySession"`          // 是否开启临时会话
	IsProhibition        bool                        `json:"isProhibition"`               // 是否开启全员禁言
	Size                 int                         `json:"size"`                        // 群规模 20 100 200 1000 2000
	MemberList           []GroupMemberModel          `gorm:"foreignKey:GroupID" json:"-"` // 群成员列表
}

func (uc GroupModel) GetQuestionCount() (count int) {
	if uc.VerificationQuestion != nil {
		if uc.VerificationQuestion.Answer1 != nil {
			count += 1
		}
		if uc.VerificationQuestion.Answer2 != nil {
			count += 1
		}
		if uc.VerificationQuestion.Answer3 != nil {
			count += 1
		}
	}
	return count
}

func (gm GroupMemberModel) GetProhibitionTime(client *redis.Client, db *gorm.DB) *int {
	if gm.ProhibitionTime == nil {
		return nil
	}
	t, err := client.TTL(fmt.Sprintf("prohibition__%d", gm.ID)).Result()
	if err != nil {
		// 查不到就说明过期了
		// 就把这个值改回去
		db.Model(&gm).Update("prohibition_time", nil)
		return nil
	}
	// -2在redis中代表键不存在
	if t == -2*time.Second {
		db.Model(&gm).Update("prohibition_time", nil)
		return nil
	}
	res := int(t / time.Minute)
	return &res
}
