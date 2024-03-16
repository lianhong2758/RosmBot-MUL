package test

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	"github.com/tidwall/gjson"
)

const api = "http://dwz.2xb.cn/zaob"

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "今日早报",
		Help: "- /今日早报",
	})
	en.AddWord("/今日早报").Handle(func(ctx *rosm.Ctx) {
		data, err := web.GetData(api, "")
		if err != nil {
			return
		}
		picURL := gjson.Get(tool.BytesToString(data), "imageUrl").String()
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text("今日早报送达~"), message.Image("url://"+picURL))
	})
}
