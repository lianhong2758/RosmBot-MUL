package mys

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

func (cc *Config) MessReceive() func(c *gin.Context) {
	return func(c *gin.Context) {
		body, _ := c.GetRawData()
		c.JSON(200, map[string]any{"message": "", "retcode": 0}) //确认接收
		sign := c.GetHeader("x-rpc-bot_sign")
		if verify(sign, helper.BytesToString(body), cc.BotToken.BotSecretConst, cc.BotToken.BotPubKey) {
			cc.process(body)
		}
	}
}

func (c *Config) process(body []byte) {
	info := new(InfoSTR)
	err := json.Unmarshal(body, info)
	if err != nil {
		log.Errorln("[info]", err)
		return
	}
	//调用消息处理件,触发中心
	switch info.Event.Type {
	default:
		log.Infoln("[info] (接收未知事件)", info.Event.ExtendData.EventData)
		return
	case 1:
		log.Debugln("[debug] (入群事件)", info.Event.ExtendData.EventData.JoinVilla)
		log.Infof("[info] (入群事件)[%d] %s(%d)", info.Event.Robot.VillaID, info.Event.ExtendData.EventData.JoinVilla.JoinUserNickname, info.Event.ExtendData.EventData.JoinVilla.JoinUID)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(info.Event.Robot.VillaID),
				User:    &rosm.UserData{ID: strconv.Itoa(info.Event.ExtendData.EventData.JoinVilla.JoinUID), Name: info.Event.ExtendData.EventData.JoinVilla.JoinUserNickname},
			},
			Bot:     c,
			Message: info,
		}
		ctx.RunEvent(rosm.Join)
	case 3:
		log.Debugln("[debug] (添加bot)", info.Event.ExtendData.EventData.CreateRobot)
		log.Infof("[info] (添加Bot事件)[%d]", info.Event.Robot.VillaID)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(info.Event.Robot.VillaID),
			},
			Message: info,
			Bot:     c,
		}
		ctx.RunEvent(rosm.Create)
	case 4:
		log.Debugln("[debug] (删除bot)", info.Event.ExtendData.EventData.DeleteRobot)
		log.Infof("[info] (删除Bot事件)[%d]", info.Event.Robot.VillaID)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(info.Event.Robot.VillaID),
			},
			Message: info,
			Bot:     c,
		}
		ctx.RunEvent(rosm.Delete)
	case 5:
		log.Debugln("[debug] (接收表态)", info.Event.ExtendData.EventData.AddQuickEmoticon)
		log.Infof("[info] (表态事件)[%d] %d:%s", info.Event.Robot.VillaID, info.Event.ExtendData.EventData.AddQuickEmoticon.UID, info.Event.ExtendData.EventData.AddQuickEmoticon.Emoticon)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(info.Event.Robot.VillaID),
				User:    &rosm.UserData{ID: strconv.Itoa(info.Event.ExtendData.EventData.AddQuickEmoticon.UID)},
				RoomID:  strconv.Itoa(info.Event.ExtendData.EventData.AddQuickEmoticon.RoomID),
			},
			Message: info,
			Bot:     c,
		}
		emoticonNext(ctx)
		ctx.RunEvent(rosm.Quick)
	//case 6:
	//回调审核
	case 2:
		log.Debugln("[debug] (接收消息)", info.Event.ExtendData.EventData.SendMessage.Content)
		u := new(MessageContent)
		err = json.Unmarshal([]byte(info.Event.ExtendData.EventData.SendMessage.Content), u)
		if err != nil {
			log.Errorln("[info]", err)
			return
		}
		log.Infof("[info] (接收消息)[%d] %s:%s", info.Event.Robot.VillaID, u.User.Name, u.Content.Text)
		id, _ := strconv.Atoi(u.User.ID)
		ctx := &rosm.CTX{
			BotType: "mys",
			Being: &rosm.Being{
				RoomID2: strconv.Itoa(info.Event.Robot.VillaID),
				RoomID:  strconv.Itoa(info.Event.ExtendData.EventData.SendMessage.RoomID),
				User:    &rosm.UserData{Name: u.User.Name, ID: strconv.Itoa(id), PortraitURI: u.User.PortraitURI},
				ATList:  u.MentionedInfo.UserIDList,
			},
			Message: u,
			Bot:     c,
		}
		ctx.Being.AtMe = true
		//消息处理(切割加去除尾部空格)
		word := strings.TrimSpace(u.Content.Text[len(info.Event.Robot.Template.Name)+2:])
		ctx.RunWord(word)
	}
}
