package qq

import (
	"encoding/base64"
	"net/url"
	"strings"
	"unsafe"
)

// UnderlineToCamel convert abc_def to AbcDef
func UnderlineToCamel(s string) string {
	sb := strings.Builder{}
	isnextupper := true
	for _, c := range []byte(strings.ToLower(s)) {
		if c == '_' {
			isnextupper = true
			continue
		}
		if isnextupper {
			sb.WriteString(strings.ToUpper(string(c)))
			isnextupper = false
			continue
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

// resolveURI github.com/wdvxdr1123/ZeroBot/driver/uri.go
func resolveURI(addr string) (network, address string) {
	network, address = "tcp", addr
	uri, err := url.Parse(addr)
	if err == nil && uri.Scheme != "" {
		scheme, ext, _ := strings.Cut(uri.Scheme, "+")
		if ext != "" {
			network = ext
			uri.Scheme = scheme // remove `+unix`/`+tcp4`
			if ext == "unix" {
				uri.Host, uri.Path, _ = strings.Cut(uri.Path, ":")
				uri.Host = base64.StdEncoding.EncodeToString(StringToBytes(uri.Host)) // special handle for unix
			}
			address = uri.String()
		}
	}
	return
}

// slice is the runtime representation of a slice.
// It cannot be used safely or portably and its representation may
// change in a later release.
//
// Unlike reflect.SliceHeader, its Data field is sufficient to guarantee the
// data it references will not be garbage collected.
type slice struct {
	data unsafe.Pointer
	len  int
	cap  int
}

// BytesToString 没有内存开销的转换
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 没有内存开销的转换
func StringToBytes(s string) (b []byte) {
	bh := (*slice)(unsafe.Pointer(&b))
	sh := (*slice)(unsafe.Pointer(&s))
	bh.data = sh.data
	bh.len = sh.len
	bh.cap = sh.len
	return b
}
