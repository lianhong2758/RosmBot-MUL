package ob11

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/tidwall/gjson"
)

// CQString 转为CQ字符串
// Deprecated: use method String instead
func (m Message) CQString() string {
	return m.String()
}

// Modified from https://github.com/catsworld/qq-bot-api

// ParseMessage parses msg, which might have 2 types, string or array,
// depending on the configuration of cqhttp, to a Message.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
func ParseMessage(msg []byte) message.Message {
	x := gjson.Parse(tool.BytesToString(msg))
	if x.IsArray() {
		return ParseMessageFromArray(x)
	} else {
		return ParseMessageFromString(x.String())
	}
}

// ParseMessageFromArray parses msg as type array to a Message.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
// ParseMessageFromArray cq字符串转化为json对象
func ParseMessageFromArray(msgs gjson.Result) message.Message {
	messagee := message.Message{}
	parse2map := func(val gjson.Result) map[string]string {
		m := map[string]string{}
		val.ForEach(func(key, value gjson.Result) bool {
			m[key.String()] = value.String()
			return true
		})
		return m
	}
	msgs.ForEach(func(_, item gjson.Result) bool {
		messagee = append(messagee, message.MessageSegment{
			Type: item.Get("type").String(),
			Data: parse2map(item.Get("data")),
		})
		return true
	})
	return messagee
}

// ParseMessageFromString parses msg as type string to a sort of MessageSegment.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
//
// CQ字符串转为消息
func ParseMessageFromString(raw string) (m message.Message) {
	var seg message.MessageSegment
	var k string
	m = message.Message{}
	for raw != "" {
		i := 0
		for i < len(raw) && !(raw[i] == '[' && i+4 < len(raw) && raw[i:i+4] == "[CQ:") {
			i++
		}

		if i > 0 {
			m = append(m, message.Text(UnescapeCQText(raw[:i])))
		}

		if i+4 > len(raw) {
			return
		}

		raw = raw[i+4:] // skip "[CQ:"
		i = 0
		for i < len(raw) && raw[i] != ',' && raw[i] != ']' {
			i++
		}

		if i+1 > len(raw) {
			return
		}
		seg.Type = raw[:i]
		seg.Data = make(map[string]string)
		raw = raw[i:]
		i = 0

		for {
			if raw[0] == ']' {
				m = append(m, seg)
				raw = raw[1:]
				break
			}
			raw = raw[1:]

			for i < len(raw) && raw[i] != '=' {
				i++
			}
			if i+1 > len(raw) {
				return
			}
			k = raw[:i]
			raw = raw[i+1:] // skip "="
			i = 0
			for i < len(raw) && raw[i] != ',' && raw[i] != ']' {
				i++
			}

			if i+1 > len(raw) {
				return
			}
			seg.Data[k] = UnescapeCQCodeText(raw[:i])
			raw = raw[i:]
			i = 0
		}
	}
	return m
}

// CQCoder 用于 log 打印 CQ 码
type CQCoder interface {
	CQCode() string
}

// EscapeCQText escapes special characters in a non-media plain message.\
//
// CQ码字符转换
func EscapeCQText(str string) string {
	str = strings.ReplaceAll(str, "&", "&amp;")
	str = strings.ReplaceAll(str, "[", "&#91;")
	str = strings.ReplaceAll(str, "]", "&#93;")
	return str
}

// UnescapeCQText unescapes special characters in a non-media plain message.
//
// CQ码反解析
func UnescapeCQText(str string) string {
	str = strings.ReplaceAll(str, "&#93;", "]")
	str = strings.ReplaceAll(str, "&#91;", "[")
	str = strings.ReplaceAll(str, "&amp;", "&")
	return str
}

// EscapeCQCodeText escapes special characters in a cqcode value.
//
// https://github.com/botuniverse/onebot-11/tree/master/message/string.md#%E8%BD%AC%E4%B9%89
//
// cq码字符转换
func EscapeCQCodeText(str string) string {
	str = strings.ReplaceAll(str, "&", "&amp;")
	str = strings.ReplaceAll(str, "[", "&#91;")
	str = strings.ReplaceAll(str, "]", "&#93;")
	str = strings.ReplaceAll(str, ",", "&#44;")
	return str
}

// UnescapeCQCodeText unescapes special characters in a cqcode value.
// https://github.com/botuniverse/onebot-11/tree/master/message/string.md#%E8%BD%AC%E4%B9%89
//
// cq码反解析
func UnescapeCQCodeText(str string) string {
	str = strings.ReplaceAll(str, "&#44;", ",")
	str = strings.ReplaceAll(str, "&#93;", "]")
	str = strings.ReplaceAll(str, "&#91;", "[")
	str = strings.ReplaceAll(str, "&amp;", "&")
	return str
}

type MessageSegment message.MessageSegment
type Message message.Message

// CQCode 将数组消息转换为CQ码
// 与 String 不同之处在于，对于
// base64 的图片消息会将其哈希
// 方便 log 打印，不可用作发送
func (m MessageSegment) CQCode() string {
	sb := strings.Builder{}
	sb.WriteString("[CQ:")
	sb.WriteString(m.Type)
	for k, v := range m.Data { // 消息参数
		// sb.WriteString("," + k + "=" + escape(v))
		sb.WriteByte(',')
		sb.WriteString(k)
		sb.WriteByte('=')
		switch m.Type {
		case "node":
			sb.WriteString(v)
		case "image":
			if strings.HasPrefix(v, "base64://") {
				v = v[9:]
				b, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					sb.WriteString(err.Error())
				} else {
					m := md5.Sum(b)
					_, _ = hex.NewEncoder(&sb).Write(m[:])
				}
				sb.WriteString(".image")
				break
			}
			fallthrough
		default:
			sb.WriteString(EscapeCQCodeText(v))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

// String impls the interface fmt.Stringer
func (m MessageSegment) String() string {
	sb := strings.Builder{}
	sb.WriteString("[CQ:")
	sb.WriteString(m.Type)
	for k, v := range m.Data { // 消息参数
		// sb.WriteString("," + k + "=" + escape(v))
		sb.WriteByte(',')
		sb.WriteString(k)
		sb.WriteByte('=')
		if m.Type == "node" {
			sb.WriteString(v)
		} else {
			sb.WriteString(EscapeCQCodeText(v))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

// CQCode 将数组消息转换为CQ码
// 与 String 不同之处在于，对于
// base64 的图片消息会将其哈希
// 方便 log 打印，不可用作发送
func (m Message) CQCode() string {
	sb := strings.Builder{}
	for _, media := range m {
		if media.Type != "text" {
			sb.WriteString(MessageSegment(media).CQCode())
		} else {
			sb.WriteString(EscapeCQText(media.Data["text"]))
		}
	}
	return sb.String()
}

// String impls the interface fmt.Stringer
func (m Message) String() string {
	sb := strings.Builder{}
	for _, media := range m {
		if media.Type != "text" {
			sb.WriteString(MessageSegment(media).String())
		} else {
			sb.WriteString(EscapeCQText(media.Data["text"]))
		}
	}
	return sb.String()
}
