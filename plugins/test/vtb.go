package test

import (
	"encoding/json"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

const url = "http://47.93.28.113/vtbwife?id="

func init() { // 插件主体
	en := rosm.Register(&rosm.PluginData{
		Name: "抽vtb老婆",
		Help: "- /抽vtb(老婆)",
	})
	en.AddRex(`^/抽(vtb|VTB)(老婆)?$`).Handle(func(ctx *rosm.CTX) {
		body, err := web.GetData(url+ctx.Being.User.ID, "")
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		var r result
		err = json.Unmarshal(body, &r)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		con, _ := web.URLToConfig(r.Imgurl)
		ctx.Send(message.AT(ctx.Being.User.ID, ctx.Being.User.Name), message.ImageUrlWithText(web.UpImgUrl(r.Imgurl), con.Width, con.Height, 0, "\n今天你的VTB老婆是: "+r.Name))
		ctx.Send(message.Text(r.Message))
	})
}

type result struct {
	Code    int    `json:"code"`
	Imgurl  string `json:"imgurl"`
	Name    string `json:"name"`
	Message string `json:"message"`
}
