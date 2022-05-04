package HighPerformanceMap

import (
	"sync"
	"unsafe"
)

type ConcurrentMap struct {
	partitions  []*ConcurrentSliceMap // 对每个桶中的数据添加map
	lenOfBucket int                   // 分桶，目的加快map查找
	free        []int                 // 用户记录删除切片的位置
	innerSlice  []*innerSlice         // 用户记录所用的值的位置
	mu          sync.RWMutex
}

type ConcurrentSliceMap struct {
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
		innerSlice:  make([]*innerSlice, 0, 1024),
	}
}

func (m *ConcurrentMap) Len() int {
	length := 0
	for _, v := range m.partitions {
		length += len(v.index)
	}
	return length
}

func (m *ConcurrentMap) Range(f func(key, value interface{}) bool) {
	m.mu.RLock()
	for _, data := range m.innerSlice {
		if !f(data.key, data.Value) {
			break
		}

	}
	m.mu.RUnlock()
}

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
		data := m.innerSlice[index]
		return *(*interface{})(data.Value), true
	}

	return nil, false
}

func (m *ConcurrentMap) Set(key Partitionable, v interface{}) {
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
