package unimessage

import (
	"fmt"

	"github.com/lianhong2758/RosmBot-MUL/message"
)

// 本实现用于匹配Message的输入
//获取到的内容填入ctx.ResultWorld

// 匹配at,uid为""时,匹配任意at
// 作为一个匹配项
func At(uid string) message.MessageSegment {
	return message.MessageSegment{
		Type: "at",
		Data: message.H{"uid": uid},
	}
}

// 普通文本
func Text(text ...any) message.MessageSegment {
	return message.MessageSegment{
		Type: "text",
		Data: message.H{"text": fmt.Sprint(text...)},
	}
}

// rex
// 作为一个匹配项
func Rex(rs string) message.MessageSegment {
	return message.MessageSegment{
		Type: "rex",
		Data: message.H{"text": rs},
	}
}

// 图片url
// 作为一个匹配项
func Image() message.MessageSegment {
	return message.MessageSegment{
		Type: "image",
		Data: message.H{},
	}
}

// 回复消息
// 作为一个匹配项
func Reply() message.MessageSegment {
	return message.MessageSegment{
		Type: "reply",
		Data: message.H{},
	}
}

// Video 短视频
// 作为一个匹配项
func Video() message.MessageSegment {
	return message.MessageSegment{
		Type: "video",
		Data: message.H{},
	}
}

// 作为一个匹配项
func Any() message.MessageSegment {
	return message.MessageSegment{
		Type: "any",
		Data: message.H{},
	}
}

// 作为一个匹配项
func Other(types string) message.MessageSegment {
	return message.MessageSegment{
		Type: "any",
		Data: message.H{types: ""},
	}
}
