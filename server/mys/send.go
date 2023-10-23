package mys

import (
	"bytes"
	"encoding/json"
	"net/http"
	"unicode/utf16"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	. "github.com/lianhong2758/RosmBot-MUL/server/mys/mysmsg"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	log "github.com/sirupsen/logrus"
)

const (
	sendMessage = "https://bbs-api.miyoushe.com/vila/api/bot/platform/sendMessage"
)

func (c *Config) BotSend(ctx *rosm.CTX, msg ...message.MessageSegment) any {
	msgContentInfo, objectStr := makeMsgContent(ctx, msg...)
	contentStr, _ := json.Marshal(msgContentInfo)
	data, _ := json.Marshal(H{"room_id": ctx.Being.RoomID, "object_name": objectStr, "msg_content": tool.BytesToString(contentStr)})
	log.Debugln("[send]", tool.BytesToString(data))
	data, err := web.Web(&http.Client{}, sendMessage, http.MethodPost, makeHeard(ctx), bytes.NewReader(data))
	if err != nil {
		log.Errorln("[send]", err)
	}
	sendState := new(SendState)
	_ = json.Unmarshal(data, sendState)
	log.Infoln("[send]["+sendState.Message+"]", tool.BytesToString(contentStr))
	log.Debugln("[send]", tool.BytesToString(data))
	return sendState
}

func makeMsgContent(ctx *rosm.CTX, msg ...message.MessageSegment) (content any, object string) {
	msgContent := new(Content)
	msgContentInfo := H{}
	for _, message := range msg {
		var text string
		if message.Data["text"] != nil {
			text = message.Data["text"].(string)
		}
		switch message.Type {
		default:
			continue
		case "text":
			msgContent.Text += text
		case "link":
			t := Entities{
				Length: len(utf16.Encode([]rune(text))),
				Offset: len(utf16.Encode([]rune(msgContent.Text))),
				Entity: H{"type": "link", "url": message.Data["url"].(string), "requires_bot_access_token": message.Data["token"].(bool)},
			}
			msgContent.Entities = append(msgContent.Entities, t)
			msgContent.Text += text
		case "villa_room_link":
			t := Entities{
				Length: len(utf16.Encode([]rune(text))),
				Offset: len(utf16.Encode([]rune(msgContent.Text))),
				Entity: H{"type": "villa_room_link", "villa_id": message.Data["villa"].(string), "room_id": message.Data["room"].(string)},
			}
			msgContent.Entities = append(msgContent.Entities, t)
			msgContent.Text += text
		case "mentioned_user":
			t := Entities{
				Length: len(utf16.Encode([]rune(text))),
				Offset: len(utf16.Encode([]rune(msgContent.Text))),
				Entity: H{"type": "mentioned_user", "user_id": message.Data["uid"].(string)},
			}
			msgContent.Entities = append(msgContent.Entities, t)
			msgContent.Text += text
			//特殊实现
			otherUID := []string{}
			if msgContentInfo["mentionedInfo"] != nil {
				otherUID = msgContentInfo["mentionedInfo"].(MentionedInfoStr).UserIDList
			}
			otherUID = append(otherUID, message.Data["uid"].(string))
			msgContentInfo["mentionedInfo"] = MentionedInfoStr{Type: 2, UserIDList: otherUID}
		case "mentioned_robot":
			t := Entities{
				Length: len(utf16.Encode([]rune(text))),
				Offset: len(utf16.Encode([]rune(msgContent.Text))),
				Entity: H{"type": "mentioned_user", "user_id": message.Data["uid"].(string)},
			}
			msgContent.Entities = append(msgContent.Entities, t)
			msgContent.Text += text
			//特殊实现
			otherUID := []string{}
			if msgContentInfo["mentionedInfo"] != nil {
				otherUID = msgContentInfo["mentionedInfo"].(MentionedInfoStr).UserIDList
			}
			otherUID = append(otherUID, message.Data["uid"].(string))
			msgContentInfo["mentionedInfo"] = MentionedInfoStr{Type: 2, UserIDList: otherUID}
		case "atall":
			t := Entities{
				Length: len(utf16.Encode([]rune("@全体成员 "))),
				Offset: len(utf16.Encode([]rune(msgContent.Text))),
				Entity: H{"type": "mention_all"},
			}
			msgContent.Entities = append(msgContent.Entities, t)
			msgContent.Text += text
			msgContentInfo["mentionedInfo"] = MentionedInfoStr{Type: 1}
		case "imagewithtext":
			msgContent.Text += text
			t := ImageStr{
				URL:  message.Data["url"].(string),
				Size: new(Size),
			}
			if w := message.Data["w"].(int); w != 0 {
				t.Size.Width = w
			}
			if h := message.Data["h"].(int); h != 0 {
				t.Size.Height = h
			}
			if s := message.Data["size"].(int); s != 0 {
				t.Size.Height = s
			}
			msgContent.Images = append(msgContent.Images, t)
		case "image":
			t := ImageStr{
				URL:  message.Data["url"].(string),
				Size: new(Size),
			}
			if w := message.Data["w"].(int); w != 0 {
				t.Size.Width = w
			}
			if h := message.Data["h"].(int); h != 0 {
				t.Size.Height = h
			}
			if s := message.Data["size"].(int); s != 0 {
				t.Size.Height = s
			}
			msgContent.ImageStr = t
		case "reply":
			id, time := message.Data["ids"].([]string)[0], message.Data["ids"].([]string)[1]
			msgContentInfo["quote"] = H{"original_message_id": id, "original_message_send_time": time, "quoted_message_id": id, "quoted_message_send_time": time}
		case "badge":
			t := message.Data["badge"].(BadgeStr)
			msgContent.Badge = &t
		case "view":
			t := message.Data["view"].(PreviewStr)
			msgContent.Preview = &t
		case "my":
			return message.Data["my"], "MHY:Text"

		}
	}
	var objectStr string
	if msgContent.URL == "" {
		objectStr = "MHY:Text"
	} else {
		objectStr = "MHY:Image"
	}
	msgContentInfo["content"] = msgContent
	return &msgContentInfo, objectStr
}
func makeHeard(ctx *rosm.CTX) func(req *http.Request) {
	return func(req *http.Request) {
		req.Header.Add("x-rpc-bot_id", ctx.Bot.(*Config).BotToken.BotID)
		req.Header.Add("x-rpc-bot_secret", ctx.Bot.(*Config).BotToken.BotSecret)
		req.Header.Add("x-rpc-bot_villa_id", ctx.Being.RoomID2)
		req.Header.Add("Content-Type", "application/json")
	}
}
