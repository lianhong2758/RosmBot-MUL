package mys

import (
	"strconv"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

// 改变发送的房间id
func ChangeSendRoom(ctx *rosm.CTX, roomid int64) {
	ctx.Being.RoomID = roomid
}

// 改变发送的大别野id
func ChangeSendVilla(ctx *rosm.CTX, villaid, roomid int64) {
	ctx.Being.RoomID2, ctx.Being.RoomID = villaid, roomid
}

// 是否是主人权限
func IsMaster(userID string) bool {
	for _, v := range MYSconfig.BotToken.Master {
		if v == userID {
			return true
		}
	}
	return false
}

func OnlyMaster(ctx *rosm.CTX) bool {
	return IsMaster(strconv.Itoa(int(ctx.Being.User.ID)))
}

// 别野房东权限以上
func OnlyOverOwner(ctx *rosm.CTX) bool {
	data, err := GetVillaData(ctx)
	if err != nil {
		log.Errorln("[ctx](", ctx.Being.RoomID2, ")获取别野信息失败:", err)
	}
	return strconv.Itoa(int(ctx.Being.User.ID)) == data.Data.Villa.OwnerUID || OnlyMaster(ctx)
}

// 触发消息是否是回复消息
func OnlyReply(ctx *rosm.CTX) bool {
	return ctx.Message.(*MessageContent).Quote.OriginalMessageSendTime != 0 && ctx.Message.(*MessageContent).Quote.QuotedMessageID != ""
}
