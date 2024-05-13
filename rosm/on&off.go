// 用于管理插件的启用/禁用
package rosm

import (
	"strconv"
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
	en.AddRex(`/启用\s*(.*)`).Handle(func(ctx *Ctx) {
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
	en.AddRex(`/禁用\s*(.*)`).Handle(func(ctx *Ctx) {
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
	en.AddRex(`/用法\s*(.*)`).Handle(func(ctx *Ctx) {
		name := ctx.Being.Rex[1]
		plugin, ok := GetPlugins()[name]
		if !ok {
			ctx.Send(message.Text("未找到插件: ", name))
			return
		}
		var msg strings.Builder
		t := strings.Repeat("*", (20-len(plugin.Name))/2)
		msg.WriteString(t)
		msg.WriteString(plugin.Name)
		msg.WriteString(t)
		msg.WriteByte('\n')
		msg.WriteString("启用状态:")
		msg.WriteString(strconv.FormatBool(PluginIsOn(plugin)(ctx)))
		msg.WriteByte('\n')
		msg.WriteString("帮助信息:")
		msg.WriteString(plugin.Help)
		msg.WriteByte('\n')
		msg.WriteString(strings.Repeat("*", 20))
		ctx.Send(message.Text(msg.String()))
	})
	//去除全局响应沉默的影响
	for _, m := range en.Matchers {
		m.rules = m.rules[1:len(m.rules)-1]
	}
}
