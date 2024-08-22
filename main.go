package main

import (
	"flag"
	"fmt"
	"im_server/core"
	"im_server/im_chat/chat_models"
	"im_server/im_file/file_model"
	"im_server/im_group/group_models"
	"im_server/im_user/user_models"
	"log"
)

type Options struct {
	DB bool
}

func main() {
	var opt Options
	flag.BoolVar(&opt.DB, "db", false, "db")
	flag.Parse()

	if opt.DB {
		//db := core.InitMysql()
		db := core.InitGorm("root:123456@tcp(120.78.234.30:3306)/im_server_db?charset=utf8mb4&parseTime=True&loc=Local")
		err := db.AutoMigrate(
			&user_models.UserModel{},                // 用户表
			&user_models.FriendModel{},              // 好友表
			&user_models.FriendVerifyModel{},        // 好友验证表
			&user_models.UserConfModel{},            // 用户配置表
			&chat_models.ChatModel{},                // 对话表
			&chat_models.TopUserModel{},             // 置顶用户表
			&chat_models.UserChatDeleteModel{},      // 用户删除聊天记录表
			&group_models.GroupModel{},              // 群组表
			&group_models.GroupMemberModel{},        // 群成员表
			&group_models.GroupMsgModel{},           // 群消息表
			&group_models.GroupVerifyModel{},        // 群验证表
			&group_models.GroupUserMsgDeleteModel{}, // 用户删除聊天记录表
			&group_models.GroupUserTopModel{},       // 用户置顶群聊表
			&file_model.FileModel{},                 // 文件表
		)

		if err != nil {
			fmt.Println("表结构生成失败", err)
			return
		}
		log.Println("表结构生成成功！")
		fmt.Println("表结构生成成功！")

	}
}
