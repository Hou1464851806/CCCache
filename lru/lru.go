// 实现了内部缓存的组织结构
// 提供了增删查改方法
// 为了加快查询，使用hashmap通过key查询value所在的列表元素
// 双向链表用于维护各元素的最近访问性

package lru

import "container/list"

// Cache
// 缓存基于LRU算法实现
// 非并发安全
type Cache struct {
	doubleList *list.List               // 双向链表
	hashmap    map[string]*list.Element // 查询用哈希表
	capacity   int64                    //最大可用空间
	used       int64                    //已用空间
	// 回调函数: 当某个k-v对被移除时触发
	callback OnEliminated
}

// OnEliminated
// 定义一个删除回调函数类型
type OnEliminated func(string, Value)

// 缓存实体节点
// k-v对：key为字符串，value为实现了Value接口的任意类型，保证通用性
type entry struct {
	key   string
	value Value
}

// Value
// 包含一个方法，返回value所占空间大小
type Value interface {
	Len() int
}

// NewCache
// Cache构造方法
// capacity为0，表示不淘汰元素，所占空间无上限
func NewCache(capacity int64, callback OnEliminated) *Cache {
	return &Cache{
		doubleList: list.New(),
		hashmap:    make(map[string]*list.Element),
		capacity:   capacity,
		used:       0,
		callback:   callback,
	}
}

// Get
// 通过key查找value
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 通过map找到key所对列表元素
	if element, ok := c.hashmap[key]; ok {
		// 取出列表元素中的值，类型断言转为entry
		kv := element.Value.(*entry)
		// 访问后将元素移至队首，表明最近访问过
		c.doubleList.MoveToFront(element)
		return kv.value, true
	}
	return
}

// RemoveOldest
// 移除最近最长时间未访问的元素
func (c *Cache) RemoveOldest() {
	// 从队尾取出元素，队尾是最近最长时间未访问的元素
	element := c.doubleList.Back()
	if element != nil {
		// 移除
		c.doubleList.Remove(element)
		kv := element.Value.(*entry)
		// 从map中移除
		delete(c.hashmap, kv.key)
		// 所占内存空间减去key的大小和value的大小
		c.used -= int64(len(kv.key) + kv.value.Len())
		// 如果回调函数不为空，执行回调函数
		if c.callback != nil {
			c.callback(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// 如果存在，则是更新操作
	if element, ok := c.hashmap[key]; ok {
		kv := element.Value.(*entry)
		// 更新value值
		kv.value = value
		// 更新所占内存空间
		c.used += int64(kv.value.Len() - value.Len())
		// 移至队首，表示最近访问过
		c.doubleList.MoveToFront(element)
	} else {
		// 新建kv对
		kv := &entry{
			key:   key,
			value: value,
		}
		// 插入队首
		ele := c.doubleList.PushFront(kv)
		// 建立mapping
		c.hashmap[key] = ele
		// 更新所占内存空间
		c.used += int64(len(key) + value.Len())
	}
	// 淘汰元素直到所占空间不超过最大空间
	for c.capacity != 0 && c.used > c.capacity {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.doubleList.Len()
}
