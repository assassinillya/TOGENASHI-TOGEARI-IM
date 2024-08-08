package chat_models

// TopUserModel 置顶用户表
type TopUserModel struct {
	UserID    uint `json:"userID"`
	TopUserID uint `json:"topUserID"`
}
