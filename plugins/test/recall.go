package test

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "撤回消息",
		Help: "- {回复消息}/撤回",
	})
	en.AddWord("/撤回").MUL("mys").Rule(mys.OnlyReply, mys.OnlyOverOwner).Handle(func(ctx *rosm.CTX) {
		if err := mys.Recall(ctx, ctx.Message.(*mys.MessageContent).Quote.OriginalMessageID, ctx.Message.(*mys.MessageContent).Quote.QuotedMessageSendTime, ctx.Being.RoomID); err != nil {
			log.Errorln("[recall]", err)
		} else {
			log.Infoln("[recall] 撤回成功,ID: ", ctx.Message.(*mys.MessageContent).Quote.OriginalMessageID)
		}
		//撤回触发者消息
		if err := mys.Recall(ctx, ctx.Being.MsgID[0], tool.Int64(ctx.Being.MsgID[1]), ctx.Being.RoomID); err != nil {
			log.Errorln("[recall]", err)
		} else {
			log.Infoln("[recall] 撤回成功,ID: ", ctx.Message.(*mys.MessageContent).Quote.OriginalMessageID)
		}
	})
}
