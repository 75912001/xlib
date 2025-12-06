package main

import (
	"fmt"
	"net"
)

var address = "127.0.0.1:6699"
var listener *net.TCPListener //监听
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if nil != err {
		ColorPrintf(Red, "ResolveTCPAddr:%v err:%v\n", tcpAddr, err)
		return
	}
	listener, err = net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		ColorPrintf(Red, "ResolveTCPAddr:%v err:%v\n", tcpAddr, err)
		return
	}
	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if nil != err {
				ColorPrintf(Red, "listen.AcceptTCP, err:%v\n", err)
				return
			}
			_ = conn
			ColorPrintf(Green, "listen.AcceptTCP:%v\n", conn.RemoteAddr())
		}
	}()
	select {}
}

// 定义颜色代码
const (
	Reset  = "\033[0m"  // 重置
	Red    = "\033[31m" // 红色
	Green  = "\033[32m" // 绿色
	Yellow = "\033[33m" // 黄色
	Blue   = "\033[34m" // 蓝色
	Purple = "\033[35m" // 紫色
	Cyan   = "\033[36m" // 青色
	White  = "\033[37m" // 白色
)

// ColorPrintf 打印带颜色的格式化文本
func ColorPrintf(color string, format string, a ...interface{}) {
	fmt.Printf(color+format+Reset, a...)
}
