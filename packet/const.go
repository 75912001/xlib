package packet

import (
	"encoding/binary"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// 字节序
type endianMode uint32

const (
	LittleEndian endianMode = 0 // 小端-模式
	BigEndian    endianMode = 1 // 大端-模式
)

var gEndianMode = LittleEndian
var GEndian binary.ByteOrder

func init() {
	GEndian = binary.LittleEndian
}

func SetEndianMode(mode endianMode) {
	gEndianMode = mode
	switch gEndianMode {
	case LittleEndian:
		GEndian = binary.LittleEndian
	case BigEndian:
		GEndian = binary.BigEndian
	}
}

// IsBigEndian 是否为大端模式
func IsBigEndian() bool {
	return gEndianMode == BigEndian
}

// IsLittleEndian 是否为小端模式
func IsLittleEndian() bool {
	return gEndianMode == LittleEndian
}

// AddPacketToData 将数据 packet 放到 data 中
func AddPacketToData(data []byte, packet IPacket) ([]byte, error) {
	packetData, err := packet.Marshal()
	if err != nil {
		return nil, errors.WithMessagef(err, "AddPacketToData packet marshal %v, %v", packet, xruntime.Location())
	}
	if len(data) == 0 { //当 data len == 0 时候, 直接发送 v.data 数据...
		data = packetData
	} else {
		data = append(data, packetData...)
	}
	return data, nil
}
