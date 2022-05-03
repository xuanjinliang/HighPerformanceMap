package HighPerformanceMap

import (
	"sync"
	"unsafe"
)

type ConcurrentMap struct {
	partitions  []*ConcurrentSliceMap // 对每个桶中的数据添加map
	lenOfBucket int                   // 分桶，目的加快map查找
	free        []int                 // 用户记录删除切片的位置
	innerSlice  []unsafe.Pointer      // 用户记录所用的值的位置
}

type ConcurrentSliceMap struct {
	index map[uint64]int
	mu    sync.RWMutex
}

type Partitionable interface {
	Value() interface{}
	PartitionKey() uint64
}

func CreateConcurrentSliceMap(lenOfBucket int) *ConcurrentMap {
	partitions := make([]*ConcurrentSliceMap, lenOfBucket)
	for i := 0; i < lenOfBucket; i++ {
		partitions[i] = &ConcurrentSliceMap{
			index: make(map[uint64]int),
		}
	}
	return &ConcurrentMap{
		partitions:  partitions,
		lenOfBucket: lenOfBucket,
		free:        make([]int, 0, 1024),
		innerSlice:  make([]unsafe.Pointer, 0, 1024),
	}
}

/*
func (m *ConcurrentSliceMap) Len() int {
	return 0
}

func (m *ConcurrentSliceMap) Range(f func(key, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
}*/

func (m *ConcurrentMap) getPartition(key Partitionable) *ConcurrentSliceMap {
	partitionID := key.PartitionKey() % uint64(m.lenOfBucket)
	return m.partitions[partitionID]
}

func (m *ConcurrentMap) Get(key Partitionable) (interface{}, bool) {
	im := m.getPartition(key)
	im.mu.RLock()
	defer im.mu.RUnlock()

	keyIndex := key.PartitionKey()
	if index, ok := im.index[keyIndex]; ok {
		return *(*interface{})(m.innerSlice[index]), true
	}

	return nil, false
}

func (m *ConcurrentMap) Set(key Partitionable, v interface{}) {
	im := m.getPartition(key)
	im.mu.Lock()
	defer im.mu.Unlock()

	keyIndex := key.PartitionKey()

	if index, ok := im.index[keyIndex]; ok {
		m.innerSlice[index] = unsafe.Pointer(&v)
		return
	}

	n := len(m.innerSlice)
	if len(m.free) > 0 {
		n = m.free[0]
		m.free = m.free[1:]
	}

	m.innerSlice = append(m.innerSlice, unsafe.Pointer(&v))
	im.index[keyIndex] = n
}

func (m *ConcurrentMap) Delete(key Partitionable) {
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
