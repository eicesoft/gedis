package dict

import (
	"fmt"
	"testing"
)

func TestSimpleDict(t *testing.T) {
	d := MakeSimple()
	key := "A1"
	key2 := "B1"
	key3 := "C1"
	val := "AAABBBB"
	d.Put(key, "abc")
	if d.Exists(key) {
		t.Logf("Exists: Key %s 存在.", key)
	} else {
		t.Errorf("xists: Key %s 不存在.", key)
	}

	if d.Put(key, "aaa") == 0 {
		t.Logf("Put: Key %s 存在, 并重设.", key)
	} else {
		t.Logf("Put: Key %s 不存在, 并重设", key)
	}

	if d.PutIfAbsent(key2, val) == 1 {
		t.Logf("PutIfAbsent: Key %s 不存在, 重设成功.", key)
	} else {
		t.Errorf("PutIfAbsent: Key %s 存在, 重设失败.", key)
	}

	if d.PutIfExists(key3, val) == 0 {
		t.Logf("PutIfExists: Key %s 不存在, 重设失败.", key)
	} else {
		t.Logf("PutIfExists: Key %s 存在, 重设成功.", key)
	}

	v, _ := d.Get(key2)
	if v == val {
		t.Logf("Get: Key值设置对应.")
	} else {
		t.Errorf("Get: Key值设置不对应.")
	}

	l := d.Len()
	if l == 2 {
		t.Logf("Len: Key数量正确.")
	} else {
		t.Errorf("Len: Key数量异常.")
	}

	for _, k := range d.Keys() {
		t.Logf("Keys: Key: %s", k)
	}

	if d.Remove(key2) == 1 {
		t.Logf("Remove: 删除Key成功.")
	} else {
		t.Errorf("Remove: 删除Key失败.")
	}
	k := d.RandomKeys(1)
	if len(k) == 1 {
		t.Logf("RandomKeys: 随机获取成功.")
	}
}

func TestConcurrentDict(t *testing.T) {
	key := "A1"
	key2 := "B1"
	key3 := "C1"
	val := "AAABBBB"
	d := MakeConcurrent(0)
	d.Put(key, "abc")
	if d.Exists(key) {
		t.Logf("Exists: Key %s 存在.", key)
	} else {
		t.Errorf("xists: Key %s 不存在.", key)
	}

	if d.Put(key, "aaa") == 0 {
		t.Logf("Put: Key %s 存在, 并重设.", key)
	} else {
		t.Logf("Put: Key %s 不存在, 并重设", key)
	}

	if d.PutIfAbsent(key2, val) == 1 {
		t.Logf("PutIfAbsent: Key %s 不存在, 重设成功.", key)
	} else {
		t.Errorf("PutIfAbsent: Key %s 存在, 重设失败.", key)
	}

	if d.PutIfExists(key3, val) == 0 {
		t.Logf("PutIfExists: Key %s 不存在, 重设失败.", key)
	} else {
		t.Logf("PutIfExists: Key %s 存在, 重设成功.", key)
	}

	v, _ := d.Get(key2)
	if v == val {
		t.Logf("Get: Key值设置对应.")
	} else {
		t.Errorf("Get: Key值设置不对应.")
	}

	l := d.Len()
	if l == 2 {
		t.Logf("Len: Key数量正确.")
	} else {
		t.Errorf("Len: Key数量异常.")
	}

	for _, k := range d.Keys() {
		t.Logf("Keys: Key: %s", k)
	}

	if d.Remove(key2) == 1 {
		t.Logf("Remove: 删除Key成功.")
	} else {
		t.Errorf("Remove: 删除Key失败.")
	}
	k := d.RandomKeys(1)
	t.Logf("RandomKeys: %v", k)
	if len(k) == 1 {
		t.Logf("RandomKeys: 随机获取成功.")
	}

	for i := 0; i < 10; i++ {
		d.Put(fmt.Sprintf("AA:%d", i), "DDD")
	}
	t.Logf("Dict: %d", d.Len())

	keys := d.ForScanKeys(func(key string, val interface{}) bool {
		return true
	}, 0, 5)
	t.Logf("%v", keys)

	keys = d.ForScanKeys(func(key string, val interface{}) bool {
		return true
	}, 5, 5)
	t.Logf("%v", keys)

	keys = d.ForScanKeys(func(key string, val interface{}) bool {
		return true
	}, 10, 5)
	t.Logf("%v", keys)
}
