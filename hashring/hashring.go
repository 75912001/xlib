package hashring

import (
	"fmt"
	"maps"
	"math"
	"sort"
)

// HashRing 哈希环
type HashRing[NODE comparable] struct {
	nodes                     []NODE                    // 物理节点列表
	ring                      map[uint32]NODE           // key: 虚拟节点哈希值 val:物理节点
	virtualNodeHashSliceOrder VirtualNodeHashSliceOrder // 所有虚拟节点的哈希点, 已排序, 便于二分查找
	weights                   map[NODE]uint32           // key: 物理节点 val: 物理节点权重. 决定虚拟节点数量
	virtualNodeCount          uint32                    // 虚拟节点基数
}

func NewHashRing[NODE comparable]() *HashRing[NODE] {
	hashRing := &HashRing[NODE]{
		nodes:                     make([]NODE, 0),
		ring:                      make(map[uint32]NODE),
		virtualNodeHashSliceOrder: make(VirtualNodeHashSliceOrder, 0),
		weights:                   make(map[NODE]uint32),
		virtualNodeCount:          DefaultVirtualNodeCount,
	}
	return hashRing
}

func (p *HashRing[NODE]) AddNode(node NODE) *HashRing[NODE] {
	return p.AddNodeWithWeight(node, 1)
}

func (p *HashRing[NODE]) AddNodeWithWeight(node NODE, weight uint32) *HashRing[NODE] {
	if p.IsNodeExist(node) { // 已存在
		return p
	}

	// 节点
	newNodes := make([]NODE, len(p.nodes), len(p.nodes)+1)
	copy(newNodes, p.nodes)
	newNodes = append(newNodes, node)

	// 权重
	newWeights := make(map[NODE]uint32)
	maps.Copy(newWeights, p.weights)
	newWeights[node] = weight

	hashRing := &HashRing[NODE]{
		nodes:                     newNodes,
		ring:                      make(map[uint32]NODE),
		virtualNodeHashSliceOrder: make(VirtualNodeHashSliceOrder, 0),
		weights:                   newWeights,
		virtualNodeCount:          p.virtualNodeCount,
	}
	hashRing.generateCircle()
	return hashRing
}

func (p *HashRing[NODE]) RemoveNode(node NODE) *HashRing[NODE] {
	if !p.IsNodeExist(node) { // 不存在
		return p
	}

	// 节点
	newNodes := make([]NODE, 0)
	for _, v := range p.nodes {
		if v != node {
			newNodes = append(newNodes, v)
		}
	}

	// 权重
	newWeights := make(map[NODE]uint32)
	for k, v := range p.weights {
		if k != node {
			newWeights[k] = v
		}
	}

	hashRing := &HashRing[NODE]{
		nodes:                     newNodes,
		ring:                      make(map[uint32]NODE),
		virtualNodeHashSliceOrder: make(VirtualNodeHashSliceOrder, 0),
		weights:                   newWeights,
		virtualNodeCount:          p.virtualNodeCount,
	}
	hashRing.generateCircle()
	return hashRing
}

func (p *HashRing[NODE]) GetNode(key string) (node NODE, ok bool) {
	pos, ok := p.getNodePos(key)
	if !ok {
		return node, false
	}
	return p.ring[p.virtualNodeHashSliceOrder[pos]], true
}

// 判断物理节点是否存在
func (p *HashRing[NODE]) IsNodeExist(node NODE) bool {
	_, ok := p.weights[node]
	return ok
}

func (p *HashRing[NODE]) getNodePos(key string) (pos int, ok bool) {
	if len(p.ring) == 0 {
		return 0, false
	}

	virtualNodeHash := p.genVirtualNodeHash(key)

	nodes := p.virtualNodeHashSliceOrder
	pos = sort.Search(len(nodes), func(i int) bool { return nodes[i] >= virtualNodeHash })

	if pos == len(nodes) { // 超出范围
		// 返回第一个节点
		return 0, true
	} else {
		return pos, true
	}
}

// 生成哈希环
func (p *HashRing[NODE]) generateCircle() {
	// 总权重
	totalWeight := uint32(0)
	for _, v := range p.weights {
		totalWeight += v
	}

	nodeCount := uint32(len(p.nodes)) // 节点数量
	for _, node := range p.nodes {
		weight := uint32(1)
		if v, ok := p.weights[node]; ok {
			weight = v
		}
		// 计算虚拟节点数量
		factor := math.Floor(float64(p.virtualNodeCount*nodeCount*weight) / float64(totalWeight))
		for i := 0; i < int(factor); i++ {
			virtualNode := fmt.Sprintf("%v.%v", node, i)
			// 16位的摘要
			virtualNodeHash16 := hashDigest(virtualNode)
			// 每个虚拟节点 生成 4段hash
			for j := 0; j < 4; j++ {
				virtualNodeHash := hashVal(virtualNodeHash16[j*4 : (j+1)*4])
				p.ring[virtualNodeHash] = node
				p.virtualNodeHashSliceOrder = append(p.virtualNodeHashSliceOrder, virtualNodeHash)
			}
		}
	}
	sort.Sort(&p.virtualNodeHashSliceOrder)
}

func (p *HashRing[NODE]) genVirtualNodeHash(data string) uint32 {
	k := hashDigest(data)
	return hashVal(k[0:4])
}
