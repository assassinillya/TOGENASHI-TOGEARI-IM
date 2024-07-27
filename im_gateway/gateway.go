package main

import (
	"bytes"
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

func gateway(res http.ResponseWriter, req *http.Request) {
	regex, _ := regexp.Compile(`/api/(.*?)/`) // 匹配请求路径 /api/user/xx
	addrList := regex.FindStringSubmatch(req.URL.Path)
	if len(addrList) != 2 {
		res.Write([]byte("err"))
		return
	}
	service := addrList[1]
	addr := etcd.GetServiceAddr(config.Etcd, service+"_api")
	log.Println("Service address from etcd:", addr) // 打印获取的地址
	if addr == "" {
		fmt.Println("不匹配的服务", service)
		res.Write([]byte("不匹配的服务"))
		return
	}

	remoteAddr := strings.Split(req.RemoteAddr, ":")
	log.Println("remoteAddr", remoteAddr)

	//请求认证服务地址
	authAddr := etcd.GetServiceAddr(config.Etcd, "auth_api")
	//authAddr:="127.0.0.1:20023"
	authUrl := fmt.Sprintf("http://%s/api/auth/authentication", authAddr)
	authReq, _ := http.NewRequest("POST", authUrl, nil)
	authReq.Header = req.Header
	authReq.Header.Set("ValidPath", req.URL.Path)

	log.Println("Token:", req.Header.Get("Authorization")) //打印请求头中的Token

	authRes, err := http.DefaultClient.Do(authReq)
	if err != nil {
		log.Println("认证服务错误 ", err)
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

	byteData, _ = io.ReadAll(req.Body)

	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(byteData))
	if err != nil {
		logx.Error(err)
		res.Write([]byte("服务异常"))
		return
	}
	proxyReq.Header = req.Header
	proxyReq.Header.Del("ValidPath")
	response, ProxyErr := http.DefaultClient.Do(proxyReq)
	if ProxyErr != nil {
		fmt.Println(ProxyErr)
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
