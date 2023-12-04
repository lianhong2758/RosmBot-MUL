package mys

import (
	"unicode/utf16"

	"github.com/lianhong2758/RosmBot-MUL/message"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

func MakeMsgContent(ctx *rosm.CTX, msg ...message.MessageSegment) (contentInfo any, object string) {
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
		case "bold", "italic", "strikethrough", "underline":
			t := Entities{
				Length: len(utf16.Encode([]rune(text))),
				Offset: len(utf16.Encode([]rune(msgContent.Text))),
				Entity: H{"type": "style", "font_style": message.Type},
			}
			msgContent.Entities = append(msgContent.Entities, t)
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
			/*	case "imagewithtext":
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
					if msgContent.Text == "" {
						msgContent.ImageStr = t
					} else {
						msgContent.Images = append(msgContent.Images, t)
					}*/
		case "imagebyte":
			if url, con := UpImgByte(ctx, message.Data["data"].([]byte)); url != "" {
				t := ImageStr{
					URL:      url,
					Size:     new(Size),
					FileSize: len(message.Data["data"].([]byte)),
				}
				if w := con.Width; w != 0 {
					t.Size.Width = w
				}
				if h := con.Height; h != 0 {
					t.Size.Height = h
				}
				if msgContent.Text == "" {
					msgContent.ImageStr = t
				} else {
					msgContent.Images = append(msgContent.Images, t)
				}
			} else {
				msgContent.Text += "\n[图片上传失败]\n"
			}
		case "image":
			if url, con := ImageAnalysis(ctx, message.Data["data"].(string)); url != "" {
				t := ImageStr{
					URL:  url,
					Size: new(Size),
				}
				if con != nil {
					if w := con.Width; w != 0 {
						t.Size.Width = w
					}
					if h := con.Height; h != 0 {
						t.Size.Height = h
					}
				}
				if msgContent.Text == "" {
					msgContent.ImageStr = t
				} else {
					msgContent.Images = append(msgContent.Images, t)
				}
			} else {
				msgContent.Text += "\n[图片上传失败]\n"
			}
		case "reply":
			id, time := message.Data["ids"].([]string)[0], message.Data["ids"].([]string)[1]
			msgContentInfo["quote"] = H{"original_message_id": id, "original_message_send_time": time, "quoted_message_id": id, "quoted_message_send_time": time}
		case "replyuser":
			msgContentInfo["quote"] = H{"original_message_id": ctx.Being.MsgID[0], "original_message_send_time": tool.Int64(ctx.Being.MsgID[1]), "quoted_message_id": ctx.Being.MsgID[0], "quoted_message_send_time": tool.Int64(ctx.Being.MsgID[1])}
		case "badge":
			t := message.Data["badge"].(BadgeStr)
			msgContent.Badge = &t
		case "view":
			t := message.Data["view"].(PreviewStr)
			msgContent.Preview = &t
		case "custom":
			return message.Data["data"], "MHY:Text"
		case "post":
			return H{"content": H{"post_id": message.Data["id"].(string)}}, "MHY:Post"
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
