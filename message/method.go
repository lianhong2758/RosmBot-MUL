package message

import "strings"

//合并消息内连续的纯文本段。
func (message *Message) Reduce() {
	for index := 0; index < len(*message)-1; {
		if (*message)[index].Type == "text" && (*message)[index+1].Type == "text" {
			(*message)[index].Data["text"] += (*message)[index+1].Data["text"]
			*message = append((*message)[:index+1], (*message)[index+2:]...)
		} else {
			index++
		}
	}
}

// ExtractPlainText 提取消息中的纯文本
func (m Message) ExtractPlainText() string {
	sb := strings.Builder{}
	for _, val := range m {
		if val.Type == "text" {
			sb.WriteString(val.Data["text"])
		}
	}
	return sb.String()
}
func (m MessageSegment) Text() string {
	sb := strings.Builder{}
	for _, val := range m.Data {
		sb.WriteString(val)
		break
	}
	return sb.String()
}
func (m MessageSegment) TrimSpaceText() string {
	return strings.TrimSpace(m.Data["text"])
}

func (m MessageSegment) AtId() string {
	sb := strings.Builder{}
	for k, val := range m.Data {
		if k == "qq" || k == "uid" {
			sb.WriteString(val)
			break
		}
	}
	return sb.String()
}
