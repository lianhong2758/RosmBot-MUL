package qqmsg

import (

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func GuildMsgContent(ctx *rosm.CTX, msg ...message.MessageSegment) *Content {
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
		case "imagewithtext":
			cnt.Text += text
			cnt.Image = message.Data["url"].(string)
		case "image":
			cnt.Image = message.Data["url"].(string)
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

func GroupMsgContent(ctx *rosm.CTX, msg ...message.MessageSegment) *Content {
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
		case "imagewithtext":
			cnt.Types = 1
			cnt.Text += text
			cnt.Image = message.Data["url"].(string)
		case "image":
			cnt.Types = 1
			cnt.Image = message.Data["url"].(string)
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
