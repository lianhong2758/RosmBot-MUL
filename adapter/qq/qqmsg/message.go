package qqmsg

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

// 发送普通图片
func ImageByte(ctx *rosm.Ctx, content, id string, params []KV) message.MessageSegment {
	ctx.State["kv"] = params
	return message.MessageSegment{
		Type: "markdown",
		Data: rosm.H{
			"content": content,
			"id":      id,
		},
	}
}
