package gscore

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RomiChan/websocket"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/send"
	log "github.com/sirupsen/logrus"
)

func (c *GsConfig) NewWebSocket() {
	for {
		conn, _, err := websocket.DefaultDialer.Dial(c.CoreUrl, nil)
		if err == nil {
			log.Info("[gscore]Core连接成功")
			c.conn = conn
			return
		}
		log.Errorln("[gscore]连接到", c.CoreUrl, "失败", err)
		time.Sleep(time.Second * 3)
	}
}
func (c *GsConfig) RecoveWebScoket() {
	Config.on = false
	Config.cancel()
	//启动ws接收
	log.Info("[gscore]Core尝试重连")
	//创建ws
	Config.NewWebSocket()
	//启动ws接收
	var ctxback context.Context
	ctxback, Config.cancel = context.WithCancel(context.Background())
	go ReadAndSendMessage(ctxback, Config.conn)
	Config.on = true
}

func SendWsMessage(SendMessage []byte, Conn *websocket.Conn) error {
	return Conn.WriteMessage(websocket.BinaryMessage, SendMessage)
}

func ReadAndSendMessage(ctxback context.Context, conn *websocket.Conn) {
	for {
		select {
		case <-ctxback.Done(): // 等待上级通知
			return
		default:
		}
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Error("[gscore]read_gscore_message_error", err)
			time.Sleep(time.Second * 5)
			Config.RecoveWebScoket()
			continue
		}
		var RecMessage RecMessageStr
		err = json.Unmarshal(p, &RecMessage)
		if err != nil {
			log.Error("[gscore]解析gscore_send失败", tool.BytesToString(p))
			continue
		}
		if RecMessage.Content[0].Type == "log_INFO" {
			log.Info("[gscore]log:", tool.BytesToString(RecMessage.Content[0].Data))
		} else {
			SendMessage(&RecMessage)
		}
	}
}
func SendMessage(RecMessage *RecMessageStr) {
	var msg []message.MessageSegment
	var nodeMessage []Message
	makeMsg := func(m []Message) {
		for _, v := range m {
			switch v.Type {
			case "text":
				var s string
				json.Unmarshal(v.Data, &s)
				msg = append(msg, message.Text(s))
			case "image":
				var image string
				_ = json.Unmarshal(v.Data, &image)
				msg = append(msg, message.Image(image))
			case "buttons":
				//***
			case "node":
				err := json.Unmarshal(v.Data, &nodeMessage)
				if err != nil {
					log.Errorf("[gscore]解析%v消息失败: %s", v.Type, tool.BytesToString(v.Data))
				}
			default:
				log.Warnf("[gscore]未知消息类型: %v ,data: %s", v.Type, v.Data)
			}
		}
	}
	makeMsg(RecMessage.Content)
	makeMsg(nodeMessage)
	ctx := cache.Get(RecMessage.MsgId)
	if ctx == nil {
		ctx = send.CTXBuild(RecMessage.BotId, "", RecMessage.TargetId)
	}
	if ctx != nil {
		ctx.Send(msg...)
	} else {
		log.Error("[gscore]空指针无法发送消息:", RecMessage)
	}
}

type MessageReceive struct {
	Bot_id      string         `default:"rosmbot" json:"bot_id"`                                          //Bot适配器类型，如onebot
	Bot_self_id string         `default:"" json:"bot_self_id"`                                            //Bot自身的QQ号
	Msg_id      string         `default:"" json:"msg_id"`                                                 //接受的消息id
	User_type   string         `default:"group" enum:"group,direct,channel,sub_channel" json:"user_type"` //消息类型 对应 群聊 私聊 频道 ？(未知)
	Group_id    string         `default:"" json:"group_id"`                                               //当消息类型为群聊消息时,此处应为群号
	User_id     string         `default:"" json:"user_id"`                                                //发送者QQ号
	User_pm     int            `default:"3" json:"user_pm"`                                               //用户权限等级，1为超级管理员，2为群管理/群主，3为普通用户
	Sender      Dictionary     `json:"sender"`
	Content     []WriteMessage `default:"" json:"content"`
}

type RecMessageStr struct {
	BotId      string `json:"bot_id"`
	BotSelfId  string `json:"bot_self_id"`
	MsgId      string `json:"msg_id"`
	TargetType string `json:"target_type"`
	TargetId   string `json:"target_id"`
	Content    []Message
}

type Message struct {
	Type string          `default:"" json:"type"`
	Data json.RawMessage `default:"" json:"data"`
}
type WriteMessage struct {
	Type string `default:"" json:"type"`
	Data any    `default:"" json:"data"`
}
type Dictionary struct {
	Age      int    `json:"age"`
	Area     string `json:"area"`
	Card     string `json:"card"`
	Level    string `json:"level"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Sex      string `json:"sex"`
	Title    string `json:"title"`
	UserID   int    `json:"user_id"`
	Avater   string `json:"avater"`
}

type GSButton struct {
	Text           string   `json:"text"`
	Data           string   `json:"data"`
	PressedText    *string  `json:"pressed_text,omitempty"`
	Style          int      `json:"style"`
	Action         int      `json:"action"`
	Permission     int      `json:"permission"`
	SpecifyRoleIds []string `json:"specify_role_ids"`
	SpecifyUserIds []string `json:"specify_user_ids"`
	UnsupportTips  string   `json:"unsupport_tips"`
}
