package test

import (
	"encoding/json"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

func init() {
	//插件注册
	en := rosm.Register(&rosm.PluginData{ //插件英文索引
		Name: "复读",      //中文插件名
		Help: "- 复读...", //插件帮助
	})
	en.AddRex("^复读(.*)").SetBlock(true).Rule(func(ctx *rosm.CTX) bool { return true }, rosm.OnlyMaster()).Handle(func(ctx *rosm.CTX) { //正则的触发方式
		ctx.Send(message.Text(ctx.Being.Rex[1])) //发送文字信息
	})
	en.AddRex("^复图(.*)文字(.*)").Handle(func(ctx *rosm.CTX) {
		con, _ := web.URLToConfig(ctx.Being.Rex[1])
		ctx.Send(message.ImageUrlWithText(web.UpImgUrl(ctx.Being.Rex[1]), con.Width, con.Height, 0, ctx.Being.Rex[2]))
	})
	en.AddRex("^复纯图(.*)").Handle(func(ctx *rosm.CTX) {
		con, _ := web.URLToConfig(ctx.Being.Rex[1])
		ctx.Send(message.ImageUrl(web.UpImgUrl(ctx.Being.Rex[1]), con.Width, con.Height, 0))
	})
	en.AddRex(`^解析([\s\S]*)$`).Handle(func(ctx *rosm.CTX) {
		info := new(map[string]any)
		err := json.Unmarshal(tool.StringToBytes(ctx.Being.Rex[1]), info)
		if err != nil {
			ctx.Send(message.Text("解析失败", err))
			return
		}
		r := ctx.Send(message.Custom(info))
		switch r.(type) {
		case *mys.SendState:
			if r.(*mys.SendState).ApiCode.Retcode != 0 {
				ctx.Send(message.Text("发送失败: ", r.(*mys.SendState).ApiCode.Message))
			}
		}
	})
}
