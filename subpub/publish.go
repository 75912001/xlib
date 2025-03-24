// 发布

package subpub

type IPublish interface {
	Publish(key uint64, parameters ...interface{}) error
}
