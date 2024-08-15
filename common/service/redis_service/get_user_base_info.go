package redis_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"im_server/common/models/ctype"
	"im_server/im_user/user_rpc/types/user_rpc"
	"time"
)

// GetUserBaseInfo 通过redis获取用户基本信息
func GetUserBaseInfo(client *redis.Client, UserRpc user_rpc.UsersClient, userID uint) (userInfo ctype.UserInfo, err error) {
	key := fmt.Sprintf("fim_server_user_%d", userID)
	str, err := client.Get(key).Result()
	if err != nil {
		//没找到,调用户rpc服务，数据库去查，查完之后再设置进redis缓存
		userBaseResponse, err1 := UserRpc.UserBaseInfo(context.Background(), &user_rpc.UserBaseInfoRequest{
			UserId: uint32(userID),
		})
		if err1 != nil {
			err = err1
			return
		}
		err = nil
		userInfo.ID = userID
		userInfo.Avatar = userBaseResponse.Avatar
		userInfo.NickName = userBaseResponse.NickName

		byteData, _ := json.Marshal(userInfo)
		//设置进缓存
		client.Set(key, string(byteData), time.Hour) //存一个小时的用户基本信息
		return
	}
	err = json.Unmarshal([]byte(str), &userInfo)
	if err != nil {
		return
	}

	return
}
