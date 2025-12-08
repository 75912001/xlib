package hashring

var DefaultVirtualNodeCount uint32 = 30 // 默认-虚拟节点基数

type VirtualNodeHashSliceOrder []uint32 // 虚拟节点-哈希值切片, 已排序, 便于二分查找

func (p *VirtualNodeHashSliceOrder) Len() int {
	return len(*p)
}

func (p *VirtualNodeHashSliceOrder) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *VirtualNodeHashSliceOrder) Less(i, j int) bool {
	return (*p)[i] < (*p)[j]
}
