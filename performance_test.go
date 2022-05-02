package HighPerformanceMap

import (
	"container/list"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

// 测试map的value用unsafe.Pointer与interface比较
func BenchmarkUnsafePointAndInterfaceA(b *testing.B) {
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData := make(map[int]unsafe.Pointer)
		for i := 0; i < 10000; i++ {
			mapData[i] = unsafe.Pointer(&i)
		}
	}
	b.StopTimer()
}

func BenchmarkUnsafePointAndInterfaceB(b *testing.B) {
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData := make(map[int]interface{})
		for i := 0; i < 10000; i++ {
			mapData[i] = i
		}
	}
	b.StopTimer()
}

// 测试map key用string，key，interface比较
func BenchmarkIntStringInterfaceA(b *testing.B) {
	sliceList := make([]int, 0, 100000)
	for i := 0; i < cap(sliceList); i++ {
		sliceList = append(sliceList, i)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData := make(map[int]int)
		for i := 0; i < len(sliceList); i++ {
			mapData[sliceList[i]] = i
		}
	}
	b.StopTimer()
}

func BenchmarkIntStringInterfaceB(b *testing.B) {
	sliceList := make([]string, 0, 100000)
	for i := 0; i < cap(sliceList); i++ {
		sliceList = append(sliceList, strconv.Itoa(i))
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData := make(map[string]int)
		for i := 0; i < len(sliceList); i++ {
			mapData[sliceList[i]] = i
		}
	}
	b.StopTimer()
}

func BenchmarkIntStringInterfaceC(b *testing.B) {
	sliceList := make([]interface{}, 0, 100000)
	for i := 0; i < cap(sliceList); i++ {
		sliceList = append(sliceList, i)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData := make(map[interface{}]int)
		for i := 0; i < len(sliceList); i++ {
			mapData[sliceList[i]] = i
		}
	}
	b.StopTimer()
}

// list与slice比较
func BenchmarkListAndSliceInsertA(b *testing.B) {
	b.ResetTimer()
	sliceList := make([]int, 0)

	for j := 0; j < b.N; j++ {
		for i := 0; i < 10000; i++ {
			sliceList = append(sliceList, i)
		}
	}

	b.StopTimer()
}

func BenchmarkListAndSliceInsertB(b *testing.B) {
	b.ResetTimer()
	containerList := list.New()
	for j := 0; j < b.N; j++ {
		for i := 0; i < 10000; i++ {
			containerList.PushFront(i)
		}
	}
	b.StopTimer()
}

func BenchmarkListAndSliceRemoveA(b *testing.B) {
	num := 10000000
	sliceList := make([]int, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, i)
	}

	b.ResetTimer()
	for i := 0; i < num; i++ {
		sliceList = sliceList[1:len(sliceList)]
	}
	b.StopTimer()
}

func BenchmarkListAndSliceRemoveB(b *testing.B) {
	num := 10000000
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushFront(i)
	}

	b.ResetTimer()
	var prev *list.Element
	for e := containerList.Back(); e != nil; e = prev {
		prev = e.Prev()
		containerList.Remove(e)
	}
	b.StopTimer()
}

func BenchmarkListAndSliceRandomInsertA(b *testing.B) {
	num := 10
	sliceList := make([]int, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, i)
	}
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		l := len(sliceList)
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(l)
		newSlice := make([]int, l+1)
		copy(newSlice, sliceList[:randNum])
		newSlice[randNum] = randNum
		copy(newSlice[randNum+1:], sliceList[randNum:])
		sliceList = newSlice
	}

	b.StopTimer()
}

func BenchmarkListAndSliceRandomInsertB(b *testing.B) {
	num := 10
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushBack(i)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		l := containerList.Len()
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(l)

		for e := containerList.Front(); e != nil; e = e.Next() {
			if e.Value.(int) == randNum {
				ef := containerList.PushBack(randNum)
				containerList.MoveBefore(ef, e)
				break
			}
		}
	}
	b.StopTimer()
}

// list 循环删除方式
func TestSliceList(t *testing.T) {
	containerList := list.New()
	for i := 0; i < 10; i++ {
		containerList.PushFront(i)
	}

	var prev *list.Element
	for e := containerList.Back(); e != nil; e = prev {
		prev = e.Prev()
		containerList.Remove(e)
	}

	for e := containerList.Back(); e != nil; e = e.Prev() {
		t.Logf("%v", e.Value.(int))
	}

	t.Logf("len --> %v", containerList.Len())
}

// slice 随机添加
func TestSliceRandomInsert(t *testing.T) {
	num := 10
	sliceList := make([]int, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, i)
	}

	l := len(sliceList)
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(l)
	t.Logf("randNum --> %v", randNum)

	newSlice := make([]int, l+1)

	copy(newSlice, sliceList[:randNum])
	newSlice[randNum] = randNum
	copy(newSlice[randNum+1:], sliceList[randNum:])

	t.Logf("%v", newSlice)
}

// slice 随机删除
func TestListRandomInsert(t *testing.T) {
	num := 10
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushBack(i)
	}

	l := containerList.Len()
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(l)
	t.Logf("randNum --> %v", randNum)

	for e := containerList.Front(); e != nil; e = e.Next() {
		if e.Value.(int) == randNum {
			ef := containerList.PushBack(randNum)
			containerList.MoveBefore(ef, e)
			break
		}
	}

	for e := containerList.Front(); e != nil; e = e.Next() {
		t.Logf("%v", e.Value.(int))
	}
}
