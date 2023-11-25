package test

import (
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "帮助菜单",
		Help: "- /help\n",
	})
	en.AddWord("/help", "/帮助").Handle(func(ctx *rosm.CTX) {
		var msg strings.Builder
		msg.WriteString("*****菜单********")
		for _, v := range rosm.GetPlugins() {
			msg.WriteString("\n")
			msg.WriteString("#")
			msg.WriteString(v.Name)
			msg.WriteString("\n")
			msg.WriteString(v.Help)
			msg.WriteString("\n")
		}
		msg.WriteString("*****************")
		ctx.Send(message.Text(msg.String()))
	})
}
