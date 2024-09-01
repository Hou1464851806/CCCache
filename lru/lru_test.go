package lru

import "testing"

// 自定义string类型，实现Value接口，从而能被Cache使用
type String string

func (d String) Len() int {
	return len(d)
}

// 测试Cache Get方法
func TestGet(t *testing.T) {
	lru := NewCache(int64(0), nil)
	lru.Add("key1", String("1111"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1111" {
		t.Fatalf("hashmap hit key1:1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("hashmap miss key2 failed")
	}
}
