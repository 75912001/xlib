// 发布

package subpub

type IPublish[KEY ISubPubKey] interface {
	Publish(key KEY, args ...any) error // 发布
}
