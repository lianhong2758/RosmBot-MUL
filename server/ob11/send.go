package ob11

import (
	"encoding/base64"
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
	msg = MakeMsgContent(ctx, msg...)
	if  ctx.Being.RoomID[0:1]!="-" {
		return rosm.H{"state": "", "id": tool.Int64ToString(SendGroupMessage(ctx, tool.StringToInt64(ctx.Being.RoomID), msg)), "code": 0}
	} else {
		return rosm.H{"state": "", "id": tool.Int64ToString(SendPrivateMessage(ctx, tool.StringToInt64(ctx.Being.RoomID[1:]), msg)), "code": 0}
	}
}

func MakeMsgContent(ctx *rosm.Ctx, msg ...message.MessageSegment)  message.Message {
	for k, message := range msg {
		switch message.Type {
		default:
			continue
		case "text","video":
			continue
			
		case "at":
			msg[k].Data = rosm.H{"qq": message.Data["uid"]}
		case "atall":
			msg[k].Type = "at"
			msg[k].Data = rosm.H{"qq": "all"}
		case "imagebyte":
			msg[k].Type = "image"
			msg[k].Data = rosm.H{"file": "base64://" + base64.StdEncoding.EncodeToString(message.Data["data"].([]byte))}
		case "image":
			msg[k].Data = rosm.H{"file": ImageAnalysis(message.Data["data"].(string))}
		case "reply":
			msg[k].Type = "reply"
			msg[k].Data = rosm.H{"id": message.Data["ids"].([]string)[0]}
		case "replyuser":
			msg[k].Type = "reply"
			msg[k].Data = rosm.H{"id": ctx.Being.MsgID[0] }
		case "link":
			msg[k].Type="text"
			msg[k].Data=rosm.H{"text": fmt.Sprintf( "%s:\n%s",message.Data["text"].(string),message.Data["url"].(string))}
		case "custom":
			return ParseMessageFromString( msg[k].Data["data"].(string))
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
	return"http://q4.qlogo.cn/g?b=qq&nk=" + ctx.Being.User.ID + "&s=640"
}
