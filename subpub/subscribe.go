// 订阅

package subpub

type ISubscribe interface {
	Subscribe(key uint64, onFunction func(...interface{}) error) error
}
