package unimessage

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
)

//本实现用于匹配Message的输入
//匹配at,uid为""时,只匹配types,uid在atlist里面获取
func At(types string, uid string) message.MessageSegment {
	return message.MessageSegment{
		Type: "at",
		Data: message.H{types: uid},
	}
}

