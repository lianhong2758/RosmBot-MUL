package ob11

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// 新建bot消息
func NewBot(botid string) rosm.Boter {
	return botMap[botid]
}

// 新建上下文
func NewCTX(botid, roomid, villaid string) *rosm.Ctx {
	return &rosm.Ctx{
		BotType: "ob11",
		Bot:     botMap[botid],
		Being: &rosm.Being{
			RoomID:  roomid,
			RoomID2: villaid,
		},
	}
}
func GetRandBot() *Config {
	for k := range botMap {
		return botMap[k]
	}
	return nil
}

// RangeBot 遍历所有bot实例
func RangeBot(fn func(id string, bot *Config) bool) {
	for k, v := range botMap {
		if !fn(k, v) {
			return
		}
	}
}

var base64Reg = regexp.MustCompile(`"type":"image","data":\{"file":"base64://[\w/\+=]+`)

// formatMessage 格式化消息数组
//
//	仅用在 log 打印
func formatMessage(msg interface{}) string {
	switch m := msg.(type) {
	case string:
		return m
	case message.CQCoder:
		return m.CQCode()
	case fmt.Stringer:
		return m.String()
	default:
		s, _ := json.Marshal(msg)
		return helper.BytesToString(base64Reg.ReplaceAllFunc(s, func(b []byte) []byte {
			buf := bytes.NewBuffer([]byte(`"type":"image","data":{"file":"`))
			b = b[40:]
			b, err := base64.StdEncoding.DecodeString(helper.BytesToString(b))
			if err != nil {
				buf.WriteString(err.Error())
			} else {
				m := md5.Sum(b)
				_, _ = hex.NewEncoder(buf).Write(m[:])
				buf.WriteString(".image")
			}
			return buf.Bytes()
		}))
	}
}
