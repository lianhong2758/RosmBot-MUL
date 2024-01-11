package test

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/qq"
)

func init() {
	en := rosm.Register(rosm.NewRegist("qqtest", "- 测试user", ""))
	en.AddWord("测试user").SetBlock(true).MUL("qq_gulid").Handle(func(ctx *rosm.Ctx) {
		u, err := qq.GetGuildUser(ctx, ctx.Being.User.ID)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text(*u))
	})
	en.AddWord("测试私聊").SetBlock(true).MUL("qq_gulid").Handle(func(ctx *rosm.Ctx) {
		guildid, chanid, err := qq.NewDms(ctx, ctx.Being.User.ID, ctx.Being.RoomID2)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Being.RoomID, ctx.Being.RoomID2 = chanid, guildid
		ctx.Being.Def["type"] = "DIRECT_MESSAGE_CREATE"
		ctx.Send(message.Text("测试的私信消息"))
	})
}
