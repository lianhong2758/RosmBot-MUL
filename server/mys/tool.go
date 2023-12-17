package mys

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

// 改变发送的房间id
func ChangeSendRoom(ctx *rosm.CTX, roomid string) {
	ctx.Being.RoomID = roomid
}

// 改变发送的大别野id
func ChangeSendVilla(ctx *rosm.CTX, villaid, roomid string) {
	ctx.Being.RoomID2, ctx.Being.RoomID = villaid, roomid
}

// 是否是主人权限
func OnlyMaster(ctx *rosm.CTX) bool {
	for _, v := range ctx.Bot.(*Config).BotToken.Master {
		if v == ctx.Being.User.ID {
			return true
		}
	}
	return false
}

// 别野房东权限以上
func OnlyOverOwner(ctx *rosm.CTX) bool {
	data, err := GetVillaData(ctx)
	if err != nil {
		log.Errorln("[ctx](", ctx.Being.RoomID2, ")获取别野信息失败:", err)
	}
	return ctx.Being.User.ID == data.Data.Villa.OwnerUID || OnlyMaster(ctx)
}

// 触发消息是否是回复消息
func OnlyReply(ctx *rosm.CTX) bool {
	return ctx.Message.(*MessageContent).Quote.OriginalMessageSendTime != 0 && ctx.Message.(*MessageContent).Quote.QuotedMessageID != ""
}

// 新建bot消息
func NewBot(botid string) rosm.Boter {
	return botMap[botid]
}

// 新建上下文
func NewCTX(botid, roomid, villaid string) *rosm.CTX {
	return &rosm.CTX{
		BotType: "mys",
		Bot:     botMap[botid],
		Being: &rosm.Being{
			RoomID:  roomid,
			RoomID2: villaid,
		},
	}
}
func GetRandBot() *Config {
	for k := range botMap {
		return botMap[k]
	}
	return nil
}

// RangeBot 遍历所有bot实例
func RangeBot(fn func(id string, bot *Config) bool) {
	for k, v := range botMap {
		if !fn(k, v) {
			return
		}
	}
}
