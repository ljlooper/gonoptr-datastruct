package list

type IfaceSlice interface {
	Get(idx int64) interface{}
	Set(data interface{}, idx int64)
	Empty(idx int64)
}

type NewIfaceSlicefunc func(length int64) IfaceSlice

type listNode struct {
	// 前驱
	prev int64
	// 后继
	next int64
	// 节点所在的位置
	idx int64
	// 值所在的位置
	valIdx int64
}

type linkList struct {
	// 存储链表节点的切片
	listSlice []listNode
	// 值切片，由使用者实现对应的接口, 自定义切片的类型使用值类型
	data IfaceSlice
	// 真实节点与值的个数
	realCount int64
}

type NoPtrList struct {
	// linkList个数
	size int64
	// 当前所在的索引 == 所有节点个数, 0号索引用来作标记
	index int64
	// 每个切片的长度
	length int64
	// 存储linkList的首地址
	addr []*linkList
	// 初始化值切片的方法，由使用者实现
	new NewIfaceSlicefunc
	// 头节点
	head *listNode
	// 尾节点
	tail *listNode
}

const (
	defineLength = 1 << 16
)

// 初始化链表，如果链表长度较短建议用原生链表
func NewNoPtrList(length int64, newSlice NewIfaceSlicefunc) *NoPtrList {
	if length <= 1 {
		length = defineLength
	}
	return &NoPtrList{
		new:    newSlice,
		length: length,
	}
}

// idx节点前插，如果链表是空则直接插入
func (list *NoPtrList) PrevInsert(val interface{}, idx int64) int64 {
	if !list.isInsert(idx) {
		return 0
	}
	list.isExpand()
	list.index++
	// 插入val
	insertSliceIdx, insertIdx := arithmeticMod(list.index, list.length)
	insertSlice := list.addr[insertSliceIdx]
	insertSlice.insertVal(val, insertIdx, list.index)
	insertSlice.realCount++
	// 第一个节点
	if list.index == 1 {
		list.firstNode(insertSlice, insertIdx)
		return list.index
	}
	// 获取idx对应的节点
	oldNode := list.getListNode(idx)
	// 头插
	if oldNode.prev == 0 {
		// 老节点的前去变为新节点的索引
		oldNode.prev = list.index
		// 之前的头节点赋值为新头节点的后继
		insertSlice.listSlice[insertIdx].next = oldNode.idx
		// 标记头节点
		insertSlice.listSlice[insertIdx].prev = 0
		list.head = &insertSlice.listSlice[insertIdx]
	} else {
		// 任意idx插入
		// 查找idx的前驱
		oldPrevNode := list.getListNode(oldNode.prev)
		// idx的前驱的后继赋值为新节点的索引
		oldPrevNode.next = list.index
		// 新节点的前驱赋值为idx前驱的索引
		insertSlice.listSlice[insertIdx].prev = oldPrevNode.idx
		// 新节点的后继赋值为idx节点的索引
		insertSlice.listSlice[insertIdx].next = oldNode.idx
		// idx的前驱赋值为新节点的索引
		oldNode.prev = list.index
	}
	return list.index
}

// idx节点后插，如果链表是空则直接插入
func (list *NoPtrList) NextInsert(val interface{}, idx int64) int64 {
	if !list.isInsert(idx) {
		return 0
	}
	list.isExpand()
	list.index++
	// 插入val
	insertSliceIdx, insertIdx := arithmeticMod(list.index, list.length)
	insertSlice := list.addr[insertSliceIdx]
	insertSlice.insertVal(val, insertIdx, list.index)
	insertSlice.realCount++
	// 如果是第一个节点
	if list.index == 1 {
		list.firstNode(insertSlice, insertIdx)
		return list.index
	}
	// 获取idx对应的节点
	oldNode := list.getListNode(idx)
	// 尾插
	if oldNode.next == 0 {
		// 之前的尾节点后继赋值为新节点的idx
		oldNode.next = list.index
		// 之前的尾节点的
		insertSlice.listSlice[insertIdx].prev = oldNode.idx
		// 标记尾节点
		insertSlice.listSlice[insertIdx].next = 0
		list.tail = &insertSlice.listSlice[insertIdx]
	} else {
		// 获取idx的后继
		oldNextNode := list.getListNode(oldNode.next)
		// 将idx的后继的前驱赋值为新节点的索引
		oldNextNode.prev = list.index
		// 新节点的前驱赋值为idx节点的索引
		insertSlice.listSlice[insertIdx].prev = oldNode.idx
		// 新节点的后继赋值为idx的后继的索引
		insertSlice.listSlice[insertIdx].next = oldNextNode.idx
		// idx节点的后继赋值为新节点的索引
		oldNode.next = list.index
	}
	return list.index
}

// 修改idx对应的值
func (list *NoPtrList) ModifyValue(val interface{}, idx int64) {
	if !list.isRight(idx) {
		return
	}
	linkListIdx, listSliceIdx := arithmeticMod(idx, list.length)
	valIdx := list.addr[linkListIdx].listSlice[listSliceIdx].valIdx
	list.addr[linkListIdx].data.Set(val, valIdx%list.length)
}

// 获取idx对应的值
func (list *NoPtrList) GetValue(idx int64) interface{} {
	if !list.isRight(idx) {
		return nil
	}
	linkListIdx, listSliceIdx := arithmeticMod(idx, list.length)
	valIdx := list.addr[linkListIdx].listSlice[listSliceIdx].valIdx
	return list.addr[linkListIdx].data.Get(valIdx % list.length)
}

// 删除
func (list *NoPtrList) Del(idx int64) {
	if !list.isRight(idx) {
		return
	}
	delNode := list.getListNode(idx)
	// 删除尾节点
	if delNode.next == 0 {
		delPrevNode := list.getListNode(delNode.prev)
		delPrevNode.next = 0
		list.tail = delPrevNode
	} else if delNode.prev == 0 {
		// 删除头节点
		delNextNode := list.getListNode(delNode.next)
		delNextNode.prev = 0
		list.head = delNextNode
	} else {
		delPrevNode := list.getListNode(delNode.prev)
		delNextNode := list.getListNode(delNode.next)
		delPrevNode.next = delNextNode.idx
		delNextNode.prev = delPrevNode.idx
	}
	// 移动最后一个节点到idx位置上，并且赋值
	list.moveLastToIdx(idx)
	list.index--
}

// 移动最后一个元素到某个位置上
func (list *NoPtrList) moveLastToIdx(idx int64) {
	lastSliceIdx, lastDataIdx := arithmeticMod(list.index, list.length)
	// 最后一个区间链表元素个数减1
	list.addr[lastSliceIdx].realCount--
	// 如果idx是最后一个元素不需要移动，置为空即可
	if idx == list.index {
		list.addr[lastSliceIdx].data.Empty(lastDataIdx)
		list.addr[lastSliceIdx].listSlice[lastDataIdx] = listNode{}
		return
	}
	// 将最后一个val移动到idx位置上，将最后一个val置为空
	lastVal := list.addr[lastSliceIdx].data.Get(lastDataIdx)
	idxSliceIdx, idxDataIdx := arithmeticMod(idx, list.length)
	list.addr[idxSliceIdx].data.Set(lastVal, idxDataIdx)
	list.addr[lastSliceIdx].data.Empty(lastDataIdx)
	// 移动节点位置
	lastNode := list.addr[lastSliceIdx].listSlice[lastDataIdx]
	// 修改lastNode的前驱
	if lastNode.prev != 0 {
		lastPrevNode := list.getListNode(lastNode.prev)
		lastPrevNode.next = idx
	}
	// 修改lastNode的后继
	if lastNode.next != 0 {
		lastNextNode := list.getListNode(lastNode.next)
		lastNextNode.prev = idx
	}
	// 更改本身的idx和valIdx
	lastNode.idx = idx
	lastNode.valIdx = idx
	// 移动到idx位置
	list.addr[idxSliceIdx].listSlice[idxDataIdx] = lastNode
	// 将最后一个位置置为空
	list.addr[lastSliceIdx].listSlice[lastDataIdx] = listNode{}
}

// 获取后继节点索引
func (list *NoPtrList) Next(idx int64) int64 {
	if !list.isRight(idx) {
		return 0
	}
	linkListIdx, listSliceIdx := arithmeticMod(idx, list.length)
	return list.addr[linkListIdx].listSlice[listSliceIdx].next
}

// 获取前驱节点索引
func (list *NoPtrList) Prev(idx int64) int64 {
	if !list.isRight(idx) {
		return 0
	}
	linkListIdx, listSliceIdx := arithmeticMod(idx, list.length)
	return list.addr[linkListIdx].listSlice[listSliceIdx].prev
}

// 获取头节点索引
func (list *NoPtrList) Head() int64 {
	if list.head == nil {
		return 0
	}
	return list.head.idx
}

// 获取尾节点索引
func (list *NoPtrList) Tail() int64 {
	if list.tail == nil {
		return 0
	}
	return list.tail.idx
}

// 获取节点个数
func (list *NoPtrList) Len() int64 {
	return list.index
}

// 链表内部结构初始化与扩容
func (list *NoPtrList) isExpand() {
	if list.index != 0 && (list.index+1)%list.length != 0 {
		return
	}
	l := &linkList{
		listSlice: make([]listNode, list.length),
		data:      list.new(list.length),
	}
	// 初始化
	if list.addr == nil {
		list.addr = []*linkList{l}
	} else {
		// 扩容
		list.addr = append(list.addr, l)
	}
	list.size++
}

// 非插入操作判断idx是否有效
func (list *NoPtrList) isRight(idx int64) bool {
	return idx > 0 && idx <= list.index
}

// 插入操作判断idx是否有效
func (list *NoPtrList) isInsert(idx int64) bool {
	return idx >= 0 && idx <= list.index
}

// 链表插入的第一个节点
func (list *NoPtrList) firstNode(firstList *linkList, insertIdx int64) {
	firstList.listSlice[insertIdx].prev = 0
	firstList.listSlice[insertIdx].next = 0
	list.head = &firstList.listSlice[insertIdx]
	list.tail = &firstList.listSlice[insertIdx]
}

// 获取链表的节点
func (list *NoPtrList) getListNode(idx int64) *listNode {
	if list.head.idx == idx {
		return list.head
	} else if list.tail.idx == idx {
		return list.tail
	}
	// 查找idx对应的节点
	oldNodeSliceIdx, oldNodeIdx := arithmeticMod(idx, list.length)
	oldNodeSlice := list.addr[oldNodeSliceIdx]
	return &oldNodeSlice.listSlice[oldNodeIdx]
}

// 插入值
func (insertSlice *linkList) insertVal(val interface{}, insertIdx, index int64) {
	// 插入val
	insertSlice.data.Set(val, insertIdx)
	// 给节点赋值value的idx
	insertSlice.listSlice[insertIdx].valIdx = index
	// 给节点赋值自己所在当前区间链表的索引
	insertSlice.listSlice[insertIdx].idx = index
}

// 计算所在的区间链表，和节点所在的位置
func arithmeticMod(idx, length int64) (int64, int64) {
	sliceIndex := idx / length
	dataIndex := idx % length
	return sliceIndex, dataIndex
}
