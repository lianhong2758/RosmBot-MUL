package ob11

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	zms "github.com/wdvxdr1123/ZeroBot/message"

	"github.com/sirupsen/logrus"
)

func (c *Config) BotSend(ctx *rosm.Ctx, msg ...message.MessageSegment) rosm.H {
	if len(msg) == 0 {
		logrus.Warn("[↑]消息为空")
		return rosm.H{}
	}
	t := MakeMsgContent(ctx, msg...)
	if ctx.Being.RoomID[0:1] != "-" {
		return rosm.H{"id": tool.Int64ToString(SendGroupMessage(ctx, tool.StringToInt64(ctx.Being.RoomID), t)), "code": "0"}
	} else {
		return rosm.H{"id": tool.Int64ToString(SendPrivateMessage(ctx, tool.StringToInt64(ctx.Being.RoomID[1:]), t)), "code": "0"}
	}
}

func (c *Config) BotSendCustom(ctx *rosm.Ctx, Count any) rosm.H {
	if Count == nil {
		logrus.Warn("[↑]消息为空")
		return rosm.H{}
	}
	if ctx.Being.RoomID[0:1] != "-" {
		return rosm.H{"id": tool.Int64ToString(SendGroupMessage(ctx, tool.StringToInt64(ctx.Being.RoomID),
			zms.UnescapeCQCodeText(Count.(string)))), "code": "0"}
	} else {
		return rosm.H{"id": tool.Int64ToString(SendPrivateMessage(ctx, tool.StringToInt64(ctx.Being.RoomID[1:]),
			zms.UnescapeCQCodeText(Count.(string)))), "code": "0"}
	}
}

// 转为符合ob11的map[string]string,但需要再调用custom消息除外
func MakeMsgContent(ctx *rosm.Ctx, msg ...message.MessageSegment) message.Message {
	for k, message := range msg {
		switch message.Type {
		default:
			continue
		case "text", "video", "node", "reply":
			continue
		case "at":
			msg[k].Data = rosm.H{"qq": message.Data["uid"]}
		case "atall":
			msg[k].Type = "at"
			msg[k].Data = rosm.H{"qq": "all"}
		case "image":
			msg[k].Data = rosm.H{"file": ImageAnalysis(message.Data["file"])}
		case "replyuser":
			msg[k].Type = "reply"
			msg[k].Data = rosm.H{"id": ctx.Being.MsgID[0]}
		case "link":
			msg[k].Type = "text"
			msg[k].Data = rosm.H{"text": fmt.Sprintf("%s:\n%s", message.Data["text"], message.Data["url"])}
			// case "custom":
			// 	return zms.UnescapeCQCodeText(msg[k].Data["data"])
		}
	}
	return msg
}

func RosmToZeroMessage(r message.Message) (z zms.Message) {
	return *(*zms.Message)(unsafe.Pointer(&r))
}

// 解析base64等data
func ImageAnalysis(data string) (url string) {
	switch parts := strings.SplitN(data, "://", 2); parts[0] {
	case "base64":
		return data
	case "file":
		return data
	case "url":
		return parts[1]
	case "consturl":
		return parts[1]

	}
	return ""
}

func (c *Config) GetPortraitURI(ctx *rosm.Ctx) string {
	return "http://q4.qlogo.cn/g?b=qq&nk=" + ctx.Being.User.ID + "&s=640"
}
