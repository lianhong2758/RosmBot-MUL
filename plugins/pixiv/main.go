package pixiv

import (
	"fmt"
	"strconv"

	"github.com/FloatTech/floatbox/file"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/ob11"
	"github.com/sirupsen/logrus"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name:       "pixiv",
		Help:       "- /搜图 pid",
		DataFolder: "pixiv",
	})

	en.AddRex(`^/搜图\s*(d+)$`).Handle(func(ctx *rosm.Ctx) {
		id, _ := strconv.ParseInt(ctx.Being.Rex[1], 10, 64)
		ctx.Send(message.Text("雪儿正在寻找中......"))
		// 获取P站插图信息
		illust, err := Works(id)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		r18 := HaveR18Pic(illust)
		if r18 {
			ctx.Send(message.Text("含有R18图片,请自行下载"))
		}
		if illust.ID > 0 {
			name := strconv.FormatInt(illust.ID, 10)
			var imgs message.Message
			for i, v := range illust.MetaPages {
				f := fmt.Sprint(file.BOTPATH, "/", en.DataFolder, "cache/", name, "_p", i, ".png")
				n := fmt.Sprint(name, "_p", i, ".png")
				if file.IsNotExist(f) {
					logrus.Debugln("[pixiv]开始下载", n)
					logrus.Debugln("[pixiv]urls:", v.ImageUrls.Original)
					err1 := DownLoadWorks(v.ImageUrls.Original, f)
					if err1 != nil {
						logrus.Debugln("[pixiv]下载err:", err1)
						continue
					}
				}
				imgs = append(imgs, message.Image("file://"+f))
			}
			txt := message.Text(
				"标题: ", illust.Title, "\n",
				"插画ID: ", illust.ID, "\n",
				"画师: ", illust.User.Name, "\n",
				"画师ID: ", illust.User.ID, "\n",
				"直链: ", "https://pixiv.re/", illust.ID, ".jpg",
			)
			if imgs != nil {
				ctx.Send(message.Message{ob11.FakeSenderForwardNode(ctx, txt),
					ob11.FakeSenderForwardNode(ctx, imgs...)}...)
			} else {
				// 图片下载失败，仅发送文字结果
				ctx.Send(txt)
			}
		} else {
			ctx.Send(message.Text("图片不存在呜..."))
		}
	})
}
