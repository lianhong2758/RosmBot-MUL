package message

import (
	"fmt"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

type H = map[string]any

type Message []MessageSegment

type MessageSegment struct {
	Type string `json:"type"`
	Data H      `json:"data"`
}

// 消息解析
func Text(text ...any) MessageSegment {
	return MessageSegment{
		Type: "text",
		Data: H{"text": fmt.Sprint(text...)},
	}
}

// at用户
func AT(uid int64, name string) MessageSegment {
	name = "@" + name + " "
	return MessageSegment{
		Type: "mentioned_user",
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

// url为图片链接,必须直链,w,h为宽高
func ImageUrlWithText(url string, w, h, size int, text ...any) MessageSegment {
	return MessageSegment{
		Type: "imagewithtext",
		Data: H{
			"text": fmt.Sprint(text...),
			"url":  url,
			"w":    w,
			"h":    h,
			"size": size,
		},
	}
}

// url为图片链接,必须直链,w,h为宽高size大小,不需要项填0
func ImageUrl(url string, w, h, size int) MessageSegment {
	return MessageSegment{
		Type: "image",
		Data: H{
			"url":  url,
			"w":    w,
			"h":    h,
			"size": size,
		},
	}
}

// 发送普通图片
func Image(img []byte) MessageSegment {
	if url, con := web.UpImgByte(img); url != "" {
		return ImageUrl(url, con.Width, con.Height, 0)
	}
	return Text("图片上传失败")
}

// 发送普通图片和文字,text必填
func ImageWithText(img []byte, text ...any) MessageSegment {
	if url, con := web.UpImgByte(img); url != "" {
		return ImageUrlWithText(url, con.Width, con.Height, 0, text...)
	}
	return Text("图片上传失败")
}

// 发送图片文件
func ImageFile(path string) MessageSegment {
	if url, con := web.UpImgfile(path); url != "" {
		return ImageUrl(url, con.Width, con.Height, 0)
	}
	return Text("图片上传失败")
}

// 发送图片文件和文字,text必填
func ImageFileWithText(path string, text ...any) MessageSegment {
	if url, con := web.UpImgfile(path); url != "" {
		return ImageUrlWithText(url, con.Width, con.Height, 0, text...)
	}
	return Text("图片上传失败")
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
