package onplugin

// 用于管理插件的启用/禁用

import (
	"strconv"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

var en = rosm.Register(&rosm.PluginData{
	Name: "插件管理",
	Help: "- /用法 xxx\n" +
		"- /禁用 xxx\n" +
		"- /启用 xxx",
})

func init() {
	en.OnRex(`^/启用\s*(.*)`).Handle(func(ctx *rosm.Ctx) {
		name := ctx.Being.ResultWord[1]
		if _, ok := rosm.GetPlugins()[name]; !ok {
			ctx.Send(message.Text("未找到插件: ", name))
			return
		}
		err := rosm.PluginDB.InsertOff(name, tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID), false)
		if err != nil {
			ctx.Send(message.Text(name, "启用失败,ERROR: ", err))
			return
		}
		ctx.Send(message.Text(name, "已启用..."))
	})
	en.OnRex(`^/禁用\s*(.*)`).Handle(func(ctx *rosm.Ctx) {
		name := ctx.Being.ResultWord[1]
		if _, ok := rosm.GetPlugins()[name]; !ok {
			ctx.Send(message.Text("未找到插件: ", name))
			return
		}
		err := rosm.PluginDB.InsertOff(name, tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID), true)
		if err != nil {
			ctx.Send(message.Text(name, "禁用失败,ERROR: ", err))
			return
		}
		ctx.Send(message.Text(name, "已禁用..."))
	})
	en.OnRex(`^/用法\s*(.*)`).Handle(func(ctx *rosm.Ctx) {
		name := ctx.Being.ResultWord[1]
		plugin, ok := rosm.GetPlugins()[name]
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
		msg.WriteString(strconv.FormatBool(rosm.PluginIsOn(plugin)(ctx)))
		msg.WriteByte('\n')
		msg.WriteString("帮助信息:")
		msg.WriteString(plugin.Help)
		msg.WriteByte('\n')
		msg.WriteString(strings.Repeat("*", 20))
		ctx.Send(message.Text(msg.String()))
	})
	en.OnWord(`早安`).SetRule(rosm.OnlyAtMe()).Handle(func(ctx *rosm.Ctx) {
		on := rosm.PluginIsOn(rosm.GetBoten())(ctx)
		if !on {
			err := rosm.PluginDB.InsertOff(rosm.GetBoten().Name, tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID), false)
			if err != nil {
				ctx.Send(message.Text("响应失败,ERROR: ", err))
				return
			}
			ctx.Send(message.Text("早安,", rosm.GetRandBotName(), "开始工作了喵~"))
		} else {
			// 已经响应了
			ctx.Send(message.Text("早安,", rosm.GetRandBotName(), "已经在认真工作了喵~"))
		}
	})
	en.OnWord(`晚安`).SetRule(rosm.OnlyAtMe()).Handle(func(ctx *rosm.Ctx) {
		on := rosm.PluginIsOn(rosm.GetBoten())(ctx)
		if on {
			err := rosm.PluginDB.InsertOff(rosm.GetBoten().Name, tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID), true)
			if err != nil {
				ctx.Send(message.Text("沉默失败,ERROR: ", err))
				return
			}
			ctx.Send(message.Text("晚安,", rosm.GetRandBotName(), "要睡觉了喵~"))
		} else {
			// 已经响应了
			ctx.Send(message.Text("晚安,Zzz~"))
		}
	})
	//去除全局响应沉默的影响
	rosm.DeleteOffRule(en)
}
