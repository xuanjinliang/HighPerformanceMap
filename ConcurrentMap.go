package HighPerformanceMap

import (
	"sync"
	"unsafe"
)

type concurrentMap struct {
	partitions  []map[uint64]int // 对每个桶中的数据添加map
	lenOfBucket int              // 分桶，目的加快map查找
	free        []int            // 用户记录删除切片的位置
	innerSlice  []*innerSlice    // 用户记录所用的值的位置
	mu          sync.RWMutex
}

type innerSlice struct {
	key   any
	Value unsafe.Pointer
}

type Partitionable interface {
	Value() any
	PartitionKey() uint64
}

func CreateConcurrentSliceMap(lenOfBucket int) *concurrentMap {
	partitions := make([]map[uint64]int, lenOfBucket)
	for i := 0; i < lenOfBucket; i++ {
		partitions[i] = make(map[uint64]int)
	}
	return &concurrentMap{
		partitions:  partitions,
		lenOfBucket: lenOfBucket,
		free:        make([]int, 0, 1024),
		innerSlice:  make([]*innerSlice, 0, 1024),
	}
}

func (m *concurrentMap) getPartition(key Partitionable) map[uint64]int {
	partitionID := key.PartitionKey() % uint64(m.lenOfBucket)
	return m.partitions[partitionID]
}

func (m *concurrentMap) getValue(v unsafe.Pointer) any {
	return *(*any)(v)
}

func (m *concurrentMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	length := 0
	for _, data := range m.partitions {
		length += len(data)
	}

	return length
}

func (m *concurrentMap) Range(f func(key, value any) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, mapData := range m.partitions {
		for _, index := range mapData {
			data := m.innerSlice[index]
			if !f(data.key, m.getValue(data.Value)) {
				return
			}
		}
	}
}

func (m *concurrentMap) FreeLen() int {
	return len(m.free)
}

func (m *concurrentMap) Get(key Partitionable) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	im := m.getPartition(key)

	keyIndex := key.PartitionKey()
	if index, ok := im[keyIndex]; ok {
		data := m.innerSlice[index]
		return m.getValue(data.Value), true
	}

	return nil, false
}

func (m *concurrentMap) Set(key Partitionable, v any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	im := m.getPartition(key)

	keyIndex := key.PartitionKey()

	if index, ok := im[keyIndex]; ok {
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
	im[keyIndex] = n
}

func (m *concurrentMap) Delete(key Partitionable) {
	m.mu.Lock()
	defer m.mu.Unlock()

	im := m.getPartition(key)

	keyIndex := key.PartitionKey()
	if index, ok := im[keyIndex]; ok {
		m.free = append(m.free, index)
		m.innerSlice[index] = nil
		delete(im, keyIndex)
	}
}
