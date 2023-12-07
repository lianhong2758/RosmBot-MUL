// 用于管理插件的启用/禁用
package rosm

import (
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

var en = Register(&PluginData{
	Name: "插件管理",
	Help: "- /用法 xxx\n" +
		"- /禁用 xxx\n" +
		"- /启用 xxx",
	DataFolder: "regulate",
})

func init() {
	en.AddRex(`/启用\s*(.*)`).Handle(func(ctx *CTX) {
		name := ctx.Being.Rex[1]
		if _, ok := GetPlugins()[name]; !ok {
			ctx.Send(message.Text("未找到插件: ", name))
			return
		}
		err := PluginDB.InsertOff(name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2), false)
		if err != nil {
			ctx.Send(message.Text(name, "启用失败,ERROR: ", err))
			return
		}
		ctx.Send(message.Text(name, "已启用..."))
	})
	en.AddRex(`/禁用\s*(.*)`).Handle(func(ctx *CTX) {
		name := ctx.Being.Rex[1]
		if _, ok := GetPlugins()[name]; !ok {
			ctx.Send(message.Text("未找到插件: ", name))
			return
		}
		err := PluginDB.InsertOff(name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2), true)
		if err != nil {
			ctx.Send(message.Text(name, "禁用失败,ERROR: ", err))
			return
		}
		ctx.Send(message.Text(name, "已禁用..."))
	})
	en.AddRex(`/用法\s*(.*)`).Handle(func(ctx *CTX) {
		name := ctx.Being.Rex[1]
		plugin, ok := GetPlugins()[name]
		if !ok {
			ctx.Send(message.Text("未找到插件: ", name))
			return
		}
		var msg strings.Builder
		msg.WriteString("******")
		msg.WriteString(plugin.Name)
		msg.WriteString("******\n")
		msg.WriteString(plugin.Help)
		msg.WriteByte('\n')
		msg.WriteString("********************")
		ctx.Send(message.Text(msg.String()))
	})
}
