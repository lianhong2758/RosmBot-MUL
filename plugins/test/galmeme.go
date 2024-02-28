package test

import (
	"fmt"
	rand "math/rand/v2"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	listURL  = "https://sticker.kungal.com/sticker/%d/__data.json"
	imageURL = "https://sticker.kungal.com/kun-galgame-stickers/telegram/KUNgal%d/%d.png"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "gal表情包",
		Help: "- /表情包\n",
	})
	en.AddWord("/gal表情包", "/表情包").Handle(func(ctx *rosm.Ctx) {
		//随机一个集合id (1-6)
		listId := rand.IntN(6) + 1
		data, err := web.GetData(fmt.Sprintf(listURL, listId), "")
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		node := gjson.ParseBytes(data).Get("nodes.1.data")
		id := rand.IntN(len(node.Get("2").Array()))
		roleindex := node.Get("2").Array()[id].String()
		roledata := node.Get(roleindex)
		name := node.Get(roledata.Get("loli").String()).String()
		gameName := node.Get(roledata.Get("game").String()).String()
		url := fmt.Sprintf(imageURL, listId, id)
		logrus.Print(name, gameName, url)
		ctx.Send(message.Text("表情包已送达:\n角色名: "+name+"\n来源: ", gameName), message.Image("url://"+url))
	})
}
