package packet

// IPacket 接口-数据包
type IPacket interface {
	// Marshal 序列化
	Marshal() (data []byte, err error)
}
