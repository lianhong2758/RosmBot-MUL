package test

import (
	"os"

	"github.com/FloatTech/imgfactory"
	"github.com/lianhong2758/RosmBot-MUL/message"
	draw "github.com/lianhong2758/RosmBot-MUL/plugins/public"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/marktoimage"
)

func init() {
	//插件注册
	en := rosm.Register(&rosm.PluginData{ //插件英文索引
		Name: "复读", //中文插件名
		Help: "- 复读..." +
			"- 用xx体渲染(抖动)文字", //插件帮助
	})
	en.OnRex(`^复读([\s\S]*)`).SetBlock(true).SetRule(func(ctx *rosm.Ctx) bool { return true }, rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) { //正则的触发方式
		ctx.Send(message.Text(ctx.Being.ResultWord[1])) //发送文字信息
	})
	en.OnRex("^复纯(.*)").Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Image("url://" + ctx.Being.ResultWord[1]))
	})
	en.OnRex(`^解析([\s\S]*)$`).Handle(func(ctx *rosm.Ctx) {
		m := ctx.Bot.BotSendCustom(ctx, ctx.Being.ResultWord[1])
		if m["code"] != "0" {
			ctx.Send(message.Text("发送失败: ", m["state"]))
		}
	})
	en.OnRex(`^用markdown渲染([\s\S]+)$`).Handle(func(ctx *rosm.Ctx) {
		data, err := os.ReadFile(draw.MaokenFontFile)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		rt := marktoimage.NewRichText(
			marktoimage.Config{
				FontData:   [][]byte{data},
				Width:      1080,
				Height:     1920,
				LineHeight: 8,
				TextMargin: 20,
				TopMargin:  10,
			})
		rt.ParseMarkdown(ctx.Being.ResultWord[1])
		rt.Draw()
		rt.Cut()
		img, _ := imgfactory.ToBytes(rt.Cov.Image())
		ctx.Send(message.ImageByte(img))
	})
	en.OnRex(`^(用.+)?渲染(抖动)?文字([\s\S]+)$`).Handle(func(ctx *rosm.Ctx) {
		font := ctx.Being.ResultWord[1]
		txt := ctx.Being.ResultWord[3]
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
		if ctx.Being.ResultWord[2] == "" {
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
