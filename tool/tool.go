package tool

import (
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
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// 转int64
func StringToInt64(ID string) int64 {
	x, _ := strconv.ParseInt(ID, 10, 64)
	return x
}

// 转字符串
func Int64ToString(ID int64) string {
	return strconv.FormatInt(ID, 10)
}

// 40位以下id字符串合并,结果称为string1
func MergePadString(one, two string) string {
	return one + strings.Repeat(" ", 40-len(one)) + two + strings.Repeat(" ", 40-len(two))
}

// 40位以下id字符串拆分,结果称为string2
func SplitPadString(only string) (one, two string) {
	return strings.TrimRight(only[:40], " "), strings.TrimRight(only[40:], " ")
}

func MergeIntToInt64(one, two int) int64 {
	return int64(two)<<32 | int64(one)
}

func SplitInt64ToInt(x int64) (one, two int) {
	one = int(x & 0xffffffff)
	two = int(x >> 32)
	return
}

// JoinTypeAndString结果和type的组合数据
func JoinTypeAndString(types, string1 string) string {
	return types + string1
}

// 解析JoinTypeAndString结果和type的组合数据
func SplitTypeAndString(value string) (types, string1 string) {
	return value[:len(value)-80], value[len(value)-80:]
}
