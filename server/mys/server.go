package mys

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	vila_bot "github.com/lianhong2758/RosmBot-MUL/server/mys/proto"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

func (cc *Config) MessReceive() func(c *gin.Context) {
	return func(c *gin.Context) {
		body, _ := c.GetRawData()
		c.JSON(200, map[string]any{"message": "", "retcode": 0}) //确认接收
		sign := c.GetHeader("x-rpc-bot_sign")
		if verify(sign, tool.BytesToString(body), cc.BotToken.BotSecretConst, cc.BotToken.BotPubKey) {
			event := new(vila_bot.RobotEvent)
			err := json.Unmarshal(body, event)
			if err != nil {
				log.Errorln("[info]", err)
				return
			}
			cc.process(event)
		}
	}
}

func (c *Config) process(event *vila_bot.RobotEvent) {
	//调用消息处理件,触发中心
	switch event.Type {
	default:
		log.Infoln("[info] (接收未知事件)", event.ExtendData.EventData)
		return
	case 1:
		log.Debugln("[debug] (入群事件)", event.ExtendData.GetJoinVilla().VillaId)
		log.Infof("[info] (入群事件)[%d] %s(%d)", event.ExtendData.GetJoinVilla().VillaId, event.ExtendData.GetJoinVilla().JoinUserNickname, event.ExtendData.GetJoinVilla().JoinUid)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(int(event.Robot.GetVillaId())),
				User:    &rosm.UserData{ID: strconv.Itoa(int(event.ExtendData.GetJoinVilla().JoinUid)), Name: event.ExtendData.GetJoinVilla().JoinUserNickname},
			},
			Bot:     c,
			Message: event,
		}
		ctx.RunEvent(rosm.Join)
	case 3:
		log.Debugln("[debug] (添加bot)", event.ExtendData.GetCreateRobot().VillaId)
		log.Infof("[info] (添加Bot事件)[%d]", event.ExtendData.GetCreateRobot().VillaId)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(int(event.Robot.GetVillaId())),
			},
			Message: event,
			Bot:     c,
		}
		ctx.RunEvent(rosm.Create)
	case 4:
		log.Debugln("[debug] (删除bot)", event.ExtendData.GetDeleteRobot().VillaId)
		log.Infof("[info] (删除Bot事件)[%d]", event.ExtendData.GetDeleteRobot().VillaId)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(int(event.Robot.GetVillaId())),
			},
			Message: event,
			Bot:     c,
		}
		ctx.RunEvent(rosm.Delete)
	case 5:
		log.Debugln("[debug] (接收表态)", event.ExtendData.GetAddQuickEmoticon().Emoticon)
		log.Infof("[info] (表态事件)[%d] %d:%s", event.Robot.GetVillaId(), event.ExtendData.GetAddQuickEmoticon().Uid, event.ExtendData.GetAddQuickEmoticon().Emoticon)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(int(event.Robot.GetVillaId())),
				User:    &rosm.UserData{ID: strconv.Itoa(int(event.ExtendData.GetAddQuickEmoticon().Uid))},
				RoomID:  strconv.Itoa(int(event.ExtendData.GetAddQuickEmoticon().RoomId)),
			},
			Message: event,
			Bot:     c,
		}
		ctx.RunEvent(rosm.Quick)
	case 6:
		//回调审核
		log.Debugln("[debug] (接收审核回调)", event.ExtendData.GetAuditCallback().AuditResult)
		log.Infof("[info] (审核回调)[%s] 审核结果: %d", event.ExtendData.GetAuditCallback().AuditId, event.ExtendData.GetAuditCallback().AuditResult)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(int(event.ExtendData.GetAuditCallback().VillaId)),
				User:    &rosm.UserData{ID: strconv.Itoa(int(event.ExtendData.GetAuditCallback().UserId))},
				RoomID:  strconv.Itoa(int(event.ExtendData.GetAuditCallback().RoomId)),
			},
			Message: event,
			Bot:     c,
		}
		ctx.RunEvent(rosm.Audit)
	case 7:
		log.Debugln("[debug] (接收回溯事件)", event.ExtendData.GetClickMsgComponent().Extra)
		log.Infof("[info] (回溯事件)[%d] %d: %s", event.ExtendData.GetClickMsgComponent().VillaId, event.ExtendData.GetClickMsgComponent().Uid, event.ExtendData.GetClickMsgComponent().Extra)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				Word:    event.ExtendData.GetClickMsgComponent().Extra,
				RoomID2: strconv.Itoa(int(event.ExtendData.GetClickMsgComponent().VillaId)),
				User:    &rosm.UserData{ID: strconv.Itoa(int(event.ExtendData.GetClickMsgComponent().Uid))},
				RoomID:  strconv.Itoa(int(event.ExtendData.GetClickMsgComponent().RoomId)),
				MsgID:   []string{event.ExtendData.GetClickMsgComponent().GetMsgUid(), event.ExtendData.GetClickMsgComponent().GetBotMsgId()},
				Def: map[string]any{
					"extra":        event.ExtendData.GetClickMsgComponent().Extra,
					"component_id": event.ExtendData.GetClickMsgComponent().ComponentId, // 机器人自定义的组件id
					"template_id":  event.ExtendData.GetClickMsgComponent().TemplateId,  // 如果该组件模板为已创建模板，则template_id不为0
				},
			},
			Message: event,
			Bot:     c,
		}
		ctx.RunEvent(rosm.Click)

	case 2:
		log.Debugln("[debug] (接收消息)", event.ExtendData.GetSendMessage().GetContent())
		u := new(MessageContent)
		err := json.Unmarshal([]byte(event.ExtendData.GetSendMessage().GetContent()), u)
		if err != nil {
			log.Errorln("[info]", err)
			return
		}
		log.Infof("[info] (接收消息)[%d] %s:%s", event.Robot.GetVillaId(), u.User.Name, u.Content.Text)
		id, _ := strconv.Atoi(u.User.ID)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(int(event.Robot.GetVillaId())),
				RoomID:  strconv.Itoa(int(event.ExtendData.GetSendMessage().RoomId)),
				User:    &rosm.UserData{Name: u.User.Name, ID: strconv.Itoa(id), PortraitURI: u.User.PortraitURI},
				ATList:  u.MentionedInfo.UserIDList,
				MsgID:   []string{event.ExtendData.GetSendMessage().MsgUid, tool.Int64ToString(event.ExtendData.GetSendMessage().SendAt)},
				Def: map[string]any{
					"Quote": &event.ExtendData.GetSendMessage().QuoteMsg, //type  MessageForQuote
				},
			},
			Message: u,
			Bot:     c,
		}
		ctx.Being.AtMe = true
		//消息处理(切割加去除尾部空格)
		word := strings.TrimSpace(u.Content.Text[len(event.Robot.Template.Name)+2:])
		ctx.RunWord(word)
	}
}
