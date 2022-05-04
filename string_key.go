package HighPerformanceMap

import (
	"hash/crc64"
)

// StringKey is for the string type key
type stringKey struct {
	key   uint64
	value string
}

func hash(str string) uint64 {
	table := crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum([]byte(str), table)
}

func (s *stringKey) PartitionKey() uint64 {
	return s.key
}

func (s *stringKey) Value() any {
	return s.value
}

func StrKey(key string) *stringKey {
	return &stringKey{hash(key), key}
}
