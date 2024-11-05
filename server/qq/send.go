package qq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FloatTech/ttl"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/qq/qqmsg"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	log "github.com/sirupsen/logrus"
)

var seqcache = ttl.NewCache[string, int](time.Minute * 5)

func (c *Config) BotSend(ctx *rosm.Ctx, msg ...message.MessageSegment) rosm.H {
	var IsGroup bool = ctx.Being.Def["type"].(string) == "GROUP_AT_MESSAGE_CREATE" || ctx.Being.Def["type"].(string) == "C2C_MESSAGE_CREATE"
	var msgContent *qqmsg.Content
	if IsGroup {
		msgContent = qqmsg.GroupMsgContent(ctx, msg...)
		//图片发送改为富文本
		if msgContent.Image != "" {
			if r, err := UpFile(ctx, msgContent.Image, 1); err == nil {
				//	c.SendMedia(ctx, r.FileINFO)
				msgContent.Types = 7
				msgContent.Media = &qqmsg.Media{FileInfo: r.FileINFO}
			} else {
				msgContent.Text += "\n[图片上传失败] ERROR:" + err.Error()
				msgContent.Types = 0
			}
			msgContent.Image = ""
		}
	} else {
		msgContent = qqmsg.GuildMsgContent(ctx, msg...)
	}
	seq := seqcache.Get(msgContent.MsgID)
	seq++
	seqcache.Set(msgContent.MsgID, seq)
	msgContent.MsgSeq = seq
	data, _ := json.Marshal(msgContent)
	url := ""
	//判断私聊
	switch ctx.Being.Def["type"].(string) {
	case "GROUP_AT_MESSAGE_CREATE":
		url = fmt.Sprintf(urlSendGroup, ctx.Being.RoomID) //群聊
	case "C2C_MESSAGE_CREATE":
		url = fmt.Sprintf(urlSendPrivate, ctx.Being.User.ID) //私聊
	case "AT_MESSAGE_CREATE", "MESSAGE_CREATE":
		url = fmt.Sprintf(urlSendGuild, ctx.Being.RoomID) //频道
	case "DIRECT_MESSAGE_CREATE":
		url = fmt.Sprintf(urlSendGuildPrivate, ctx.Being.RoomID2) //频道私聊
	}
	log.Infoln("[↑]["+url+"]", tool.BytesToString(data))
	data, err := web.Web(clientConst, host+url, http.MethodPost, makeHeard(c.access, c.BotToken.AppId), bytes.NewReader(data))
	if err != nil {
		log.Errorln("[send][", host+url, "]", err)
	}
	log.Debugln("[send-result]", tool.BytesToString(data))
	sendState := new(qqmsg.SendState)
	_ = json.Unmarshal(data, sendState)
	return rosm.H{"id": sendState.MsgID, "code": func(b bool) string {
		if b {
			return "1"
		}
		return "0"
	}(err != nil)}
}

func (c *Config) BotSendCustom(ctx *rosm.Ctx, count any) rosm.H {
	url := ""
	//判断私聊
	switch ctx.Being.Def["type"].(string) {
	case "GROUP_AT_MESSAGE_CREATE":
		url = fmt.Sprintf(urlSendGroup, ctx.Being.RoomID) //群聊
	case "C2C_MESSAGE_CREATE":
		url = fmt.Sprintf(urlSendPrivate, ctx.Being.User.ID) //私聊
	case "AT_MESSAGE_CREATE", "MESSAGE_CREATE":
		url = fmt.Sprintf(urlSendGuild, ctx.Being.RoomID) //频道
	case "DIRECT_MESSAGE_CREATE":
		url = fmt.Sprintf(urlSendGuildPrivate, ctx.Being.RoomID2) //频道私聊
	}
	log.Infoln("[qq] [↑]["+url+"]",count.(string))
	data, err := web.Web(clientConst, host+url, http.MethodPost, makeHeard(c.access, c.BotToken.AppId), bytes.NewReader(tool.StringToBytes(count.(string))))
	if err != nil {
		log.Errorln("[send][", host+url, "]", err)
	}
	log.Debugln("[send-result]", tool.BytesToString(data))
	sendState := new(qqmsg.SendState)
	_ = json.Unmarshal(data, sendState)
	return rosm.H{"id": sendState.MsgID, "code": func(b bool) string {
		if b {
			return "1"
		}
		return "0"
	}(err != nil)}
}

func (c *Config) GetPortraitURI(ctx *rosm.Ctx) string {
	if r, ok := ctx.Message.(*RawGuildMessage); ok {
		return r.Author.Avatar
	}
	return ""
}
