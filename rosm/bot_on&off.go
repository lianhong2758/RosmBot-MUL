// 用于管理插件的启用/禁用
package rosm

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

// 仅作为存储开关数据使用,没有Matcher
var boten = Register(&PluginData{
	Name: "响应管理",
	Help: "- @bot/早安\n" +
		"- @bot/晚安",
	//借用插件管理的存储器
})

func init() {
	en.AddWord(`/早安`).Rule(OnlyAtMe()).Handle(func(ctx *Ctx) {
	 on:=PluginIsOn(boten)(ctx)
		if !on {
			err := PluginDB.InsertOff(boten.Name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2), false)
			if err != nil {
				ctx.Send(message.Text("响应失败,ERROR: ", err))
				return
			}
			ctx.Send(message.Text("早安,", ctx.Bot.Card().BotName, "开始工作了喵~"))
		} else {
			// 已经响应了
			ctx.Send(message.Text("早安,", ctx.Bot.Card().BotName, "已经在认真工作了喵~"))
		}
	})
	en.AddWord(`/晚安`).Rule(OnlyAtMe()).Handle(func(ctx *Ctx) {
		on:=PluginIsOn(boten)(ctx)
		if on {
			err := PluginDB.InsertOff(boten.Name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2), true)
			if err != nil {
				ctx.Send(message.Text("沉默失败,ERROR: ", err))
				return
			}
			ctx.Send(message.Text("晚安,", ctx.Bot.Card().BotName, "要睡觉了喵~"))
		} else {
			// 已经响应了
			ctx.Send(message.Text("晚安,Zzz~"))
		}
	})

}
