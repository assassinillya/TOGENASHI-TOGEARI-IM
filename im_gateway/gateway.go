package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"im_server/common/etcd"
	"io"
	"log"
	http "net/http"
	"regexp"
	"strings"
)

var serviceMap = map[string]string{
	"auth": "http://127.0.0.1:20021",
	"user": "http://127.0.0.1:20022",
}

func gateway(res http.ResponseWriter, req *http.Request) {
	// 匹配请求路径 /api/user/xx
	regex, _ := regexp.Compile(`/api/(.*?)/`)
	addrList := regex.FindStringSubmatch(req.URL.Path)
	if len(addrList) != 2 {
		res.Write([]byte("err"))
		return
	}

	service := addrList[1]

	addr := etcd.GetServiceAddr(config.Etcd, service+"_api")
	if addr == "" {
		fmt.Println("不匹配的服务", service)
		res.Write([]byte("err"))
		return
	}

	remoteAddr := strings.Split(req.RemoteAddr, ":")
	log.Println("remoteAddr", remoteAddr)

	//请求认证服务地址
	authAddr := etcd.GetServiceAddr(config.Etcd, "auth_api")
	authUrl := fmt.Sprintf("http://%s/api/auth/authentication", authAddr)
	authReq, _ := http.NewRequest("POST", authUrl, req.Body)
	authReq.Header.Set("X-Forwarded-For", remoteAddr[0])
	authRes, err := http.DefaultClient.Do(authReq)
	if err != nil {
		logx.Error(err)
		res.Write([]byte("认证服务错误"))
		return
	}

	type Response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	var authResponse Response
	byteData, _ := io.ReadAll(authRes.Body)
	authErr := json.Unmarshal(byteData, &authResponse)
	if authErr != nil {
		logx.Error(err)
		res.Write([]byte("认证服务错误"))
		return
	}

	// 认证不通过
	if authResponse.Code != 0 {
		res.Write(byteData)
		return
	}

	url := fmt.Sprintf("http://%s%s", addr, req.URL.String())
	fmt.Println(url)
	proxyReq, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		logx.Error(err)
		res.Write([]byte("服务异常"))
		return
	}
	response, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		fmt.Println(err)
		res.Write([]byte("服务异常"))
		return
	}
	io.Copy(res, response.Body)
}

var configFile = flag.String("f", "settings.yaml", "the config file")

type Config struct {
	Addr string
	Etcd string
}

var config Config

func main() {
	flag.Parse()
	conf.MustLoad(*configFile, &config)

	// 回调函数
	http.HandleFunc("/", gateway)
	fmt.Printf("gateway running %s\n", config.Addr)
	// 绑定服务
	http.ListenAndServe(config.Addr, nil)
}
