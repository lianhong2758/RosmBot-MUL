package qq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/qq/qqmsg"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	log "github.com/sirupsen/logrus"
)

func (c *Config) BotSend(ctx *rosm.CTX, msg ...message.MessageSegment) any {
	var IsGroup bool = ctx.Being.Def["type"].(string) == "GROUP_AT_MESSAGE_CREATE" || ctx.Being.Def["type"].(string) == "C2C_MESSAGE_CREATE"
	msgContent := makeMsgContent(ctx, IsGroup, msg...)
	data, _ := json.Marshal(msgContent)
	log.Debugln("[send]", tool.BytesToString(data))
	url := ""
	//判断私聊
	if ctx.Being.RoomID == "" {
		if IsGroup {
			url = fmt.Sprintf(urlSendPrivate, ctx.Being.User.ID) //私聊
		} else {
			url = fmt.Sprintf(urlSendGuildPrivate, ctx.Being.User.ID) //频道私聊
		}
	} else {
		if IsGroup {
			url = fmt.Sprintf(urlSendGroup, ctx.Being.RoomID) //群聊
		} else {
			url = fmt.Sprintf(urlSendGuild, ctx.Being.RoomID) //频道
		}
	}
	data, err := web.Web(clientConst, host+url, http.MethodPost, makeHeard(c.access, c.BotToken.AppId), bytes.NewReader(data))
	if err != nil {
		log.Errorln("[send][", host+url, "]", err)
	}
	sendState := new(qqmsg.SendState)
	_ = json.Unmarshal(data, sendState)
	log.Infoln("[send]["+sendState.MsgID+"]", tool.BytesToString(data))
	return sendState
}
func makeMsgContent(ctx *rosm.CTX, IsGroup bool, msg ...message.MessageSegment) *qqmsg.Content {
	cnt := new(qqmsg.Content)
	for _, message := range msg {
		var text string
		if IsGroup {
			cnt.Types = 0
		}
		if message.Data["text"] != nil {
			text = message.Data["text"].(string)
		}
		switch message.Type {
		default:
			continue
		case "text":
			cnt.Text += text
		case "mentioned_user", "mentioned_robot":
			if IsGroup {
				cnt.Types = 5
			}
			cnt.Text += `<@!` + message.Data["uid"].(string) + `>`
		case "atall":
			if IsGroup {
				cnt.Types = 5
			}
			cnt.Text += "@everyone"
		case "imagewithtext":
			if IsGroup {
				cnt.Types = 1
			}
			cnt.Text += text
			cnt.Image = message.Data["url"].(string)
		case "image":
			if IsGroup {
				cnt.Types = 1
			}
			cnt.Image = message.Data["url"].(string)
		case "reply":
			cnt.Reference = &qqmsg.ReferenceS{ID: message.Data["ids"].([]string)[0], NeedError: true}
		}
	}
	if ctx.Being.Def["id"] != nil {
		cnt.MsgID = ctx.Being.Def["id"].(string)
	}
	return cnt
}
