package HighPerformanceMap

import (
	"container/list"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

// test map value ues unsafe.Pointer and interface type compare
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

// test map key use int、int64、string and interface type compare
func BenchmarkIntInt64StringInterfaceA(b *testing.B) {
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

func BenchmarkIntInt64StringInterfaceB(b *testing.B) {
	sliceList := make([]int64, 0, 100000)
	for i := 0; i < cap(sliceList); i++ {
		sliceList = append(sliceList, int64(i))
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData := make(map[int64]int)
		for i := 0; i < len(sliceList); i++ {
			mapData[sliceList[i]] = i
		}
	}
	b.StopTimer()
}

func BenchmarkIntInt64StringInterfaceC(b *testing.B) {
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

func BenchmarkIntInt64StringInterfaceD(b *testing.B) {
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

// list and slice add compare
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

func TestListAndSliceInsertA(t *testing.T) {
	sliceList := make([]int, 0)

	for i := 0; i < 10; i++ {
		sliceList = append(sliceList, i)
	}
	t.Logf("%v", sliceList)
}

func BenchmarkListAndSliceInsertB(b *testing.B) {
	b.ResetTimer()
	containerList := list.New()
	for j := 0; j < b.N; j++ {
		for i := 0; i < 10000; i++ {
			containerList.PushBack(i)
		}
	}
	b.StopTimer()
}

func TestListAndSliceInsertB(t *testing.T) {
	containerList := list.New()
	for i := 0; i < 10; i++ {
		containerList.PushBack(i)
	}
	for e := containerList.Front(); e != nil; e = e.Next() {
		t.Logf("%v", e.Value.(int))
	}
}

// list and slice remove compare
func BenchmarkListAndSliceRemoveA(b *testing.B) {
	num := 100
	sliceList := make([]int, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, i)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		a := sliceList[:]
		for i := 0; i < num; i++ {
			a = a[1:len(a)]
		}
	}
	b.StopTimer()
}

func TestListAndSliceRemoveA(t *testing.T) {
	num := 10
	sliceList := make([]int, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, i)
	}

	a := sliceList[:]
	for i := 0; i < num; i++ {
		a = a[1:len(a)]

		t.Logf("%v", a)
	}
}

func BenchmarkListAndSliceRemoveB(b *testing.B) {
	num := 100
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushFront(i)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		copyContainerList := list.New()
		copyContainerList.PushBackList(containerList)
		var prev *list.Element
		for e := copyContainerList.Back(); e != nil; e = prev {
			prev = e.Prev()
			copyContainerList.Remove(e)
		}
	}
	b.StopTimer()
}

func TestListAndSliceRemoveB(t *testing.T) {
	num := 10
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushFront(i)
	}

	copyContainerList := list.New()
	copyContainerList.PushBackList(containerList)
	var prev *list.Element
	for e := copyContainerList.Back(); e != nil; e = prev {
		t.Logf("%v", e.Value.(int))
		prev = e.Prev()
		copyContainerList.Remove(e)
	}
	t.Logf("len --> %v", copyContainerList.Len())
}

// If you know the index, list and slice random add compare
func BenchmarkListAndSliceRandomIndexInsertA(b *testing.B) {
	num := 10000
	sliceList := make([]int, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, i)
	}
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		l := len(sliceList)
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(l) // random index
		newSlice := make([]int, l+1)
		copy(newSlice, sliceList[:randNum])
		newSlice[randNum] = randNum
		copy(newSlice[randNum+1:], sliceList[randNum:])
		sliceList = newSlice
	}

	b.StopTimer()
}

func TestListAndSliceRandomIndexInsertA(t *testing.T) {
	num := 10
	sliceList := make([]int, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, i)
	}

	for j := 0; j < 2; j++ {
		l := len(sliceList)
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(l) // random index
		t.Logf("randNum --> %v", randNum)
		newSlice := make([]int, l+1)
		copy(newSlice, sliceList[:randNum])
		newSlice[randNum] = randNum
		copy(newSlice[randNum+1:], sliceList[randNum:])
		sliceList = newSlice
	}
	t.Logf("%v", sliceList)
}

func BenchmarkListAndSliceRandomIndexInsertB(b *testing.B) {
	num := 10000
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushBack(i)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		l := containerList.Len()
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(l)

		index := 0
		for e := containerList.Front(); e != nil; e = e.Next() {
			if index == randNum {
				ef := containerList.PushBack(randNum)
				containerList.MoveBefore(ef, e)
				break
			} else {
				index += 1
			}
		}
	}
	b.StopTimer()
}

func TestListAndSliceRandomIndexInsertB(t *testing.T) {
	num := 10
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushBack(i)
	}

	for j := 0; j < 2; j++ {
		l := containerList.Len()
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(l) // random index
		t.Logf("randNum --> %v", randNum)

		index := 0
		for e := containerList.Front(); e != nil; e = e.Next() {
			if index == randNum {
				ef := containerList.PushBack(randNum)
				containerList.MoveBefore(ef, e)
				break
			} else {
				index += 1
			}
		}
	}

	for e := containerList.Front(); e != nil; e = e.Next() {
		t.Logf("%v", e.Value.(int))
	}
}

// If you know the index, list and slice random remove compare
func BenchmarkListAndSliceRandomIndexRemoveA(b *testing.B) {
	num := 10000
	sliceList := make([]int, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, i)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		a := sliceList[:]

		for i := 0; i < num/2; i++ {
			l := len(a)
			rand.Seed(time.Now().UnixNano())
			randNum := rand.Intn(l) // random index
			newSlice := make([]int, l-1)
			copy(newSlice, a[:randNum])
			copy(newSlice[randNum:], a[randNum+1:])
			a = newSlice
		}

	}
	b.StopTimer()
}

func TestListAndSliceRandomIndexRemoveA(t *testing.T) {
	num := 10
	sliceList := make([]int, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, i)
	}

	a := sliceList[:]

	for i := 0; i < num/2; i++ {
		l := len(a)
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(l) // random index
		t.Logf("randNum --> %v", randNum)
		newSlice := make([]int, l-1)
		copy(newSlice, a[:randNum])
		copy(newSlice[randNum:], a[randNum+1:])
		a = newSlice
	}
	t.Logf("%v", a)
}

func BenchmarkListAndSliceRandomIndexRemoveB(b *testing.B) {
	num := 10000
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushBack(i)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		copyContainerList := list.New()
		copyContainerList.PushBackList(containerList)

		for i := 0; i < num/2; i++ {
			l := copyContainerList.Len()
			rand.Seed(time.Now().UnixNano())
			randNum := rand.Intn(l)

			index := 0
			for e := copyContainerList.Front(); e != nil; e = e.Next() {
				if index == randNum {
					copyContainerList.Remove(e)
					break
				} else {
					index += 1
				}
			}
		}
	}
	b.StopTimer()
}

func TestListAndSliceRandomIndexRemoveB(t *testing.T) {
	num := 10
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushBack(i)
	}

	copyContainerList := list.New()
	copyContainerList.PushBackList(containerList)

	for i := 0; i < num/2; i++ {
		l := copyContainerList.Len()
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(l)

		t.Logf("randNum --> %v", randNum)

		index := 0
		for e := copyContainerList.Front(); e != nil; e = e.Next() {
			if index == randNum {
				copyContainerList.Remove(e)
				break
			} else {
				index += 1
			}
		}
	}

	for e := copyContainerList.Front(); e != nil; e = e.Next() {
		t.Logf("%v", e.Value.(int))
	}
}

// big size, slice and list compare
type IntBig struct {
	Num1 int64
	Num2 int64
	Num3 int64
	Num4 int64
	Num5 int64
}

func BenchmarkListAndSliceBigRandomIndexInsertA(b *testing.B) {
	num := 100000
	sliceList := make([]IntBig, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, IntBig{
			Num1: int64(i),
			Num2: int64(i),
			Num3: int64(i),
			Num4: int64(i),
			Num5: int64(i),
		})
	}
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		l := len(sliceList)
		rand.Seed(time.Now().UnixNano())
		randNum := int64(rand.Intn(l)) // random index
		newSlice := make([]IntBig, l+1)
		copy(newSlice, sliceList[:randNum])
		newSlice[randNum] = IntBig{
			Num1: randNum,
			Num2: randNum,
			Num3: randNum,
			Num4: randNum,
			Num5: randNum,
		}
		copy(newSlice[randNum+1:], sliceList[randNum:])
		sliceList = newSlice
	}

	b.StopTimer()
}

func TestListAndSliceBigRandomIndexInsertA(t *testing.T) {
	num := 10
	sliceList := make([]IntBig, 0)
	for i := 0; i < num; i++ {
		sliceList = append(sliceList, IntBig{
			Num1: int64(i),
			Num2: int64(i),
			Num3: int64(i),
			Num4: int64(i),
			Num5: int64(i),
		})
	}

	for i := 0; i < num/2; i++ {
		l := len(sliceList)
		rand.Seed(time.Now().UnixNano())
		randNum := int64(rand.Intn(l)) // random index
		t.Logf("randNum --> %v", randNum)
		newSlice := make([]IntBig, l+1)
		copy(newSlice, sliceList[:randNum])
		newSlice[randNum] = IntBig{
			Num1: randNum,
			Num2: randNum,
			Num3: randNum,
			Num4: randNum,
			Num5: randNum,
		}
		copy(newSlice[randNum+1:], sliceList[randNum:])
		sliceList = newSlice
	}
	t.Logf("%v", sliceList)
}

func BenchmarkListAndSliceBigRandomIndexInsertB(b *testing.B) {
	num := 100000
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushBack(IntBig{
			Num1: int64(i),
			Num2: int64(i),
			Num3: int64(i),
			Num4: int64(i),
			Num5: int64(i),
		})
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		l := containerList.Len()
		rand.Seed(time.Now().UnixNano())
		randNum := int64(rand.Intn(l)) // random index

		index := int64(0)
		for e := containerList.Front(); e != nil; e = e.Next() {
			if index == randNum {
				ef := containerList.PushBack(IntBig{
					Num1: randNum,
					Num2: randNum,
					Num3: randNum,
					Num4: randNum,
					Num5: randNum,
				})
				containerList.MoveBefore(ef, e)
				break
			} else {
				index += 1
			}
		}
	}
	b.StopTimer()
}

func TestListAndSliceBigRandomIndexInsertB(t *testing.T) {
	num := 10
	containerList := list.New()
	for i := 0; i < num; i++ {
		containerList.PushBack(IntBig{
			Num1: int64(i),
			Num2: int64(i),
			Num3: int64(i),
			Num4: int64(i),
			Num5: int64(i),
		})
	}

	for j := 0; j < 2; j++ {
		l := containerList.Len()
		rand.Seed(time.Now().UnixNano())
		randNum := int64(rand.Intn(l)) // random index
		t.Logf("randNum --> %v", randNum)

		index := int64(0)
		for e := containerList.Front(); e != nil; e = e.Next() {
			if index == randNum {
				ef := containerList.PushBack(IntBig{
					Num1: randNum,
					Num2: randNum,
					Num3: randNum,
					Num4: randNum,
					Num5: randNum,
				})
				containerList.MoveBefore(ef, e)
				break
			} else {
				index += 1
			}
		}
	}

	for e := containerList.Front(); e != nil; e = e.Next() {
		t.Logf("%v", e.Value)
	}
}

func BenchmarkListAndSliceBigInsertA(b *testing.B) {
	b.ResetTimer()
	sliceList := make([]IntBig, 0)

	for j := 0; j < b.N; j++ {
		for i := 0; i < 1000; i++ {
			sliceList = append(sliceList, IntBig{
				Num1: int64(i),
				Num2: int64(i),
				Num3: int64(i),
				Num4: int64(i),
				Num5: int64(i),
			})
		}
	}

	b.StopTimer()
}

func BenchmarkListAndSliceBigInsertB(b *testing.B) {
	b.ResetTimer()
	containerList := list.New()
	for j := 0; j < b.N; j++ {
		for i := 0; i < 1000; i++ {
			containerList.PushBack(IntBig{
				Num1: int64(i),
				Num2: int64(i),
				Num3: int64(i),
				Num4: int64(i),
				Num5: int64(i),
			})
		}
	}
	b.StopTimer()
}

// list for remove element
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
