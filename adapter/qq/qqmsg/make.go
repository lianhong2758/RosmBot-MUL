package qqmsg

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

func GuildMsgContent(ctx *rosm.Ctx, msg ...message.MessageSegment) *Content {
	cnt := new(Content)
	for _, message := range msg {
		var text string
		if message.Data["text"] != "" {
			text = message.Data["text"]
		}
		switch message.Type {
		default:
			continue
		case "text":
			cnt.Text += text
		case "at", "mentioned_robot":
			cnt.Text += `<@!` + message.Data["uid"] + `>`
		case "atall":
			cnt.Text += "@everyone"

		case "image":
			cnt.Text += text
			if url, _ := web.ImageAnalysis(message.Data["file"]); url != "" {
				cnt.Image = url
			} else {
				cnt.Text += "\n[图片上传失败]\n"
			}
		case "reply":
			cnt.Reference = &ReferenceS{ID: message.Data["id"], NeedError: true}
		case "replyuser":
			cnt.Reference = &ReferenceS{ID: ctx.Being.MsgID, NeedError: true}
		}
	}
	cnt.Text = tool.HideURL(cnt.Text)
	if ctx.State["id"] != nil {
		cnt.MsgID = ctx.State["id"].(string)
	}
	return cnt
}

func GroupMsgContent(ctx *rosm.Ctx, msg ...message.MessageSegment) *Content {
	cnt := new(Content)
	cnt.Types = 0
	for _, message := range msg {
		var text string
		if message.Data["text"] != "" {
			text = message.Data["text"]
		}
		switch message.Type {
		default:
			continue
		case "text":
			cnt.Text += text
		case "image":
			cnt.Types = 1
			cnt.Text += text
			if url, _ := web.ImageAnalysis(message.Data["file"]); url != "" {
				cnt.Image = url
			} else {
				cnt.Text += "\n[图片上传失败]\n"
			}
		case "reply":
			cnt.Reference = &ReferenceS{ID: message.Data["id"], NeedError: true}
		case "replyuser":
			cnt.Reference = &ReferenceS{ID: ctx.Being.MsgID, NeedError: true}
		case "markdown":
			var parm []KV = nil
			if t, ok := ctx.State["kv"]; ok {
				parm = t.([]KV)
			}
			cnt.MarkDown = &MarkDownS{Content: message.Data["content"], ID: message.Data["id"], Params: parm}
		}
	}
	cnt.Text = tool.HideURL(cnt.Text)
	if ctx.State["id"] != nil {
		cnt.MsgID = ctx.State["id"].(string)
	}
	return cnt
}
