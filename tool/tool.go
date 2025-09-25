package tool

import (
	"hash/crc64"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var BOTPATH string

func init() {
	BOTPATH, _ = os.Getwd()
	if runtime.GOOS == "windows" {
		BOTPATH = strings.ReplaceAll(BOTPATH, "\\", "/")
	}
}

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

// IsExist 文件/路径存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsNotExist 文件/路径不存在
func IsNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
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

// 40位以下id字符串合并,结果称PadString
func MergePadString(group, guild string) string {
	return group + strings.Repeat(" ", 40-len(group)) + guild + strings.Repeat(" ", 40-len(guild))
}

// 40位以下id字符串拆分
func SplitPadString(only string) (group, guild string) {
	return strings.TrimRight(only[:40], " "), strings.TrimRight(only[40:], " ")
}

func MergeIntToInt64(group, guild int) int64 {
	return int64(guild)<<32 | int64(group)
}

func SplitInt64ToInt(x int64) (group, guild int) {
	group = int(x & 0xffffffff)
	guild = int(x >> 32)
	return
}

// JoinTypeAndPadString结果和type的组合数据
func JoinTypeAndPadString(types, padString string) string {
	return types + padString
}

// 解析JoinTypeAndString结果和type的组合数据
func SplitTypeAndPadString(value string) (types, padString string) {
	return value[:len(value)-80], value[len(value)-80:]
}

// 等待其他init加载完毕
func WaitWhile() {
	time.Sleep(time.Second * 2)
}

// HideURL 转义 URL 以避免审核
func HideURL(s string) string {
	s = strings.ReplaceAll(s, ".", "…")
	s = strings.ReplaceAll(s, "http://", "🔗📄:")
	s = strings.ReplaceAll(s, "https://", "🔗🔒:")
	return s
}

// RandSenderPerDayN 每个用户每天随机数  github.com/FloatTech/floatbox/ctxext
func RandSenderPerDayN(uid int64, n int) int {
	sum := crc64.New(crc64.MakeTable(crc64.ISO))
	sum.Write(StringToBytes(time.Now().Format("20060102")))
	sum.Write((*[8]byte)(unsafe.Pointer(&uid))[:])
	r := rand.New(rand.NewSource(int64(sum.Sum64())))
	return r.Intn(n)
}
