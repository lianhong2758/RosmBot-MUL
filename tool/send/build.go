package send

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
	"github.com/lianhong2758/RosmBot-MUL/server/qq"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

// roomid为两个room的结合,请使用tool.String221的结果
func CTXBuild(types, botid, roomid21 string) (ctx *rosm.CTX) {
	switch types {
	case "mys":
		room, villa := tool.SplitPadString(roomid21)
		if botid == "" {
			botid = mys.GetBot().BotToken.BotID
		}
		ctx = mys.NewCTX(botid, room, villa)
	case "qq":
		id1, id2 := tool.SplitPadString(roomid21)
		if botid == "" {
			botid = qq.GetBot().BotID
		}
		ctx = qq.NewCTX(botid, id1, id2)
	}
	return ctx
}
