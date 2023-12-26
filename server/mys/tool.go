package mys

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

// 改变发送的房间id
func ChangeSendRoom(ctx *rosm.Ctx, roomid string) {
	ctx.Being.RoomID = roomid
}

// 改变发送的大别野id
func ChangeSendVilla(ctx *rosm.Ctx, villaid, roomid string) {
	ctx.Being.RoomID2, ctx.Being.RoomID = villaid, roomid
}

// 新建bot消息
func NewBot(botid string) rosm.Boter {
	return botMap[botid]
}

// 新建上下文
func NewCTX(botid, roomid, villaid string) *rosm.Ctx {
	return &rosm.Ctx{
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
