package chat_models

// UserChatDeleteModel 用户删除聊天记录表
type UserChatDeleteModel struct {
	UserID uint `json:"user_id"`
	ChatID uint `json:"chat_id"` //聊天记录的id

}
