package test

import (
	"github.com/lianhong2758/RosmBot-MUL/adapter/ob11"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "撤回消息",
		Help: "- 撤回",
	})
	en.OnRex(`^\[CQ:reply,id=(-?[0-9]+)\].*`).MUL("ob11").SetRule(rosm.OnlyMaster(), rosm.KeyWords("撤回")).Handle(func(ctx *rosm.Ctx) {
		ob11.DeleteMessage(ctx, ctx.Being.ResultWord[1])
		ob11.DeleteMessage(ctx, ctx.Being.MsgID)
		logrus.Info("[delete]撤回消息", ctx.Being.ResultWord[1], " - ", ctx.Being.MsgID[0])
	})
	//跟随撤回
	en.OnNoticeWithType(rosm.FriendRecall, rosm.GroupRecall).MUL("ob11").Handle(func(ctx *rosm.Ctx) {
		if id, ok := ctx.Message.(*zero.Event).MessageID.(int64); ok {
			if sids := rosm.GetMessageIDFormMapCache(tool.Int64ToString(id)); len(sids) > 0 {
				for _, sid := range sids {
					tool.WaitWhile()
					ob11.DeleteMessage(ctx, sid)
				}
			}
		}
	})
}
