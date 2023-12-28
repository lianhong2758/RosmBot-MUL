package qqmsg

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

func GuildMsgContent(ctx *rosm.Ctx, msg ...message.MessageSegment) *Content {
	cnt := new(Content)
	for _, message := range msg {
		var text string
		if message.Data["text"] != nil {
			text = message.Data["text"].(string)
		}
		switch message.Type {
		default:
			continue
		case "text":
			cnt.Text += text
		case "mentioned_user", "mentioned_robot":
			cnt.Text += `<@!` + message.Data["uid"].(string) + `>`
		case "atall":
			cnt.Text += "@everyone"
		case "imagebyte":
			cnt.Text += text
			if url, _ := web.UpImgByte(message.Data["data"].([]byte)); url != "" {
				cnt.Image = url
			} else {
				cnt.Text += "\n[图片上传失败]\n"
			}
		case "image":
			cnt.Text += text
			if url, _ := web.ImageAnalysis(message.Data["data"].(string)); url != "" {
				cnt.Image = url
			} else {
				cnt.Text += "\n[图片上传失败]\n"
			}
		case "reply":
			cnt.Reference = &ReferenceS{ID: message.Data["ids"].([]string)[0], NeedError: true}
		case "replyuser":
			cnt.Reference = &ReferenceS{ID: ctx.Being.MsgID[0], NeedError: true}
		}
	}
	if ctx.Being.Def["id"] != nil {
		cnt.MsgID = ctx.Being.Def["id"].(string)
	}
	return cnt
}

func GroupMsgContent(ctx *rosm.Ctx, msg ...message.MessageSegment) *Content {
	cnt := new(Content)
	cnt.Types = 0
	for _, message := range msg {
		var text string
		if message.Data["text"] != nil {
			text = message.Data["text"].(string)
		}
		switch message.Type {
		default:
			continue
		case "text":
			cnt.Text += text
			/*
				case "mentioned_user", "mentioned_robot":
						cnt.Types = 5
						cnt.Text += `<@!` + message.Data["uid"].(string) + `>`
					case "atall":
						cnt.Types = 5
						cnt.Text += "@everyone"
			*/
		case "imagebyte":
			cnt.Types = 1
			cnt.Text += text
			if url, _ := web.UpImgByte(message.Data["data"].([]byte)); url != "" {
				cnt.Image = url
			} else {
				cnt.Text += "\n[图片上传失败]\n"
			}
		case "image":
			cnt.Types = 1
			cnt.Text += text
			if url, _ := web.ImageAnalysis(message.Data["data"].(string)); url != "" {
				cnt.Image = url
			} else {
				cnt.Text += "\n[图片上传失败]\n"
			}
		case "reply":
			cnt.Reference = &ReferenceS{ID: message.Data["ids"].([]string)[0], NeedError: true}
		case "replyuser":
			cnt.Reference = &ReferenceS{ID: ctx.Being.MsgID[0], NeedError: true}
		}
	}
	if ctx.Being.Def["id"] != nil {
		cnt.MsgID = ctx.Being.Def["id"].(string)
	}
	return cnt
}
