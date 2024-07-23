// Code generated by goctl. DO NOT EDIT.
package types

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type OpenLoginInfoResponse struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Href string `json:"href"` // 跳转地址
}

type OpenLoginRequest struct {
	Code string `json:"code"`
	Flag string `json:"flag"` // 登录标志, 第三方登录标识
}
