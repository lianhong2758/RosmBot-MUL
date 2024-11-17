package ob11

import (
	"encoding/json"
	"strconv"
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
		atList := c.preprocessMessageEvent(e)
		mess := Message(e.Message).CQString()
		log.Debug("Message: ", mess)
		switch e.MessageType {
		// 私聊信息
		case "private":
			ctx := &rosm.Ctx{
				Being: &rosm.Being{
					ATList:  atList,
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
			uid := tool.Int64ToString(e.Sender.ID)
			if e.MessageType == "guild" {
				uid = e.TinyID
			}
			ctx := &rosm.Ctx{

				Being: &rosm.Being{
					ATList:  atList,
					GuildID: e.ChannelID,
					GroupID: tool.Int64ToString(e.GroupID) + e.GuildID,
					RawWord: mess,
					User: &rosm.UserData{
						Name: e.Sender.NickName,
						ID:   uid,
					},
					MsgID: tool.BytesToString(e.RawMessageID),
				},
				State:   map[string]any{"event": e},
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
	event.MessageID = tool.BytesToString(event.RawMessageID)
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
func (c *Config) preprocessMessageEvent(e *Event) []string {
	e.Message = ParseMessage(e.NativeMessage)
	atList := []string{}
	processAt := func() { // 处理是否at机器人
		//索引纠正
		var ioffset = 0
		e.IsToMe = false
		for i, m := range e.Message {
			if m.Type == "at" {
				qq, _ := strconv.ParseInt(m.Data["qq"], 10, 64)
				if qq == e.SelfID && !e.IsToMe {
					e.IsToMe = true
					e.Message = append(e.Message[:i+ioffset], e.Message[i+ioffset+1:]...)
					ioffset--
					continue
				}
				if qq != 0 {
					atList = append(atList, m.Data["qq"])
					e.Message = append(e.Message[:i+ioffset], e.Message[i+ioffset+1:]...)
					ioffset--
					continue
				}
			}
		}
		if len(e.Message) == 0 || e.Message[0].Type != "text" {
			return
		}
		first := e.Message[0]
		first.Data["text"] = strings.TrimLeft(first.Data["text"], " ") // Trim!
		text := first.Data["text"]
		if strings.HasPrefix(text, c.Card().BotName) {
			e.IsToMe = true
			first.Data["text"] = text[len(c.Card().BotName):]
			return
		}
	}
	switch {
	case e.DetailType == "group":
		log.Infof("[ob11] [↓][群(%v)消息][%v] : %v", e.GroupID, e.Sender.String(), e.RawMessage)
		processAt()
	case e.DetailType == "guild" && e.SubType == "channel":
		log.Infof("[ob11] [↓][频道(%v)(%v-%v)消息][%v] : %v", e.GroupID, e.GuildID, e.ChannelID, e.Sender.String(), e.Message)
		processAt()
	default:
		processAt()
		e.IsToMe = true // 私聊也判断为at
		log.Infof("[ob11] [↓][私聊消息][%v] : %v", e.Sender.String(), e.RawMessage)
	}
	return atList
}
