// Code generated by goctl. DO NOT EDIT.
package types

type AddFriendRequest struct {
	UserID               uint                  `header:"User-ID"`
	FriendID             uint                  `json:"friendID"`
	Verify               string                `json:"verify,optional"`               // 验证消息
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion,optional"` // 问题和答案
}

type AddFriendResponse struct {
}

type FriendDeleteRequest struct {
	UserID   uint `header:"User-ID"`
	FriendID uint `json:"friendID"`
}

type FriendDeleteResponse struct {
}

type FriendInfoRequest struct {
	UserID   uint `header:"User-ID"`
	Role     int8 `header:"Role"`
	FriendID uint `form:"friendID"` // 好友的用户ID
}

type FriendInfoResponse struct {
	UserID   uint   `json:"userID`
	NickName string `json:"nickname`
	Abstract string `json:"abstract"`
	Avatar   string `json:"avatar"`
	Notice   string `json:"notice"`
	IsOnline bool   `json:"isOnline"` // 是否在线
}

type FriendListRequest struct {
	UserID uint `header:"User-ID"`
	Role   int8 `header:"Role"`
	Page   int  `form:"page,optional"`
	Limit  int  `form:"limit,optional"`
}

type FriendListResponse struct {
	List  []FriendInfoResponse `json:"list"`
	Count int                  `json:"count"`
}

type FriendNoticeUpdateRequest struct {
	UserID   uint   `header:"User-ID"`
	FriendID uint   `json:"friendID"`
	Notice   string `json:"notice"` // 备注
}

type FriendNoticeUpdateResponse struct {
}

type FriendValidInfo struct {
	UserID               uint                  `json:"userID`
	NickName             string                `json:"nickname`
	Avatar               string                `json:"avatar"`
	AdditionalMessages   string                `json:"additionalMessages"`   // 附加消息
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion"` // 验证问题 为3/4的时候需要回答
	Status               int8                  `json:"status"`               // 状态 0:未操作 1:同意 2:拒绝 3:忽略
	Verification         int8                  `json:"verification"`         // 好友验证
	ID                   uint                  `json:"id"`                   // 验证记录的id
	Flag                 string                `json:"flag"`                 // send 我是发起方 rev 我是接收方
	CreatedAt            string                `json:"createdAt"`
}

type FriendValidRequest struct {
	UserID uint `header:"User-ID"`
	Page   int  `form:"page,optional"`
	Limit  int  `form:"limit,optional"`
}

type FriendValidResponse struct {
	List  []FriendValidInfo `json:"list"`
	Count int64             `json:"count"`
}

type FriendValidStatusRequest struct {
	UserID   uint `header:"User-ID"`
	VerifyID uint `json:"verifyID"`
	Status   int8 `json:"status"` // 状态 0:未操作 1:同意 2:拒绝 3:忽略 4:删除
}

type SearchInfo struct {
	UserID   uint   `json:"userID`
	NickName string `json:"nickname`
	Abstract string `json:"abstract"`
	Avatar   string `json:"avatar"`
	IsFriend bool   `json:"isFriend"` // 是否为好友
}

type SearchRequest struct {
	Page   int    `form:"page,optional"`
	Limit  int    `form:"limit,optional"`
	Key    string `form:"key,optional"` // 用户id和昵称
	UserID uint   `header:"User-ID"`
	Online bool   `form:"online,optional"` // 是否在线
}

type SearchResponse struct {
	List  []SearchInfo `json:"list"`
	Count int64        `json:"count"`
}

type UserInfoRequest struct {
	UserID uint `header:"User-ID"`
	Role   int8 `header:"Role"`
}

type UserInfoResponse struct {
	UserID               uint                  `json:"userID`
	NickName             string                `json:"nickname`
	Abstract             string                `json:"abstract"`
	Avatar               string                `json:"avatar"`
	RecallMessage        *string               `json:"recallMessage"`
	FriendOnline         bool                  `json:"friendOnline"`
	Sound                bool                  `json:"sound"`
	SecureLink           bool                  `json:"secureLink"`
	SavePwd              bool                  `json:"savePwd"`
	SearchUser           int8                  `json:"searchUser"`
	Verification         int8                  `json:"verification"`
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion"`
}

type UserInfoUpdateRequest struct {
	UserID               uint                  `header:"User-ID"`
	NickName             *string               `json:"nickname,optional" user:"nickname"`
	Abstract             *string               `json:"abstract,optional" user:"abstract"`
	Avatar               *string               `json:"avatar,optional" user:"avatar"`
	RecallMessage        *string               `json:"recallMessage,optional" user_conf:"recall_message"`
	FriendOnline         *bool                 `json:"friendOnline,optional" user_conf:"friend_online"`
	Sound                *bool                 `json:"sound,optional" user_conf:"sound"`
	SecureLink           *bool                 `json:"secureLink,optional" user_conf:"secure_link"`
	SavePwd              *bool                 `json:"savePwd,optional" user_conf:"save_pwd"`
	SearchUser           *int8                 `json:"searchUser,optional" user_conf:"search_user"`
	Verification         *int8                 `json:"verification,optional" user_conf:"verification"`
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion,optional" user_conf:"verification_question"`
}

type UserInfoUpdateResponse struct {
}

type UserValidRequest struct {
	UserID   uint `header:"User-ID"`
	FriendID uint `json:"friendID"`
}

type UserValidResponse struct {
	Verification         int8                 `json:"verification"`         // 好友验证
	VerificationQuestion VerificationQuestion `json:"verificationQuestion"` // 问题和答案, 注意不要暴露答案
}

type VerificationQuestion struct {
	Problem1 *string `json:"problem1,optional" user_conf:"problem1"`
	Problem2 *string `json:"problem2,optional" user_conf:"problem2"`
	Problem3 *string `json:"problem3,optional" user_conf:"problem3"`
	Answer1  *string `json:"answer1,optional" user_conf:"answer1"`
	Answer2  *string `json:"answer2,optional" user_conf:"answer2"`
	Answer3  *string `json:"answer3,optional" user_conf:"answer3"`
}
