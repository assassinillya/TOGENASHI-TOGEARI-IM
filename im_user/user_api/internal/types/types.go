// Code generated by goctl. DO NOT EDIT.
package types

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

type VerificationQuestion struct {
	Problem1 *string `json:"problem1,optional" user_conf:"problem1"`
	Problem2 *string `json:"problem2,optional" user_conf:"problem2"`
	Problem3 *string `json:"problem3,optional" user_conf:"problem3"`
	Answer1  *string `json:"answer1,optional" user_conf:"answer1"`
	Answer2  *string `json:"answer2,optional" user_conf:"answer2"`
	Answer3  *string `json:"answer3,optional" user_conf:"answer3"`
}
