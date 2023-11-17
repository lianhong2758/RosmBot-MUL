package tool

import (
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// BytesToString 没有内存开销的转换
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 没有内存开销的转换
func StringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

// 转int64
func Int64(ID string) int64 {
	x, _ := strconv.ParseInt(ID, 10, 64)
	return x
}

// 转字符串
func String(ID int64) string {
	return strconv.FormatInt(ID, 10)
}

// 20位以下id字符串合并
func String221(one, two string) string {
	return one + strings.Repeat(" ", 20-len(one)) + two + strings.Repeat(" ", 20-len(two))
}

// 20位以下id字符串拆分
func String122(only string) (one, two string) {
	return strings.TrimRight(only[:20], " "), strings.TrimRight(only[20:], " ")
}

func Int64_221(one, two int) int64 {
	var x int64
	x = int64(two) << 28
	x |= int64(one)
	return x
}

func Int64_122(x int64) (one, two int) {
	two = int(x >> 28)
	one = int(x & 0xfffffff)
	return
}
