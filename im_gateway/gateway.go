package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var ServiceMap = map[string]string{
	"auth":  "http://127.0.0.1:20021",
	"user":  "http://127.0.0.1:20022",
	"chat":  "http://127.0.0.1:20023",
	"group": "http://127.0.0.1:20024",
}

type Data struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func toJson(data Data) []byte {
	byteData, _ := json.Marshal(data)
	return byteData
}

func gateway(res http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	//   发送请求 /api/user/xxx/
	//   list[0]=  /api/user
	//   list[1]=  user
	regex, _ := regexp.Compile(`/api/(.*?)/`)
	list := regex.FindStringSubmatch(p)
	if len(list) != 2 {
		res.Write(toJson(Data{Code: 7, Msg: "服务错误"}))
		return
	}

	addr, ok := ServiceMap[list[1]]
	if !ok {
		log.Println("不匹配的服务")
		res.Write(toJson(Data{Code: 7, Msg: "服务错误"}))
		return
	}
	// 转发到实际服务上
	url := addr + req.URL.String()

	// 读取请求体
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新放入Body

	log.Printf("请求体: %s", string(bodyBytes)) // 记录请求体

	proxyReq, _ := http.NewRequest(req.Method, url, bytes.NewBuffer(bodyBytes))
	proxyReq.Header = req.Header
	remoteAddr := strings.Split(req.RemoteAddr, ":")
	//if len(remoteAddr) != 2 {
	//	log.Println("err:" ,req.RemoteAddr)
	//	res.Write(toJson(Data{Code: 7, Msg: "服务错误"}))
	//	return
	//}
	log.Printf(`%s%s`, addr, req.URL.String())

	fmt.Printf("%s %s =>  %s\n", remoteAddr[0], list[1], url)
	proxyReq.Header.Set("X-Forwarded-For", remoteAddr[0])
	proxyResponse, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		log.Println("服务异常", err)
		res.Write(toJson(Data{Code: 7, Msg: "服务错误"}))
		return
	}
	io.Copy(res, proxyResponse.Body)
	return
}

func main() {
	// 回调函数
	http.HandleFunc("/", gateway)
	// 绑定服务
	fmt.Printf("fim_gateway 运行在：%s\n", "http://127.0.0.1:9000")
	http.ListenAndServe(":9000", nil)
}
