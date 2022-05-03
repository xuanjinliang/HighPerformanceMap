package HighPerformanceMap

import "testing"

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