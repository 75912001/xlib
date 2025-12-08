package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

var address = "127.0.0.1:5602"

func main() {
	// 生成一个 1 秒链接一次的客户端
	cnt := 0
	var conns []*net.TCPConn // 创建conn slice
	sumCnt := 0
	failCnt := 0
	for {
		time.Sleep(1 * time.Second)
		for i := 0; i < 10; i++ {
			cnt++
			tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
			if nil != err {
				failCnt++
				ColorPrintf(Red, "ResolveTCPAddr:%v cnt:%v err:%v\n", tcpAddr, cnt, err)
				continue
			}
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if nil != err {
				failCnt++
				ColorPrintf(Red, "DialTCP:%v cnt:%v err:%v\n", tcpAddr, cnt, err)
				continue
			}
			_ = conn
			sumCnt++
			conns = append(conns, conn) // 将新连接添加到 slice 中
			fmt.Printf("DialTCP:%v cnt:%v success\n", tcpAddr, cnt)
		}
		// 将 conns 乱序
		rand.Shuffle(len(conns), func(i, j int) {
			conns[i], conns[j] = conns[j], conns[i]
		})
		// 每次循环完后，清理多余的连接，只保留两个
		if len(conns) > 2 {
			// 关闭多余的连接
			for i := 2; i < len(conns); i++ {
				conns[i].Close()
			}
			// 只保留前两个连接
			conns = conns[:2]
		}
		fmt.Printf("Cleaned up connections, remaining: %d sum:%v fail:%v\n",
			len(conns), sumCnt, failCnt)
	}
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
