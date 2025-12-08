package timer

import "container/list"

// IExpire 过期接口
type IExpire interface {
	GetExpire() int64 // 获取过期时间
}

// 按照过期时间排序
//
//	移动最后一个元素到合适的位置,移动到大于他的元素的前面[实现按照时间排序,加入顺序排序]
//	e.g.: 1,2,2,3,4,4,3 => 1,2,2,3,3,4,4 [将最后一个元素移动到4的前面]
func sortByExpire(l *list.List) {
	lastElement := l.Back() // 获取最后一个元素
	target := lastElement.Value.(IExpire)
	var element *list.Element
	for element = lastElement.Prev(); element != nil; element = element.Prev() {
		current := element.Value.(IExpire)
		if current.GetExpire() <= target.GetExpire() {
			l.MoveAfter(lastElement, element)
			return
		}
	}
	// 如果没有找到比目标小或等于的元素，将目标元素移动到列表的前面
	l.MoveToFront(lastElement)
}
