// 发布

package subpub

type IPublish interface {
	Publish(key uint64, args ...any) error // 发布
}
