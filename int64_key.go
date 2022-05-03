package HighPerformanceMap

type Int64Key struct {
	value uint64
}

func (i *Int64Key) PartitionKey() uint64 {
	return i.value
}

// Value is the raw string
func (i *Int64Key) Value() interface{} {
	return i.value
}

func I64Key(key int64) *Int64Key {
	return &Int64Key{uint64(key)}
}
