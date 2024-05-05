package test

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/ob11"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "撤回消息",
		Help: "- 撤回\n",
	})
	en.AddRex(`^\[CQ:reply,id=(-?[0-9]+)\].*`).MUL("ob11").Rule(rosm.OnlyMaster(), rosm.KeyWords("撤回")).Handle(func(ctx *rosm.Ctx) {
		ob11.DeleteMessage(ctx, ctx.Being.Rex[1])
		ob11.DeleteMessage(ctx, ctx.Being.MsgID[0])
	})
}
