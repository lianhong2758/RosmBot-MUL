package qq

import (
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
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
				uri.Host = base64.StdEncoding.EncodeToString(tool.StringToBytes(uri.Host)) // special handle for unix
			}
			address = uri.String()
		}
	}
	return
}

// 新建bot消息
func NewBot(botid string) rosm.Boter {
	return botMap[botid]
}

// 新建上下文
func NewCTX(botid, types, GroupID, GuildID string) *rosm.Ctx {
	return &rosm.Ctx{
		BotType: types,
		Bot:     botMap[botid],
		Being: &rosm.Being{
			GroupID: GroupID,
			GuildID: GuildID,
		},
	}
}
func GetRandBot() *Config {
	for k := range botMap {
		return botMap[k]
	}
	return nil
}
