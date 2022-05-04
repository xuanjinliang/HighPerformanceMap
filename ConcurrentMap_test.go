package HighPerformanceMap

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

// A functional test
func TestCreateConcurrentSliceMapString(t *testing.T) {
	mapData := CreateConcurrentSliceMap(99)
	v, ok := mapData.Get(StrKey("Hello"))
	if ok == false {
		t.Logf("Hello is not exist")
	}

	mapData.Set(StrKey("Hello"), 123)

	v, ok = mapData.Get(StrKey("Hello"))
	if v.(int) != 123 || ok != true {
		t.Error("set/get failed.")
	}
	mapData.Delete(StrKey("Hello"))

	v, ok = mapData.Get(StrKey("Hello"))
	if v != nil || ok != false {
		t.Error("del failed")
	}
}

func TestCreateConcurrentSliceMapInt64(t *testing.T) {
	mapData := CreateConcurrentSliceMap(99)
	v, ok := mapData.Get(I64Key(111))
	if ok == false {
		t.Logf("111 is not exist")
	}

	mapData.Set(I64Key(111), "jinjin")

	v, ok = mapData.Get(I64Key(111))
	if v.(string) != "jinjin" || ok != true {
		t.Error("set/get failed.")
	}
	mapData.Delete(I64Key(111))

	v, ok = mapData.Get(I64Key(111))
	if v != nil || ok != false {
		t.Error("del failed")
	}
}

func TestCreateConcurrentSliceMapStringLen(t *testing.T) {
	num := 1000

	sliceList := make([]string, num)
	for i := 0; i < num; i++ {
		sliceList[i] = strconv.Itoa(i)
	}

	mapData := CreateConcurrentSliceMap(99)

	for _, data := range sliceList {
		mapData.Set(StrKey(data), data)
	}

	t.Logf("%v", mapData.Len())
}

func TestCreateConcurrentSliceMapStringRange(t *testing.T) {
	mapData := CreateConcurrentSliceMap(99)

	mapData.Set(StrKey("Hello"), 123)
	mapData.Set(StrKey("Hello World"), 123)
	mapData.Set(StrKey("Hello World1"), 123)

	mapData.Range(func(key, value interface{}) bool {
		t.Logf("key --> %v, value --> %v \n", key.(string), value.(int))
		if key.(string) == "Hello" {
			return false
		}
		return true
	})
}

// goroutine Test
func TestGoroutineSet(t *testing.T) {
	num := 10000

	sliceList := make([]string, num)
	for i := 0; i < num; i++ {
		sliceList[i] = strconv.Itoa(i)
	}

	mapData := CreateConcurrentSliceMap(99)
	goroutineNum := 100
	ch := make(chan string, goroutineNum)
	wg := sync.WaitGroup{}

	for i := 0; i < goroutineNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if data, ok := <-ch; ok {
					mapData.Set(StrKey(data), data)
				} else {
					break
				}
			}
		}()
	}

	for _, v := range sliceList {
		ch <- v
	}

	close(ch)
	wg.Wait()

	t.Logf("%v", mapData.Len())
}

func TestGoroutineRead(t *testing.T) {
	num := 10000

	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		mapData.Set(StrKey(strconv.Itoa(i)), i)
	}

	goroutineNum := 100
	wg := sync.WaitGroup{}

	for i := 0; i < goroutineNum; i++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			if data, ok := mapData.Get(StrKey(key)); ok {
				t.Logf("key --> %v, value --> %v \n", key, data)
			} else {
				t.Errorf("get key %v Error", key)
			}

		}(strconv.Itoa(i))
	}

	wg.Wait()
}

func TestGoroutineDelete(t *testing.T) {
	num := 10000

	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		mapData.Set(StrKey(strconv.Itoa(i)), i)
	}

	goroutineNum := 100
	wg := sync.WaitGroup{}

	for i := 0; i < goroutineNum; i++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			mapData.Delete(StrKey(key))
			if data, ok := mapData.Get(StrKey(key)); ok {
				t.Errorf("key --> %v, value --> %v \n", key, data)
			}
		}(strconv.Itoa(i))
	}

	wg.Wait()
	t.Logf("len --> %v", mapData.Len())
	t.Logf("free len --> %v", mapData.FreeLen())
}

// performance Test write
func BenchmarkSyncAndMapAndPMapSetA(b *testing.B) {
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData := make(map[string]int)
		num := 100
		mu := sync.RWMutex{}
		wg := sync.WaitGroup{}

		for i := 0; i < num; i++ {
			wg.Add(1)
			go func(index int) {
				mu.Lock()
				defer func() {
					mu.Unlock()
					wg.Done()
				}()
				mapData[strconv.Itoa(index)] = index
			}(i)
		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapSetB(b *testing.B) {
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData := sync.Map{}
		num := 100
		wg := sync.WaitGroup{}
		for i := 0; i < num; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				mapData.Store(strconv.Itoa(index), index)
			}(i)

		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapSetC(b *testing.B) {
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData := CreateConcurrentSliceMap(99)
		num := 100
		wg := sync.WaitGroup{}
		for i := 0; i < num; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				mapData.Set(StrKey(strconv.Itoa(index)), index)
			}(i)
		}
		wg.Wait()
	}
	b.StopTimer()
}

// performance Test Read
func BenchmarkSyncAndMapAndPMapReadA(b *testing.B) {
	num := 1000000
	mapData := make(map[string]int)
	for i := 0; i < num; i++ {
		mapData[strconv.Itoa(i)] = i
	}

	randNum := 100
	mu := sync.RWMutex{}

	b.ResetTimer()
	wg := sync.WaitGroup{}
	for j := 0; j < b.N; j++ {
		wg.Add(1)
		go func(index int) {
			mu.RLock()
			defer func() {
				mu.RUnlock()
				wg.Done()
			}()
			if _, ok := mapData[strconv.Itoa(index+randNum)]; ok {

			} else {
				b.Errorf("error %v", index+randNum)
			}
		}(j)
	}
	wg.Wait()
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapReadB(b *testing.B) {
	num := 1000000
	mapData := sync.Map{}
	for i := 0; i < num; i++ {
		mapData.Store(strconv.Itoa(i), i)
	}

	randNum := 100

	b.ResetTimer()
	wg := sync.WaitGroup{}
	for j := 0; j < b.N; j++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			if _, ok := mapData.Load(strconv.Itoa(index + randNum)); ok {

			} else {
				b.Errorf("error %v", index+randNum)
			}
		}(j)

	}
	wg.Wait()
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapReadC(b *testing.B) {
	num := 1000000
	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		mapData.Set(StrKey(strconv.Itoa(i)), i)
	}

	randNum := 100

	b.ResetTimer()
	wg := sync.WaitGroup{}
	for j := 0; j < b.N; j++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			if _, ok := mapData.Get(StrKey(strconv.Itoa(index + randNum))); ok {

			} else {
				b.Errorf("error %v", index+randNum)
			}
		}(j)
	}
	wg.Wait()
	b.StopTimer()
}

// performance Test GC recycle
func BenchmarkSyncAndMapAndPMapGCA(b *testing.B) {
	num := 1000000
	mapData := make(map[string]int)
	for i := 0; i < num; i++ {
		mapData[strconv.Itoa(i)] = i
	}

	randNum := 100
	wg := sync.WaitGroup{}
	mu := sync.RWMutex{}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		wg.Add(1)
		go func(index int) {
			mu.Lock()
			defer func() {
				mu.Unlock()
				wg.Done()
			}()
			delete(mapData, strconv.Itoa(index+randNum))
		}(j)
		runtime.GC()
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapGCB(b *testing.B) {
	num := 1000000
	mapData := sync.Map{}
	for i := 0; i < num; i++ {
		mapData.Store(strconv.Itoa(i), i)
	}

	randNum := 100

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData.Delete(strconv.Itoa(j + randNum))
		runtime.GC()
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapGCC(b *testing.B) {
	num := 1000000
	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		mapData.Set(StrKey(strconv.Itoa(i)), i)
	}

	randNum := 100

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData.Delete(StrKey(strconv.Itoa(j + randNum)))
		runtime.GC()
	}
	b.StopTimer()
}

// big Data
type intBig struct {
	Num1 int64
	Num2 int64
	Num3 int64
	Num4 int64
	Num5 int64
}

// performance Test write
func BenchmarkSyncAndMapAndPMapBigSetA(b *testing.B) {
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		num := 1000
		mapData := make(map[string]*intBig)
		mu := sync.Mutex{}
		for i := 0; i < num; i++ {
			mu.Lock()
			iBig := int64(i)
			d := &intBig{
				iBig, iBig, iBig, iBig, iBig,
			}
			mapData[strconv.Itoa(i)] = d
			mu.Unlock()
		}
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapBigSetB(b *testing.B) {
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		num := 1000
		mapData := sync.Map{}
		for i := 0; i < num; i++ {
			iBig := int64(i)
			d := &intBig{
				iBig, iBig, iBig, iBig, iBig,
			}
			mapData.Store(strconv.Itoa(i), d)
		}
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapBigSetC(b *testing.B) {
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		num := 1000
		mapData := CreateConcurrentSliceMap(99)
		for i := 0; i < num; i++ {
			iBig := int64(i)
			d := &intBig{
				iBig, iBig, iBig, iBig, iBig,
			}
			mapData.Set(StrKey(strconv.Itoa(i)), d)
		}
	}
	b.StopTimer()
}

// performance Test Read
func BenchmarkSyncAndMapAndPMapBigReadA(b *testing.B) {
	num := 1000000
	mapData := make(map[string]*intBig)
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData[strconv.Itoa(i)] = d
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mu := sync.RWMutex{}
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(num)
		mu.RLock()
		d := mapData[strconv.Itoa(randNum)]
		mu.RUnlock()
		fmt.Sprintf("%v", d)
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapBigReadB(b *testing.B) {
	num := 1000000
	mapData := sync.Map{}
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData.Store(strconv.Itoa(i), d)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(num)
		if d, ok := mapData.Load(strconv.Itoa(randNum)); ok {
			fmt.Sprintf("%v", d)
		} else {
			b.Errorf("error %v", randNum)
		}
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapBigReadC(b *testing.B) {
	num := 1000000
	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData.Set(StrKey(strconv.Itoa(i)), d)
	}

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(num)
		if d, ok := mapData.Get(StrKey(strconv.Itoa(randNum))); ok {
			fmt.Sprintf("%v", d)
		} else {
			b.Errorf("error %v", randNum)
		}
	}
	b.StopTimer()
}

// performance Test GC recycle
func BenchmarkSyncAndMapAndPMapBigGCA(b *testing.B) {
	num := 100000
	mapData := make(map[string]*intBig)
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData[strconv.Itoa(i)] = d
	}

	randNum := 100

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mu := sync.RWMutex{}
		mu.Lock()
		delete(mapData, strconv.Itoa(j+randNum))
		mu.Unlock()
		runtime.GC()
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapBigGCB(b *testing.B) {
	num := 100000
	mapData := sync.Map{}
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData.Store(strconv.Itoa(i), d)
	}

	randNum := 100

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData.Delete(strconv.Itoa(j + randNum))
		runtime.GC()
	}
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapBigGCC(b *testing.B) {
	num := 100000
	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData.Set(StrKey(strconv.Itoa(i)), d)
	}

	randNum := 100

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		mapData.Delete(StrKey(strconv.Itoa(j + randNum)))
		runtime.GC()
	}
	b.StopTimer()
}
