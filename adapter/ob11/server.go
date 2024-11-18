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
		log.Debug("Message: ", mess)
		switch e.MessageType {
		// 私聊信息
		case "private":
			ctx := &rosm.Ctx{
				Being: &rosm.Being{
					GuildID: "",
					GroupID: "-" + tool.Int64ToString(e.Sender.ID),
					RawWord: mess,
					User: &rosm.UserData{
						Name: e.Sender.NickName,
						ID:   tool.Int64ToString(e.Sender.ID),
					},
					MsgID: e.MessageID,
				},
				State:   map[string]any{"event": e},
				Message: e.Message,
				Bot:     c,
			}
			ctx.Being.IsAtMe = true
			ctx.RunWord()
		// 群聊信息
		case "group":
			uid := tool.Int64ToString(e.Sender.ID)
			if e.MessageType == "guild" {
				uid = e.TinyID
			}

			ctx := &rosm.Ctx{
				Being: &rosm.Being{
					GuildID: e.ChannelID,
					GroupID: tool.Int64ToString(e.GroupID) + e.GuildID,
					RawWord: mess,
					User: &rosm.UserData{
						Name: e.Sender.NickName,
						ID:   uid,
					},
					MsgID: e.MessageID,
				},
				State:   map[string]any{"event": e, "reply": ""},
				Message: e.Message,
				Bot:     c,
			}

			ctx.Being.IsAtMe = e.IsToMe
			e.IsToMe = ctx.Being.IsAtMe
			//log.Println(ctx.Being.Word)
			ctx.RunWord()
		case "guild":

		default:
			log.Warningf("Cannot Parse 'message' event -> %s", e.MessageType)
		}

		// 通知事件
	case "notice":
		preprocessNoticeEvent(e)
		//https://github.com/botuniverse/onebot-11/blob/master/event/notice.md
		ctx := &rosm.Ctx{

			Being: &rosm.Being{
				GuildID: e.ChannelID,
				GroupID: tool.Int64ToString(e.GroupID) + e.GuildID,
				User: &rosm.UserData{
					ID: tool.Int64ToString(e.UserID),
				},
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
	e.Message.Reduce()
	e.IsToMe = false
	processReply := func() { // 处理是否是回复消息
		//索引纠正
		for i, m := range e.Message {
			if m.Type == "reply" {
				e.ReplyMessageID = m.Data["id"]
				e.Message = append(e.Message[:i], e.Message[i+1:]...)
				return
			}
		}

		if len(e.Message) == 0 {
			return
		}

		//判断at
		if e.Message[0].Type == "at" && tool.Int64ToString(e.SelfID) == e.Message[0].Data["id"] {
			e.IsToMe = true
			e.Message = e.Message[1:]
		}
		if e.Message[0].Type == "text" {
			e.Message[0].Data["text"] = strings.TrimLeft(e.Message[0].Data["text"], " ") // Trim!
			text := e.Message[0].Data["text"]
			for _, nickname := range rosm.GetRosmConfig().BotName {
				if strings.HasPrefix(text, nickname) {
					e.IsToMe = true
					e.Message[0].Data["text"] =  strings.TrimLeft(text[len(nickname):]," ")//Trim!
					return
				}
			}
		}

	}
	switch {
	case e.DetailType == "group":
		log.Infof("[ob11] [↓][群(%v)消息][%v] : %v", e.GroupID, e.Sender.String(), e.RawMessage)
		processReply()
	case e.DetailType == "guild" && e.SubType == "channel":
		log.Infof("[ob11] [↓][频道(%v)(%v-%v)消息][%v] : %v", e.GroupID, e.GuildID, e.ChannelID, e.Sender.String(), e.Message)
		processReply()
	default:
		processReply()
		e.IsToMe = true // 私聊也判断为at
		log.Infof("[ob11] [↓][私聊消息][%v] : %v", e.Sender.String(), e.RawMessage)
	}
}
