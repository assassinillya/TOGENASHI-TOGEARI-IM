package main

import (
	"flag"
	"fmt"
	"im_server/core"
	"im_server/im_chat/chat_models"
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
		db := core.InitMysql()
		err := db.AutoMigrate(
			&user_models.UserModel{},
			&user_models.FriendModel{},
			&user_models.FriendVerifyModel{},
			&user_models.UserConfModel{},
			&chat_models.ChatModel{},
			&group_models.GroupModel{},
			&group_models.GroupMemberModel{},
			&group_models.GroupMsgModel{},
			&group_models.GroupVerifyModel{},
		)

		if err != nil {
			fmt.Println("表结构生成失败", err)
			return
		}
		log.Println("表结构生成成功！")
		fmt.Println("表结构生成成功！")

	}
}
