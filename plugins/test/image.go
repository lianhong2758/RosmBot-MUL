package test

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

const (
	referer = "https://weibo.com/"
	shouer  = "https://iw233.cn/api.php?sort=cat&referer"
	baisi   = "https://api.iw233.cn/seapi.php?sort=bs"
	heisi   = "https://api.iw233.cn/seapi.php?sort=hs"
	siwa    = "https://api.iw233.cn/seapi.php?sort=hbs"
	bizhi   = "https://iw233.cn/api.php?sort=iw233"
	baimao  = "https://iw233.cn/api.php?sort=yin"
	xing    = "https://iw233.cn/api.php?sort=xing"
	sese    = "https://api.iw233.cn/seapi.php?sort=setu"
	manghe  = "https://iw233.cn/api.php?sort=random"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "图片",
		Help: "- /随机壁纸" + " | " + "兽耳" + " | " + "星空" + " | " + "白毛" + " | " + "我要涩涩" + " | " + "涩涩达咩" + " | " + "白丝" + "\n" +
			"黑丝" + " | " + "丝袜" + " | " + "盲盒" + " | " + "开盲盒",
	})
	en.OnWord("随机壁纸", "兽耳", "星空", "白毛", "我要涩涩", "涩涩达咩", "白丝", "黑丝", "丝袜", "盲盒", "开盲盒").
		Handle(func(ctx *rosm.Ctx) {
			var url string
			switch ctx.Being.RawWord[1:] {
			case "兽耳":
				url = shouer
			case "随机壁纸":
				url = bizhi
			case "白毛":
				url = baimao
			case "星空":
				url = xing
			case "我要涩涩", "涩涩达咩":
				url = sese
			case "白丝":
				url = baisi
			case "黑丝":
				url = heisi
			case "丝袜":
				url = siwa
			case "盲盒", "开盲盒":
				url = manghe
			default:
				return
			}
			url2, err := web.GetRealURL(url)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			data, err := web.RequestDataWith(web.NewDefaultClient(), url2, "", referer, "", nil)
			if err != nil {
				ctx.Send(message.Text("获取图片失败惹"))
				return
			}
			ctx.Send(message.ImageByte(data))
		})
}
