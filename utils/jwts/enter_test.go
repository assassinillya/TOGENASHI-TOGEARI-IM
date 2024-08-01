package jwts

import (
	"fmt"
	"testing"
)

// 测试生成Token的函数
func TestGenToken(t *testing.T) {
	payload := JwtPayLoad{
		UserID:   1,
		Nickname: "asily",
		Role:     2,
	}
	accessSecret := "123456"
	expires := 10000000 // token有效期1小时

	token, err := GenToken(payload, accessSecret, expires)
	fmt.Println(token)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	if token == "" {
		t.Fatalf("Generated token is empty")
	}
}

// 测试解析Token的函数
func TestParseToken(t *testing.T) {
	payload := JwtPayLoad{
		UserID:   1,
		Nickname: "testuser",
		Role:     2,
	}
	accessSecret := "123456"
	expires := 1 // token有效期1小时

	token, err := GenToken(payload, accessSecret, expires)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	parsedClaims, err := ParseToken(token, accessSecret)
	if err != nil {
		t.Fatalf("Error parsing token: %v", err)
	}

	if parsedClaims.UserID != payload.UserID || parsedClaims.Nickname != payload.Nickname || parsedClaims.Role != payload.Role {
		t.Fatalf("Parsed claims do not match original payload")
	}
}
