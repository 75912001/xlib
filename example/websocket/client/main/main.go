package main

import (
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	// 连接到 WebSocket 服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()

	// 发送消息
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
	if err != nil {
		log.Fatal("发送消息失败:", err)
	}

	// 读取响应
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("读取消息失败:", err)
	}
	log.Printf("收到响应: %s", message)
}
