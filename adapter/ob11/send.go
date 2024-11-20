package ob11

import (
	"fmt"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"

	"github.com/sirupsen/logrus"
)

func (c *Config) BotSend(ctx *rosm.Ctx, msg ...message.MessageSegment) rosm.H {
	if len(msg) == 0 {
		logrus.Warn("[↑]消息为空")
		return rosm.H{}
	}
	t := MakeMsgContent(ctx, msg...)
	if ctx.Being.GroupID[0:1] != "-" {
		return rosm.H{"id": tool.Int64ToString(SendGroupMessage(ctx, tool.StringToInt64(ctx.Being.GroupID), t)), "code": "0"}
	} else {
		return rosm.H{"id": tool.Int64ToString(SendPrivateMessage(ctx, tool.StringToInt64(ctx.Being.GroupID[1:]), t)), "code": "0"}
	}
}

func (c *Config) BotSendCustom(ctx *rosm.Ctx, Count any) rosm.H {
	if Count == nil {
		logrus.Warn("[↑]消息为空")
		return rosm.H{}
	}
	if c, ok := Count.(string); ok {
		Count = UnescapeCQCodeText(c)
	}
	if ctx.Being.GroupID[0:1] != "-" {
		return rosm.H{"id": tool.Int64ToString(SendGroupMessage(ctx, tool.StringToInt64(ctx.Being.GroupID),
			UnescapeCQCodeText(Count.(string)))), "code": "0"}
	} else {
		return rosm.H{"id": tool.Int64ToString(SendPrivateMessage(ctx, tool.StringToInt64(ctx.Being.GroupID[1:]),
			UnescapeCQCodeText(Count.(string)))), "code": "0"}
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
			if message.Data["uid"] == "all" {
				msg[k].Data = rosm.H{"qq": "all"}
				continue
			}
			msg[k].Data = rosm.H{"qq": message.Data["uid"]}
		case "image":
			msg[k].Data = rosm.H{"file": ImageAnalysis(message.Data["file"])}
		case "replyuser":
			msg[k].Type = "reply"
			msg[k].Data = rosm.H{"id": ctx.Being.MsgID}
		case "link":
			msg[k].Type = "text"
			msg[k].Data = rosm.H{"text": fmt.Sprintf("%s:\n%s", message.Data["text"], message.Data["url"])}
			// case "custom":
			// 	return  UnescapeCQCodeText(msg[k].Data["data"])
		}
	}
	return msg
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
