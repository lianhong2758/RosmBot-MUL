package test

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func init() {
	en := rosm.Register(rosm.NewRegist("@回复", "- @机器人", ""))
	en.AddWord("").SetBlock(true).Rule(rosm.OnlyAtMe()).Handle(func(ctx *rosm.CTX) {
		ctx.Send(message.Text(ctx.Bot.Name(), "不在呢~"))
	})
}
