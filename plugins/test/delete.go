package test

import (
	"github.com/lianhong2758/RosmBot-MUL/adapter/ob11"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/sirupsen/logrus"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "撤回消息",
		Help: "- 撤回",
	})
	en.OnWord("撤回").MUL(ob11.BotType).SetRule(rosm.OnlyMaster(), rosm.OnlyReply()).Handle(func(ctx *rosm.Ctx) {
		ob11.DeleteMessage(ctx, ctx.State["reply"].(string))
		ob11.DeleteMessage(ctx, ctx.Being.MsgID)
		logrus.Info("[delete]撤回消息", ctx.State["reply"].(string), " - ", ctx.Being.MsgID)
	})
	//跟随撤回
	en.OnNoticeWithType(rosm.FriendRecall, rosm.GroupRecall).MUL(ob11.BotType).Handle(func(ctx *rosm.Ctx) {
		id := ctx.Being.MsgID
		if sids := rosm.GetMessageIDFormMapCache(id); len(sids) > 0 {
			for _, sid := range sids {
				tool.WaitWhile()
				ob11.DeleteMessage(ctx, sid)
				logrus.Info("[delete]跟随撤回消息", ctx.Being.MsgID, " - ", sid)
			}
		}
	})
}
