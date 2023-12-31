package gscore

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/RomiChan/websocket"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
	"github.com/lianhong2758/RosmBot-MUL/tool"
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
	//发送信息
	var msg []message.MessageSegment
	var p *mys.InfoContent
	for _, v := range RecMessage.Content {
		switch v.Type {
		case "text":
			msg = append(msg, message.Text(tool.BytesToString(v.Data)))
		case "image":
			var image string
			_ = json.Unmarshal(v.Data, &image)
			msg = append(msg, message.Image(image))
		case "buttons":
			var buttons [][]GSButton
			if v.Data[0] != v.Data[1] {
				//二级目录
				v.Data = append(append([]byte{91}, v.Data...), 93)
			}
			err := json.Unmarshal(v.Data, &buttons)
			if err != nil {
				log.Errorf("[gscore]解析%v消息失败: %v", v.Type, tool.BytesToString(v.Data))
			}
			if RecMessage.BotId == "mys" {
				p = mys.NewPanel()
				for l, buttonArry := range buttons {
					for i, button := range buttonArry {
						if RecMessage.BotId == "mys" {
							switch len([]rune(button.Text)) {
							case 1, 2:
								p.Small(i == 0, &mys.Component{
									ID:           strconv.Itoa(l) + strconv.Itoa(i),
									Text:         button.Text,
									Type:         1,
									CType:        2,
									InputContent: button.Data,
									Extra:        "",
								})
							case 3, 4:
								p.Mid(i == 0, &mys.Component{
									ID:           strconv.Itoa(l) + strconv.Itoa(i),
									Text:         button.Text,
									Type:         1,
									CType:        2,
									InputContent: button.Data,
									Extra:        "",
								})
							default:
								p.Big(i == 0, &mys.Component{
									ID:           strconv.Itoa(l) + strconv.Itoa(i),
									Text:         button.Text,
									Type:         1,
									CType:        2,
									InputContent: button.Data,
									Extra:        "",
								})
							}

						}
					}
				}
			}
		case "node":
			var m []Message
			err := json.Unmarshal(v.Data, &m)
			if err != nil {
				log.Errorf("[gscore]解析%v消息失败: %v", v.Type, tool.BytesToString(v.Data))
			}
			RecMessage.Content = append(RecMessage.Content, m...)
		}
	}
	ctx := cache.Get(RecMessage.MsgId)
	if ctx == nil {
		switch RecMessage.BotId {
		case "mys":
			room, villa := tool.SplitPadString(RecMessage.TargetId)
			ctx = mys.NewCTX(RecMessage.BotSelfId, room, villa)
		}
	}
	if p != nil {
		p.TextBuild(ctx, msg...)
		if p.Content.Text == "" {
			p.TextBuild(ctx, append([]message.MessageSegment{message.Text("喵~")}, msg...)...)
		}
		msg = []message.MessageSegment{message.Custom(p)}
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
