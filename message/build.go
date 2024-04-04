// ########################################
// 消息组件
// 在Send()中进行发送,不限长度,Custom()除外
// 部分组件可以含有多余的参数
// 适平台要求,可以将部分参数传空

// 不同的server可以有独属于自己的参数
// 请谨慎使用

// 如果有无法识别的组件信息
// server可以选择处理或者丢弃

// 本包下的所有组件server必须实现
// 平台有限制无法实现的除外
// ########################################
package message

import (
	"fmt"
)

type H = map[string]any

type Message []MessageSegment

type MessageSegment struct {
	Type string `json:"type"`
	Data H      `json:"data"`
}

// 普通文本
func Text(text ...any) MessageSegment {
	return MessageSegment{
		Type: "text",
		Data: H{"text": fmt.Sprint(text...)},
	}
}

// at用户
func AT(uid, name string) MessageSegment {
	name = "@" + name + " "
	return MessageSegment{
		Type: "at",
		Data: H{
			"text": name,
			"uid":  uid,
		},
	}
}

// atbot
func ATBot(botid, botname string) MessageSegment {
	botname = "@" + botname + " "
	return MessageSegment{
		Type: "mentioned_robot",
		Data: H{
			"text": botname,
			"uid":  botid,
		},
	}
}

// at all
func ATAll() MessageSegment {
	return MessageSegment{
		Type: "atall",
		Data: H{
			"text": "@全体成员 ",
		},
	}
}

// 发送普通图片
func ImageByte(img []byte) MessageSegment {
	return MessageSegment{
		Type: "imagebyte",
		Data: H{
			"data": img,
		},
	}
}

// 发送图片
// 支持的格式:base64://,file://,url://,consturl://
func Image(data string) MessageSegment {
	return MessageSegment{
		Type: "image",
		Data: H{
			"data": data,
		},
	}
}

// 蓝色跳转链接
func Link(url string, haveToken bool, text ...any) MessageSegment {
	return MessageSegment{
		Type: "link",
		Data: H{
			"text":  fmt.Sprint(text...),
			"url":   url,
			"token": haveToken,
		},
	}
}

// 根据平台增加参数个数,米游社参数为id,time
func ReplyOther(some ...string) MessageSegment {
	return MessageSegment{
		Type: "reply",
		Data: H{
			"ids": some,
		},
	}
}

// 根据平台增加参数个数,米游社参数为id,time
func Reply() MessageSegment {
	return MessageSegment{
		Type: "replyuser",
	}
}

// 自定义全量消息内容
func Custom(messageData any) MessageSegment {
	return MessageSegment{
		Type: "custom",
		Data: H{
			"data": messageData,
		},
	}
}
