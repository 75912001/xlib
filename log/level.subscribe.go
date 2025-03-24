package log

type levelSubscribe struct {
	subMap       map[uint32]struct{}
	callBackFunc CallBackFunc
}

func newLevelSubscribe() *levelSubscribe {
	return &levelSubscribe{
		subMap: make(map[uint32]struct{}),
	}
}

// CallBackFunc 回调函数
type CallBackFunc func(level uint32, outString string)

// 是否订阅
func (p *levelSubscribe) isSubscribe(level uint32) bool {
	_, ok := p.subMap[level]
	return ok
}
