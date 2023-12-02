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
	var msgContent *qqmsg.Content
	if IsGroup {
		msgContent = qqmsg.GroupMsgContent(ctx, msg...)
		//图片发送改为富文本
		if msgContent.Image != "" {
			if r, err := UpFile(ctx, msgContent.Image, 1); err == nil {
				msgContent.Image = ""
				msgContent.Media = &qqmsg.Media{FileInfo: r.FileINFO}
			}
		}
	} else {
		msgContent = qqmsg.GuildMsgContent(ctx, msg...)
	}
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
