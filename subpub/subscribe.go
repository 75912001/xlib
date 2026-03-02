// 订阅

package subpub

import xcontrol "github.com/75912001/xlib/control"

type ISubscribe[KEY ISubPubKey] interface {
	Subscribe(key KEY, onFunction xcontrol.OnFunction) error   // 订阅
	Unsubscribe(key KEY, onFunction xcontrol.OnFunction) error // 取消订阅
}
