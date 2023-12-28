package test

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
	log "github.com/sirupsen/logrus"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "撤回消息",
		Help: "- {回复消息}/撤回",
	})
	en.AddWord("/撤回").MUL("mys").Rule(rosm.OnlyReply(), rosm.OnlyOverHost()).Handle(func(ctx *rosm.Ctx) {
		if err := mys.Recall(ctx, ctx.Being.Def["Content"].(*mys.MessageContent).Quote.OriginalMessageID, ctx.Being.Def["Content"].(*mys.MessageContent).Quote.QuotedMessageSendTime, ctx.Being.RoomID); err != nil {
			log.Errorln("[recall]", err)
		} else {
			log.Infoln("[recall] 撤回成功,ID: ", ctx.Being.Def["Content"].(*mys.MessageContent).Quote.OriginalMessageID)
		}
		//撤回触发者消息
		if err := mys.Recall(ctx, ctx.Being.MsgID[0], ctx.Being.MsgID[1], ctx.Being.RoomID); err != nil {
			log.Errorln("[recall]", err)
		} else {
			log.Infoln("[recall] 撤回成功,ID: ", ctx.Being.Def["Content"].(*mys.MessageContent).Quote.OriginalMessageID)
		}
	})
}
