package timer

import (
	xcontrol "github.com/75912001/xlib/control"
)

// Millisecond 毫秒级定时器
type Millisecond struct {
	xcontrol.ISwitchButton       // 有效(false:不执行,扫描时自动删除)
	xcontrol.ICallBack           // 到期-回调函数
	expire                 int64 // 过期时间
	xcontrol.IOut                // 到期-输出
}

// Delete 删除毫秒级定时器
func (p *Millisecond) Delete() {
	p.ISwitchButton.Off()
}

func (p *Millisecond) GetExpire() int64 {
	return p.expire
}
