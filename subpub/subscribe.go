// 订阅

package subpub

type ISubscribe interface {
	Subscribe(key uint64, onFunction func(...any) error) error   // 订阅
	Unsubscribe(key uint64, onFunction func(...any) error) error // 取消订阅
}
