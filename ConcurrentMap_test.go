package HighPerformanceMap

import (
	"strconv"
	"sync"
	"testing"
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

	mapData.Range(func(key, value interface{}) bool {
		t.Logf("key --> %v", key.(string))
		t.Logf("value --> %v", value.(int))
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
