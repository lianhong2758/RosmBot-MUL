package test

import (
	"encoding/json"

	"github.com/lianhong2758/RosmBot-MUL/message"
	draw "github.com/lianhong2758/RosmBot-MUL/plugins/public"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

func init() {
	//插件注册
	en := rosm.Register(&rosm.PluginData{ //插件英文索引
		Name: "复读",      //中文插件名
		Help: "- 复读...", //插件帮助
	})
	en.AddRex("^复读(.*)").SetBlock(true).Rule(func(ctx *rosm.Ctx) bool { return true }, rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) { //正则的触发方式
		ctx.Send(message.Text(ctx.Being.Rex[1])) //发送文字信息
	})
	en.AddRex("^复纯(.*)").Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Image("url://" + ctx.Being.Rex[1]))
	})
	en.AddRex(`^解析([\s\S]*)$`).Handle(func(ctx *rosm.Ctx) {
		info := new(map[string]any)
		err := json.Unmarshal(tool.StringToBytes(ctx.Being.Rex[1]), info)
		if err != nil {
			ctx.Send(message.Text("解析失败", err))
			return
		}
		m := ctx.Send(message.Custom(info))
		if m["code"] != 0 {
			ctx.Send(message.Text("发送失败: ", m["state"]))
		}
	})
	en.AddRex(`^(用.+)?渲染(抖动)?文字([\s\S]+)$`).Handle(func(ctx *rosm.Ctx) {
		font := ctx.Being.Rex[1]
		txt := ctx.Being.Rex[3]
		switch font {
		case "用终末体":
			font = draw.SyumatuFontFile
		case "用终末变体":
			font = draw.NisiFontFile
		case "用紫罗兰体":
			font = draw.VioletEvergardenFontFile
		case "用樱酥体":
			font = draw.SakuraFontFile
		case "用Consolas体":
			font = draw.ConsolasFontFile
		case "用苹方体":
			font = draw.FontFile
		case "猫啃体":
			font = draw.MaokenFontFile
		default:
			font = draw.MaokenFontFile
		}
		if ctx.Being.Rex[2] == "" {
			image, err := draw.StringToPic(txt, font)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			ctx.Send(message.ImageByte(image))
		} else {
			image, err := draw.StringToShake(txt, font)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			ctx.Send(message.ImageByte(image))
		}
	})
}
