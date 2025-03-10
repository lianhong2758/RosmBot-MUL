package test

import (
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/message"
	draw "github.com/lianhong2758/RosmBot-MUL/plugins/public"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "帮助菜单",
		Help: "- /help\n",
	})
	en.OnWord("/help", "/帮助").SetRule(rosm.OnlyAtMe()).Handle(func(ctx *rosm.Ctx) {
		var msg strings.Builder
		msg.WriteString("****RosmBot菜单*****")
		for _, v := range rosm.GetPlugins() {
			msg.WriteString("\n")
			msg.WriteString("#")
			msg.WriteString(v.Name)
			msg.WriteString("\n")
			msg.WriteString(v.Help)
			msg.WriteString("\n")
		}
		msg.WriteString("*******************")
		//ctx.Send(message.Text(msg.String()))
		image, err := draw.StringToPic(msg.String(), draw.MaokenFontFile)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.ImageByte(image))
	})
}
