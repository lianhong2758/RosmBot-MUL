// ########################################
// 消息组件
// 在Send()中进行发送,不限长度
//BotSendCustom也需要实现,用于自定义一个特殊的消息发送,一般用于直接发送消息体
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
	"encoding/base64"
	"fmt"
)

type H = map[string]string

type Message  []MessageSegment

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
func AT(uid string) MessageSegment {
	return MessageSegment{
		Type: "at",
		Data: H{
			"uid": uid,
		},
	}
}

// at all
func ATAll() MessageSegment {
	return AT("all")
}

// 发送普通图片
func ImageByte(img []byte) MessageSegment {
	return MessageSegment{
		Type: "image",
		Data: H{
			"file": "base64://" + base64.StdEncoding.EncodeToString(img),
		},
	}
}

// 发送图片
// 支持的格式:base64://,file://,url://,consturl://
func Image(data string) MessageSegment {
	return MessageSegment{
		Type: "image",
		Data: H{
			"file": data,
		},
	}
}

// 蓝色跳转链接
func Link(url string, text ...any) MessageSegment {
	return MessageSegment{
		Type: "link",
		Data: H{
			"text": fmt.Sprint(text...),
			"url":  url,
		},
	}
}

// 回复其余人
func ReplyOther(id string) MessageSegment {
	return MessageSegment{
		Type: "reply",
		Data: H{
			"id": id,
		},
	}
}

// 回复消息
func Reply() MessageSegment {
	return MessageSegment{
		Type: "replyuser",
	}
}

// Video 短视频
func Video(file string) MessageSegment {
	return MessageSegment{
		Type: "video",
		Data: H{
			"file": file,
		},
	}
}
