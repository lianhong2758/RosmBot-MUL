package test

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "@回复",
		Help: "- @机器人",
	})
	en.OnWord("").SetBlock(true).SetRule(rosm.OnlyAtMe()).Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Text(rosm.GetRandBotName(), "不在呢~"))
	})
}
