package HighPerformanceMap

import (
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
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
	mapData := make(map[string]int)
	mu := sync.RWMutex{}
	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			mu.Lock()
			mapData[strconv.Itoa(id)] = id
			mu.Unlock()
		}
	})
}

func BenchmarkSyncAndMapAndPMapSetB(b *testing.B) {
	mapData := sync.Map{}
	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			mapData.Store(strconv.Itoa(id), id)
		}
	})
}

func BenchmarkSyncAndMapAndPMapSetC(b *testing.B) {
	mapData := CreateConcurrentSliceMap(99)
	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			mapData.Set(StrKey(strconv.Itoa(id)), id)
		}
	})
}

// performance Test Read
func BenchmarkSyncAndMapAndPMapReadA(b *testing.B) {
	num := 100
	mapData := make(map[string]int)
	for i := 0; i < num; i++ {
		mapData[strconv.Itoa(i)] = i
	}

	mu := sync.RWMutex{}
	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			mu.RLock()
			if _, ok := mapData[strconv.Itoa(id)]; ok {

			} else {
				b.Errorf("error %v", id)
			}
			mu.RUnlock()
		}
	})
}

func BenchmarkSyncAndMapAndPMapReadB(b *testing.B) {
	num := 100
	mapData := sync.Map{}
	for i := 0; i < num; i++ {
		mapData.Store(strconv.Itoa(i), i)
	}

	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			if _, ok := mapData.Load(strconv.Itoa(id)); ok {

			} else {
				b.Errorf("error %v", id)
			}
		}
	})
}

func BenchmarkSyncAndMapAndPMapReadC(b *testing.B) {
	num := 100
	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		mapData.Set(StrKey(strconv.Itoa(i)), i)
	}

	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			if _, ok := mapData.Get(StrKey(strconv.Itoa(id))); ok {

			} else {
				b.Errorf("error %v", id)
			}
		}
	})
}

// performance Test delete
func BenchmarkSyncAndMapAndPMapDeleteA(b *testing.B) {
	num := 100
	mapData := make(map[string]int)
	for i := 0; i < num; i++ {
		mapData[strconv.Itoa(i)] = i
	}

	mu := sync.RWMutex{}
	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			mu.RLock()
			_, ok := mapData[strconv.Itoa(id)]
			mu.RUnlock()
			if ok {
				mu.Lock()
				delete(mapData, strconv.Itoa(id))
				mu.Unlock()
			}
		}
	})
	b.StopTimer()
}

func BenchmarkSyncAndMapAndPMapDeleteB(b *testing.B) {
	num := 100
	mapData := sync.Map{}
	for i := 0; i < num; i++ {
		mapData.Store(strconv.Itoa(i), i)
	}

	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			mapData.Delete(strconv.Itoa(id))
		}
	})
}

func BenchmarkSyncAndMapAndPMapDeleteC(b *testing.B) {
	num := 100
	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		mapData.Set(StrKey(strconv.Itoa(i)), i)
	}

	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			mapData.Delete(StrKey(strconv.Itoa(id)))
		}
	})
}

// performance Test GC recycle
func TestSyncAndMapAndPMapGCA(t *testing.T) {
	num := 1000000
	mapData := make(map[string]int)
	for i := 0; i < num; i++ {
		mapData[strconv.Itoa(i)] = i
	}
	for i := 0; i < 100; i++ {
		delete(mapData, strconv.Itoa(i+100))
	}

	now := time.Now()
	runtime.GC()
	t.Logf("With a map of strings, GC took: %s\n", time.Since(now))

	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	t.Logf("numGC --> %v, PauseTotal --> %v", stats.NumGC, stats.PauseTotal)
	t.Logf("len --> %v", len(mapData))
}

func TestSyncAndMapAndPMapGCB(t *testing.T) {
	num := 1000000
	mapData := sync.Map{}
	for i := 0; i < num; i++ {
		mapData.Store(strconv.Itoa(i), i)
	}

	for i := 0; i < 100; i++ {
		mapData.Delete(strconv.Itoa(i + 100))
	}

	now := time.Now()
	runtime.GC()
	t.Logf("With a sync map of strings, GC took: %s\n", time.Since(now))

	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	t.Logf("numGC --> %v, PauseTotal --> %v", stats.NumGC, stats.PauseTotal)

	_, ok := mapData.Load(strconv.Itoa(num - 1))
	t.Logf("ok --> %v", ok)
}

func TestSyncAndMapAndPMapGCC(t *testing.T) {
	num := 1000000
	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		mapData.Set(StrKey(strconv.Itoa(i)), i)
	}

	for i := 0; i < 100; i++ {
		mapData.Delete(StrKey(strconv.Itoa(i + 100)))
	}

	now := time.Now()
	runtime.GC()
	t.Logf("With a PMap of strings, GC took: %s\n", time.Since(now))

	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	t.Logf("numGC --> %v, PauseTotal --> %v", stats.NumGC, stats.PauseTotal)
	t.Logf("len --> %v", mapData.Len())
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
	mapData := make(map[string]*intBig)
	mu := sync.RWMutex{}
	i := int64(0)
	id := int64(0)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = atomic.AddInt64(&i, 1)
		for pb.Next() {
			mu.Lock()
			d := &intBig{
				id, id, id, id, id,
			}
			mapData[strconv.FormatInt(id, 10)] = d
			mu.Unlock()
		}
	})
}

func BenchmarkSyncAndMapAndPMapBigSetB(b *testing.B) {
	mapData := sync.Map{}
	i := int64(0)
	id := int64(0)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = atomic.AddInt64(&i, 1)
		for pb.Next() {
			d := &intBig{
				id, id, id, id, id,
			}
			mapData.Store(strconv.FormatInt(id, 10), d)
		}
	})
}

func BenchmarkSyncAndMapAndPMapBigSetC(b *testing.B) {
	mapData := CreateConcurrentSliceMap(99)
	i := int64(0)
	id := int64(0)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = atomic.AddInt64(&i, 1)
		for pb.Next() {
			d := &intBig{
				id, id, id, id, id,
			}
			mapData.Set(StrKey(strconv.FormatInt(id, 10)), d)
		}
	})
}

// performance Test Read
func BenchmarkSyncAndMapAndPMapBigReadA(b *testing.B) {
	num := 100
	mapData := make(map[string]*intBig)
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData[strconv.Itoa(i)] = d
	}

	mu := sync.RWMutex{}
	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			mu.RLock()
			if _, ok := mapData[strconv.Itoa(id)]; ok {

			} else {
				b.Errorf("error %v", id)
			}
			mu.RUnlock()
		}
	})
}

func BenchmarkSyncAndMapAndPMapBigReadB(b *testing.B) {
	num := 100
	mapData := sync.Map{}
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData.Store(strconv.Itoa(i), d)
	}

	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			if _, ok := mapData.Load(strconv.Itoa(id)); ok {

			} else {
				b.Errorf("error %v", id)
			}
		}
	})
}

func BenchmarkSyncAndMapAndPMapBigReadC(b *testing.B) {
	num := 100
	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData.Set(StrKey(strconv.Itoa(i)), d)
	}

	i := int64(0)
	id := 0

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id = int(atomic.AddInt64(&i, 1))
		for pb.Next() {
			if _, ok := mapData.Get(StrKey(strconv.Itoa(id))); ok {

			} else {
				b.Errorf("error %v", id)
			}
		}
	})
}

// performance Test GC recycle
func TestSyncAndMapAndPMapBigGCA(t *testing.T) {
	num := 1000000
	mapData := make(map[string]*intBig)
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData[strconv.Itoa(i)] = d
	}
	for i := 0; i < 10000; i++ {
		delete(mapData, strconv.Itoa(i))
	}

	now := time.Now()
	runtime.GC()
	t.Logf("With a map of strings, GC took: %s\n", time.Since(now))

	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	t.Logf("numGC --> %v, PauseTotal --> %v", stats.NumGC, stats.PauseTotal)
	t.Logf("len --> %v", len(mapData))
}

func TestSyncAndMapAndPMapBigGCB(t *testing.T) {
	num := 1000000
	mapData := sync.Map{}
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData.Store(strconv.Itoa(i), d)
	}

	for i := 0; i < 10000; i++ {
		mapData.Delete(strconv.Itoa(i))
	}

	now := time.Now()
	runtime.GC()
	t.Logf("With a sync map of strings, GC took: %s\n", time.Since(now))

	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	t.Logf("numGC --> %v, PauseTotal --> %v", stats.NumGC, stats.PauseTotal)

	_, ok := mapData.Load(strconv.Itoa(num - 1))
	t.Logf("ok --> %v", ok)
}

func TestSyncAndMapAndPMapBigGCC(t *testing.T) {
	num := 1000000
	mapData := CreateConcurrentSliceMap(99)
	for i := 0; i < num; i++ {
		iBig := int64(i)
		d := &intBig{
			iBig, iBig, iBig, iBig, iBig,
		}
		mapData.Set(StrKey(strconv.Itoa(i)), d)
	}

	for i := 0; i < 10000; i++ {
		mapData.Delete(StrKey(strconv.Itoa(i)))
	}

	now := time.Now()
	runtime.GC()
	t.Logf("With a PMap of strings, GC took: %s\n", time.Since(now))

	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	t.Logf("numGC --> %v, PauseTotal --> %v", stats.NumGC, stats.PauseTotal)
	t.Logf("len --> %v", mapData.Len())
}
