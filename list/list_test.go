package list

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestNoPtrList(t *testing.T) {
	for i := 0; i < 33; i++ {
		globalNoPtrList.PrevInsert(i, globalNoPtrList.Len())
	}
	for i := 100; i < 152; i++ {
		globalNoPtrList.NextInsert(i, globalNoPtrList.Len())
	}
	globalNoPtrList.ModifyValue(999, 3)
	head := globalNoPtrList.Head()
	globalNoPtrList.Del(32)
	globalNoPtrList.Del(head)
	globalNoPtrList.Del(globalNoPtrList.Tail())
	globalNoPtrList.Del(3)
	head = globalNoPtrList.Head()
	length := globalNoPtrList.Len()
	for i := int64(0); i < length; i++ {
		fmt.Println(globalNoPtrList.GetValue(head).(int))
		head = globalNoPtrList.Next(head)
	}
}

type srcList struct {
	val  int
	next *srcList
}

const passLength = 10000000

func createNoPtrList() *NoPtrList {
	var list = NewNoPtrList(100000, func(length int64) IfaceSlice {
		des := make([]int, 100000)
		return listSlice(des)
	})
	for i := 0; i < passLength; i++ {
		list.PrevInsert(i, list.Len())
	}
	return list
}

// 创建passLength 的noptr链表与gc时间
func TestGcNoPtrList(t *testing.T) {
	now := time.Now()
	cur := createNoPtrList()
	fmt.Printf("NoPtrList addr %v create %d node time: %v\n", &cur, cur.Len(), time.Since(now))
	now = time.Now()
	runtime.GC()
	fmt.Printf("NoPtrList gc mark time: %v\n", time.Since(now))
}

func createSrcList() *srcList {
	list := new(srcList)
	head := list
	for i := 0; i < passLength; i++ {
		head.next = new(srcList)
		head.next.val = i
		head = head.next
	}
	return list
}

// 创建passLength 的原生链表与gc时间
func TestSrcList(t *testing.T) {
	now := time.Now()
	cur := createSrcList()
	fmt.Printf("src list addr %v create %d node time: %v\n", &cur, passLength, time.Since(now))
	now = time.Now()
	runtime.GC()
	fmt.Printf("src list gc mark time: %v\n", time.Since(now))
}

// 遍历链表，可以两个协程分别从首尾开始遍历
func TestRangeList(t *testing.T) {
	for i := 0; i < 100; i++ {
		globalNoPtrList.NextInsert(i, globalNoPtrList.Len())
	}
	length := globalNoPtrList.Len()
	head := globalNoPtrList.Head()
	for i := int64(0); i < length; i++ {
		fmt.Println(globalNoPtrList.GetValue(head).(int))
		head = globalNoPtrList.Next(head)
	}
}

var globalNoPtrList = NewNoPtrList(100000, func(length int64) IfaceSlice {
	des := make([]int, 100000)
	return listSlice(des)
})

type listSlice []int

func (des listSlice) Get(idx int64) interface{} {
	if idx < 0 || idx >= int64(len(des)) {
		return nil
	}
	return des[idx]
}

func (des listSlice) Set(data interface{}, idx int64) {
	if idx < 0 || idx >= int64(len(des)) {
		return
	}
	if de, ok := data.(int); ok {
		des[idx] = de
	}
}

func (des listSlice) Empty(idx int64) {
	if idx < 0 || idx >= int64(len(des)) {
		return
	}
	des[idx] = 0
}
