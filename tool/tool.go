package tool

import (
	"reflect"
	"strconv"
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
