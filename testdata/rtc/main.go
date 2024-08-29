package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var offerClient *websocket.Conn
var answerClient *websocket.Conn

func checkStart() {
	if offerClient != nil && answerClient != nil {
		offerClient.WriteJSON(map[string]string{
			"type": "create_offer",
		})
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 鉴权 true表示放行，false表示拦截
			return true
		},
	}
	fmt.Println("建立ws连接")
	// 将http升级至websocket
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	for {
		var obj map[string]any
		err := conn.ReadJSON(&obj)
		if err != nil {
			log.Println("Error reading JSON:", err)
			break
		}

		log.Println("recv:", obj)

		switch obj["type"] {
		case "connect":
			if offerClient == nil {
				offerClient = conn
				conn.WriteJSON(map[string]interface{}{
					"type":    "connect",
					"code":    200,
					"message": "connect success",
				})
				checkStart()
			} else if answerClient == nil {
				answerClient = conn
				conn.WriteJSON(map[string]interface{}{
					"type":    "connect",
					"code":    200,
					"message": "connect success",
				})
				checkStart()
			} else {
				conn.WriteJSON(map[string]interface{}{
					"type":    "connect",
					"code":    -1,
					"message": "connect failed",
				})
				conn.Close()
			}
		case "offer":
			if answerClient != nil {
				answerClient.WriteJSON(obj)
			}
		case "answer":
			if offerClient != nil {
				offerClient.WriteJSON(obj)
			}
		case "offer_ice":
			if answerClient != nil {
				answerClient.WriteJSON(obj)
			}
		case "answer_ice":
			if offerClient != nil {
				offerClient.WriteJSON(obj)
			}
		}
	}

	if conn == offerClient {
		log.Println("remove offerClient")
		offerClient = nil
	} else if conn == answerClient {
		log.Println("remove answerClient")
		answerClient = nil
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(123)
	byteData, err := os.ReadFile("C:/Users/Assas/Desktop/code/TOGENASHI-TOGEARI-IM/TOGENASHI-TOGEARI-IM/testdata/rtc/static/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(byteData)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ws", wsHandler)
	log.Println("Server running on :9004")
	log.Fatal(http.ListenAndServe(":9004", nil))
}
