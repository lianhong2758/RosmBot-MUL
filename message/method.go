package message

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
