package test

import (
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "房间列表",
		Help: "- /房间列表",
	})
	en.AddWord("/房间列表").Handle(func(ctx *rosm.CTX) {
		result, err := mys.GetRoomList(ctx)
		if err != nil {
			ctx.Send(message.Text("获取信息失败", err))
		}
		var msg strings.Builder
		msg.WriteString("别野")
		msg.WriteString(ctx.Being.RoomID2)
		msg.WriteString(":\n")
		for _, v := range result.Data.List {
			if v.GroupID == "0" {
				continue
			}
			if msg.String() != "" {
				msg.WriteByte('\n')
			}
			msg.WriteString("#" + v.GroupName)
			msg.WriteString("(" + v.GroupID + "):")
			for _, vv := range v.RoomList {
				msg.WriteString("\n" + vv.RoomName)
				msg.WriteString("(" + vv.RoomID + ")")
			}
		}
		ctx.Send(message.Text(msg.String()))
	})
}
