package send

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/adapter/ob11"
	"github.com/lianhong2758/RosmBot-MUL/adapter/qq"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

// roomid12为两个room的结合,请使用tool.MergePadString的结果
func CTXBuild(types, botid, roomid12 string) (ctx *rosm.Ctx) {
	switch types {
	// case "mys":
	// 	room, villa := tool.SplitPadString(roomid12)
	// 	if botid == "" {
	// 		botid = mys.GetRandBot().BotToken.BotID
	// 	}
	// 	ctx = mys.NewCTX(botid, room, villa)
	case "qq_group", "qq_guild":
		id1, id2 := tool.SplitPadString(roomid12)
		if botid == "" {
			botid = qq.GetRandBot().BotID
		}
		ctx = qq.NewCTX(botid, types, id1, id2)
	case "ob11":
		group, _ := tool.SplitPadString(roomid12)
		if botid == "" {
			botid = ob11.GetRandBot().BotID
		}
		ctx = ob11.NewCTX(botid, group)
	}
	return ctx
}
