package Cache

import "bytes"

// byteView 模块定义读取缓存结果
// 实际上 byteView 只是简单的封装了byte slice，让其只读。
// 试想一下，直接返回slice，在golang里，一切参数按值传递。
// slice底层只是一个struct，记录着ptr/len/cap，相当于
// 复制了一份这三者的值。因此[]byte底层指向同一片内存区域
// 我们的缓存底层是存储在LRU的双向链表的Element里，因此
// 可以被恶意修改。因此需要将slice封装成只读的ByteView

// ByteView
// 一个只读的字节数组
// byte以二进制方式存储任意类型数据
type ByteView struct {
	b []byte
}

// Len
// 被缓存对象需要实现Value，即实现Len()方法
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice
// 返回byteView字节数组的拷贝，byteView需要保证只读
func (v ByteView) ByteSlice() []byte {
	return bytes.Clone(v.b)
}

// String
// 将字节数组转换为string类型
func (v ByteView) String() string {
	return string(v.b)
}
