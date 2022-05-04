package HighPerformanceMap

type int64Key struct {
	value uint64
}

func (i *int64Key) PartitionKey() uint64 {
	return i.value
}

// Value is the raw string
func (i *int64Key) Value() any {
	return i.value
}

func I64Key(key int64) *int64Key {
	return &int64Key{uint64(key)}
}
