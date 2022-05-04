package HighPerformanceMap

import (
	"sync"
	"unsafe"
)

type concurrentMap struct {
	partitions  []*concurrentSliceMap // 对每个桶中的数据添加map
	lenOfBucket int                   // 分桶，目的加快map查找
	free        []int                 // 用户记录删除切片的位置
	innerSlice  []*innerSlice         // 用户记录所用的值的位置
	mu          sync.RWMutex
}

type concurrentSliceMap struct {
	index map[uint64]int
	mu    sync.RWMutex
}

type innerSlice struct {
	key   interface{}
	Value unsafe.Pointer
}

type Partitionable interface {
	Value() interface{}
	PartitionKey() uint64
}

func CreateConcurrentSliceMap(lenOfBucket int) *concurrentMap {
	partitions := make([]*concurrentSliceMap, lenOfBucket)
	for i := 0; i < lenOfBucket; i++ {
		partitions[i] = &concurrentSliceMap{
			index: make(map[uint64]int),
		}
	}
	return &concurrentMap{
		partitions:  partitions,
		lenOfBucket: lenOfBucket,
		free:        make([]int, 0, 1024),
		innerSlice:  make([]*innerSlice, 0, 1024),
	}
}

func (m *concurrentMap) getPartition(key Partitionable) *concurrentSliceMap {
	partitionID := key.PartitionKey() % uint64(m.lenOfBucket)
	return m.partitions[partitionID]
}

func (m *concurrentMap) getValue(v unsafe.Pointer) interface{} {
	return *(*interface{})(v)
}

func (m *concurrentMap) Len() int {
	return len(m.innerSlice)
}

func (m *concurrentMap) Range(f func(key, value interface{}) bool) {
	m.mu.RLock()
	for _, data := range m.innerSlice {
		if !f(data.key, m.getValue(data.Value)) {
			break
		}

	}
	m.mu.RUnlock()
}

func (m *concurrentMap) Get(key Partitionable) (interface{}, bool) {
	im := m.getPartition(key)
	im.mu.RLock()
	defer im.mu.RUnlock()

	keyIndex := key.PartitionKey()
	if index, ok := im.index[keyIndex]; ok {
		data := m.innerSlice[index]
		return m.getValue(data.Value), true
	}

	return nil, false
}

func (m *concurrentMap) Set(key Partitionable, v interface{}) {
	im := m.getPartition(key)
	im.mu.Lock()
	defer im.mu.Unlock()

	keyIndex := key.PartitionKey()

	if index, ok := im.index[keyIndex]; ok {
		m.innerSlice[index] = &innerSlice{
			key.Value(),
			unsafe.Pointer(&v),
		}
		return
	}

	n := len(m.innerSlice)
	if len(m.free) > 0 {
		n = m.free[0]
		m.free = m.free[1:]
	}

	m.innerSlice = append(m.innerSlice, &innerSlice{
		key.Value(),
		unsafe.Pointer(&v),
	})
	im.index[keyIndex] = n
}

func (m *concurrentMap) Delete(key Partitionable) {
	im := m.getPartition(key)
	im.mu.Lock()
	defer im.mu.Unlock()

	keyIndex := key.PartitionKey()
	if index, ok := im.index[keyIndex]; ok {
		m.free = append(m.free, index)
		m.innerSlice[index] = nil
		delete(im.index, keyIndex)
	}
}
