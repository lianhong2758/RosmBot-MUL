package ob11

import (
	"encoding/json"
	"hash/crc64"
	"strconv"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func (c *Config) process(e *zero.Event) {
	mess := e.Message.CQString()
	log.Debug("Message: ", mess)
	switch e.PostType {
	// 消息事件
	case "message":
		switch e.MessageType {
		// 私聊信息
		case "private":
			ctx := &rosm.Ctx{
				BotType: "ob11",
				Being: &rosm.Being{
					RoomID2: "",
					RoomID:  "-" + tool.Int64ToString(e.Sender.ID),
					Word:    mess,
					User: &rosm.UserData{
						Name: e.Sender.NickName,
						ID:   tool.Int64ToString(e.Sender.ID),
						//	PortraitURI: u.User.PortraitURI,
					},
					MsgID: []string{tool.BytesToString(e.RawMessageID)},
					Def:   map[string]any{},
				},
				Message: e,
				Bot:     c,
			}
			ctx.Being.AtMe = true
			ctx.RunWord(ctx.Being.Word)
		// 群聊信息
		case "group":
			ctx := &rosm.Ctx{
				BotType: "ob11",
				Being: &rosm.Being{
					RoomID2: e.ChannelID,
					RoomID:  tool.Int64ToString(e.GroupID) + e.GuildID,
					Word:    mess,
					User: &rosm.UserData{
						Name: e.Sender.NickName,
						ID:   tool.Int64ToString(e.Sender.ID),
						//	PortraitURI: u.User.PortraitURI,
					},
					MsgID: []string{tool.BytesToString(e.RawMessageID)},
					Def:   map[string]any{},
				},
				Message: e,
				Bot:     c,
			}

			ctx.Being.AtMe = e.IsToMe
			e.IsToMe = ctx.Being.AtMe
			//log.Println(ctx.Being.Word)
			ctx.RunWord(ctx.Being.Word)
		case "guild":

		default:
			//	log.Warning(fmt.Sprintf("Cannot Parse 'message' event -> %s", receive))
		}

		// 通知事件
	case "notice":
		ctx := &rosm.Ctx{
			BotType: "ob11",
			Being: &rosm.Being{
				RoomID2: e.ChannelID,
				RoomID:  tool.Int64ToString(e.GroupID) + e.GuildID,
				User:    &rosm.UserData{
					// ID: tool.Int64ToString(e.Sender.ID),
					// Name: e.Sender.NickName,
				},
				Def: map[string]any{},
			},
			Message: e,
			Bot:     c,
		}
		log.Debug(ctx)
		//ctx.RunEvent()
	default:
	}
}

func (c *Config) processEvent() func([]byte, zero.APICaller) {
	return func(response []byte, caller zero.APICaller) {
		var event zero.Event
		_ = json.Unmarshal(response, &event)
		event.RawEvent = gjson.Parse(tool.BytesToString(response))
		//var msgid message.MessageID
		messageID, err := strconv.ParseInt(tool.BytesToString(event.RawMessageID), 10, 64)
		if err == nil {
			event.MessageID = messageID
			//	msgid = message.NewMessageIDFromInteger(messageID)
		} else if event.MessageType == "guild" {
			// 是 guild 消息，进行如下转换以适配非 guild 插件
			// MessageID 填为 string
			event.MessageID, _ = strconv.Unquote(tool.BytesToString(event.RawMessageID))
			// 伪造 GroupID
			crc := crc64.New(crc64.MakeTable(crc64.ISO))
			crc.Write(tool.StringToBytes(event.GuildID))
			crc.Write(tool.StringToBytes(event.ChannelID))
			r := int64(crc.Sum64() & 0x7fff_ffff_ffff_ffff) // 确保为正数
			if r <= 0xffff_ffff {
				r |= 0x1_0000_0000 // 确保不与正常号码重叠
			}
			event.GroupID = r
			// 伪造 UserID
			crc.Reset()
			crc.Write(tool.StringToBytes(event.TinyID))
			r = int64(crc.Sum64() & 0x7fff_ffff_ffff_ffff) // 确保为正数
			if r <= 0xffff_ffff {
				r |= 0x1_0000_0000 // 确保不与正常号码重叠
			}
			event.UserID = r
			if event.Sender != nil {
				event.Sender.ID = r
			}
			//	msgid = message.NewMessageIDFromString(event.MessageID.(string))
		}
		switch event.PostType { // process DetailType
		case "message", "message_sent":
			event.DetailType = event.MessageType
		case "notice":
			event.DetailType = event.NoticeType
			preprocessNoticeEvent(&event)
		case "request":
			event.DetailType = event.RequestType
		}
		if event.PostType == "message" {
			c.preprocessMessageEvent(&event)
		}
		go c.process(&event)
	}
}

// preprocessNoticeEvent 更新事件
func preprocessNoticeEvent(e *zero.Event) {
	if e.SubType == "poke" || e.SubType == "lucky_king" {
		e.IsToMe = e.TargetID == e.SelfID
	} else {
		e.IsToMe = e.UserID == e.SelfID
	}
}

// preprocessMessageEvent 返回信息事件
func (c *Config) preprocessMessageEvent(e *zero.Event) {
	e.Message = message.ParseMessage(e.NativeMessage)

	processAt := func() { // 处理是否at机器人
		e.IsToMe = false
		for i, m := range e.Message {
			if m.Type == "at" {
				qq, _ := strconv.ParseInt(m.Data["qq"], 10, 64)
				if qq == e.SelfID {
					e.IsToMe = true
					e.Message = append(e.Message[:i], e.Message[i+1:]...)
					return
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
	if len(e.Message) > 0 && e.Message[0].Type == "text" { // Trim Again!
		e.Message[0].Data["text"] = strings.TrimLeft(e.Message[0].Data["text"], " ")
	}
}
