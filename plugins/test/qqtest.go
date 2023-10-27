package test

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/qq"
)

func init() {
	en := rosm.Register(rosm.NewRegist("qqtest", "- 测试user", ""))
	en.AddWord("测试user").SetBlock(true).MUL("qq").Handle(func(ctx *rosm.CTX) {
		u, err := qq.GetGuildUser(ctx, ctx.Being.User.ID)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text(*u))
	})
}
