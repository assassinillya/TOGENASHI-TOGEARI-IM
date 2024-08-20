// Code generated by goctl. DO NOT EDIT.
package types

type AddGroupRequest struct {
	UserID               uint                  `header:"User-ID"`
	GroupID              uint                  `json:"groupID"`
	Verify               string                `json:"verify,optional"`               // 验证消息
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion,optional"` // 问题和答案
}

type AddGroupResponse struct {
}

type GroupFriendsResponse struct {
	UserID    uint   `json:"userId"`
	Avatar    string `json:"avatar"`
	Nickname  string `json:"nickname"`
	IsInGroup bool   `json:"isInGroup"` //是否在群里面
}

type GroupMemberInfo struct {
	UserID         uint   `json:"userId"`
	UserNickname   string `json:"userNickname"`
	Avatar         string `json:"avatar"`
	IsOnline       bool   `json:"isOnline"`
	Role           int8   `json:"role"`
	MemberNickname string `json:"memberNickname"`
	CreatedAt      string `json:"createdAt"`
	NewMsgDate     string `json:"newMsgDate"`
}

type GroupSearchResponse struct {
	GroupID         uint   `json:"groupId`
	Title           string `json:"title`
	Abstract        string `json:"abstract"`
	Avatar          string `json:"avatar"`
	IsInGroup       bool   `json:"isInGroup"`       // 我是否在群里
	UserCount       int    `json:"userCount"`       // 群用户总数
	UserOnlineCount int    `json:"userOnlineCount"` // 群在线用户总数
}

type UserInfo struct {
	UserID   uint   `header:"User-ID"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

type VerificationQuestion struct {
	Problem1 *string `json:"problem1,optional" conf:"problem1"`
	Problem2 *string `json:"problem2,optional" conf:"problem2"`
	Problem3 *string `json:"problem3,optional" conf:"problem3"`
	Answer1  *string `json:"answer1,optional" conf:"answer1"`
	Answer2  *string `json:"answer2,optional" conf:"answer2"`
	Answer3  *string `json:"answer3,optional" conf:"answer3"`
}

type GroupCreateRequest struct {
	UserID     uint   `header:"User-ID"`
	Mode       int8   `json:"mode,optional"` // 模式 1 直接创建模式 2 选人创建模式
	Name       string `json:"name,optional"` // 群聊名字
	IsSearch   bool   `json:"isSearch,optional"`
	Size       int    `json:"size,optional"`       // 群规模
	UserIDList []uint `json:"userIdList,optional"` // 用户id列表
}

type GroupCreateResponse struct {
}

type GroupFriendsListRequest struct {
	UserID uint `header:"User-ID"` //自己的id
	ID     uint `form:"id"`        //群id
}

type GroupFriendsListResponse struct {
	Count int                    `json:"count"`
	List  []GroupFriendsResponse `json:"list"`
}

type GroupInfoRequest struct {
	UserID uint `header:"User-ID"`
	ID     uint `path:"id"` //群id
}

type GroupInfoResponse struct {
	GroupID           uint       `json:"groupId"`           //群id
	Title             string     `json:"title"`             //群名称
	Abstract          string     `json:"abstract"`          //群简介
	Avatar            string     `json:"avatar"`            //群头像
	Creator           UserInfo   `json:"creator"`           //群主
	MemberCount       int        `json:"memberCount"`       //群聊用户总数
	MemberOnlineCount int        `json:"memberOnlineCount"` //在线用户数量
	AdminList         []UserInfo `json:"adminList"`         //管理员列表
	Role              int8       `json:"role"`              // 群角色 1 群主 2 管理员 3 群成员
}

type GroupMemberAddRequest struct {
	UserID       uint   `header:"User-ID"`
	ID           uint   `json:"id"`           //群id
	MemberIDList []uint `json:"memberIdList"` //成员id列表
}

type GroupMemberAddResponse struct {
}

type GroupMemberNicknameUpdateRequest struct {
	UserID   uint   `header:"User-ID"`
	ID       uint   `json:"id"`       //群id
	MemberID uint   `json:"memberId"` //成员id
	Nickname string `json:"nickname"` //成员昵称
}

type GroupMemberNicknameUpdateResponse struct {
}

type GroupMemberRemoveRequest struct {
	UserID   uint `header:"User-ID"`
	ID       uint `form:"id"`       //群id
	MemberID uint `form:"memberId"` //成员id
}

type GroupMemberRemoveResponse struct {
}

type GroupMemberRequest struct {
	UserID uint   `header:"User-ID"`
	ID     uint   `form:"id"` //群id
	Page   int    `form:"page,optional"`
	Limit  int    `form:"limit,optional"`
	Sort   string `form:"sort,optional"`
}

type GroupMemberResponse struct {
	List  []GroupMemberInfo `json:"list"`
	Count int               `json:"count"`
}

type GroupMemberRoleUpdateRequest struct {
	UserID   uint `header:"User-ID"`
	ID       uint `json:"id"`       //群id
	MemberID uint `json:"memberId"` //成员id
	Role     int8 `json:"role"`
}

type GroupMemberRoleUpdateResponse struct {
}

type GroupRemoveRequest struct {
	UserID uint `header:"User-ID"`
	ID     uint `path:"id"` //群id
}

type GroupRemoveResponse struct {
}

type GroupSearchListResponse struct {
	List  []GroupSearchResponse `json:"list"`
	Count int                   `json:"count"`
}

type GroupSearchRequest struct {
	Page   int    `form:"page,optional"`
	Limit  int    `form:"limit,optional"`
	Key    string `form:"key,optional"` // 用户id和昵称
	UserID uint   `header:"User-ID"`
}

type GroupUpdateRequest struct {
	UserID               uint                  `header:"User-ID"`
	ID                   uint                  `json:"id"`                                                      // 群id
	IsSearch             *bool                 `json:"isSearch,optional" conf:"is_search"`                      // 是否可以被搜索
	Avatar               *string               `json:"avatar,optional" conf:"avatar"`                           // 群头像
	Abstract             *string               `json:"abstract,optional" conf:"abstract"`                       // 群简介
	Title                *string               `json:"title,optional" conf:"title"`                             // 群名称
	Verification         *int8                 `json:"verification,optional" conf:"verification"`               // 群验证
	IsInvite             *bool                 `json:"isInvite,optional" conf:"is_invite"`                      // 是否可邀请好友
	IsTemporarySession   *bool                 `json:"isTemporarySession,optional" conf:"is_temporary_session"` // 是否可临时会话
	IsProhibition        *bool                 `json:"isProhibition,optional" conf:"is_prohibition"`            // 是否全员禁言
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion,optional" conf:"verification_question"`
}

type GroupUpdateResponse struct {
}

type GroupValidRequest struct {
	UserID  uint `header:"User-ID"`
	GroupID uint `form:"groupID"`
}

type GroupValidResponse struct {
	Verification         int8                 `json:"verification"`         // 好友验证
	VerificationQuestion VerificationQuestion `json:"verificationQuestion"` // 问题和答案, 答案不返回
}
