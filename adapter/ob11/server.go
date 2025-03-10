package ob11

import (
	"encoding/json"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func (c *Config) process(e *Event) {

	switch e.PostType {
	// 消息事件
	case "message", "message_sent":
		c.preprocessMessageEvent(e)
		mess := Message(e.Message).CQString()
		log.Debugf("Event: %+v\n", e)
		switch e.MessageType {
		// 私聊信息
		case "private":
			log.Infof("[ob11] [↓][私聊消息][%v] : %v", e.Sender.String(), e.RawMessage)
			ctx := &rosm.Ctx{
				Being: &rosm.Being{
					GuildID: "",
					GroupID: "-" + tool.Int64ToString(e.Sender.ID),
					RawWord: mess,
					User: &rosm.UserData{
						Name: e.Sender.NickName,
						ID:   tool.Int64ToString(e.Sender.ID),
					},
					MsgID: tool.BytesToString(e.RawMessageID),
				},
				State:   map[string]any{"event": e},
				Message: e.Message,
				Bot:     c,
			}
			ctx.Being.IsAtMe = true
			ctx.RunWord()
		// 群聊信息
		case "group":
			log.Infof("[ob11] [↓][群(%v)消息][%v] : %v", e.GroupID, e.Sender.String(), e.RawMessage)
			ctx := &rosm.Ctx{
				Being: &rosm.Being{
					GuildID: e.ChannelID,
					GroupID: tool.Int64ToString(e.GroupID) + e.GuildID,
					RawWord: mess,
					User: &rosm.UserData{
						Name: e.Sender.NickName,
						ID:   tool.Int64ToString(e.Sender.ID),
					},
					MsgID:  tool.BytesToString(e.RawMessageID),
					ATList: e.AtList,
				},
				State:   map[string]any{"event": e, "reply": e.ReplyMessageID},
				Message: e.Message,
				Bot:     c,
			}
			ctx.Being.IsAtMe = e.IsToMe
			//log.Debugf("%+v\n", ctx.Being)
			//log.Println(ctx.Being.Word)
			ctx.RunWord()
		case "guild":
			log.Infof("[ob11] [↓][频道(%v)(%v-%v)消息][%v] : %v", e.GroupID, e.GuildID, e.ChannelID, e.Sender.String(), e.Message)
			// if e.MessageType == "guild" {
			// 	uid = e.TinyID
			// }
			//strconv.Unquote(helper.BytesToString(event.RawMessageID))
		default:
			log.Warningf("Cannot Parse 'message' event -> %s", e.MessageType)
		}

		// 通知事件
	case "notice":
		preprocessNoticeEvent(e)
		log.Infof("[ob11] [↓][事件 %v]: %v", e.GroupID, e.NoticeType)
		//https://github.com/botuniverse/onebot-11/blob/master/event/notice.md
		ctx := &rosm.Ctx{
			Being: &rosm.Being{
				GuildID: e.ChannelID,
				GroupID: tool.Int64ToString(e.GroupID) + e.GuildID,
				User: &rosm.UserData{
					ID: tool.Int64ToString(e.UserID),
				},
				MsgID: tool.BytesToString(e.RawMessageID),
			},
			State:   map[string]any{"event": e, "notice_type": e.NoticeType},
			Message: e.Message,
			Bot:     c,
		}
		log.Debug(ctx)
		ctx.RunEvent("notice")
	case "request": //好有请求
		ctx := &rosm.Ctx{

			Being: &rosm.Being{
				User: &rosm.UserData{
					ID: tool.Int64ToString(e.UserID),
				},
			},
			State:   map[string]any{"event": e},
			Message: e.Message,
			Bot:     c,
		}
		log.Debug(ctx)
		ctx.RunEvent("request")
	default:
		log.Warningf("Cannot Parse 'message' event -> %s", e.PostType)
	}
}

func (c *Config) processEvent(response []byte, caller APICaller) {
	var event Event
	_ = json.Unmarshal(response, &event)
	event.RawEvent = gjson.Parse(tool.BytesToString(response))
	go c.process(&event)
}

// preprocessNoticeEvent 更新事件
func preprocessNoticeEvent(e *Event) {
	if e.SubType == "poke" || e.SubType == "lucky_king" {
		e.IsToMe = e.TargetID == e.SelfID
	} else {
		e.IsToMe = e.UserID == e.SelfID
	}
}

// preprocessMessageEvent 返回信息事件
func (c *Config) preprocessMessageEvent(e *Event) {
	e.Message = ParseMessage(e.NativeMessage)
	e.IsToMe = false
	// 处理是否是回复消息,at消息
	var f = func() {}
	for i, m := range e.Message {
		if m.Type == "reply" {
			e.ReplyMessageID = m.Data["id"]
			f = func() { e.Message = append(e.Message[:i], e.Message[i+1:]...) }
			continue
		}
		if m.Type == "at" {
			e.AtList = append(e.AtList, m.Data["qq"])
		}
	}
	f()
	if len(e.Message) == 0 {
		return
	}
	//判断at
	if e.Message[0].Type == "at" && tool.Int64ToString(e.SelfID) == e.Message[0].Data["qq"] {
		e.IsToMe = true
		e.Message = e.Message[1:]
		e.AtList = e.AtList[1:]
		return
	}
	if e.Message[0].Type == "text" {
		e.Message[0].Data["text"] = strings.TrimLeft(e.Message[0].Data["text"], " ") // Trim!
		text := e.Message[0].Data["text"]
		for _, nickname := range rosm.GetRosmConfig().BotName {
			if strings.HasPrefix(text, nickname) {
				e.IsToMe = true
				e.Message[0].Data["text"] = strings.TrimLeft(text[len(nickname):], " ") //Trim!
				return
			}
		}
	}
}
