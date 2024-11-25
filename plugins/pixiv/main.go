package pixiv

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/jozsefsallai/gophersauce"
	"github.com/lianhong2758/RosmBot-MUL/adapter/ob11"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	"github.com/sirupsen/logrus"
)

var (
	saucenaocli *gophersauce.Client
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "pixiv",
		Help: "- /搜图 pid\n" +
			"- /以图搜图\n" +
			"- 设置 saucenao api key [apikey]",
		DataFolder: "pixiv",
	})
	if file.IsNotExist(en.DataFolder + "cache") {
		_ = os.MkdirAll(en.DataFolder+"cache", 0755)
	}

	en.OnRex(`^/搜图\s*(\d+)$`).MUL(ob11.BotType).Handle(func(ctx *rosm.Ctx) {
		id, _ := strconv.ParseInt(ctx.Being.ResultWord[1], 10, 64)
		ctx.Send(message.Text("雪儿正在寻找中......"))
		// 获取P站插图信息
		illust, err := Works(id)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		if illust.R18 {
			ctx.Send(message.Text("含有R18图片,请自行下载"))
		}
		remessage := "直链: "
		if illust.Pid > 0 {
			name := strconv.FormatInt(illust.Pid, 10)
			var imgs message.Message
			for i, v := range illust.URLs {
				f := fmt.Sprint(file.BOTPATH, "/", en.DataFolder, "cache/", name, "-", i+1, ".png")
				n := fmt.Sprint(name, "-", i+1, ".png")
				remessage += fmt.Sprint("\nhttps://pixiv.re/", n)
				//下载非R18图
				if file.IsNotExist(f) {
					logrus.Debugln("[pixiv]开始下载", n)
					logrus.Debugln("[pixiv]urls:", v)
					err1 := DownLoadWorks(v, f)
					if err1 != nil {
						logrus.Debugln("[pixiv]下载err:", err1)
						continue
					}
				}
				///发送非R18图
				if !illust.R18 {
					imgs = append(imgs, message.Image("file://"+f))
				}
			}
			if len(illust.URLs) == 1 {
				remessage = fmt.Sprint("直链: ", "https://pixiv.re/", illust.Pid, ".png")
			}
			txt := message.Text(
				"标题: ", illust.Title, "\n",
				"插画ID: ", illust.Pid, "\n",
				"画师: ", illust.UserName, "\n",
				"画师ID: ", illust.UID, "\n",
				remessage,
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
	en.OnWord("以图搜图", "以图识图").MUL(ob11.BotType).Handle(func(ctx *rosm.Ctx) {
		pics := GetMustPic(ctx)
		if len(pics) == 0 {
			ctx.Send(message.Text("雪儿没有收到图片唔..."))
			return
		}
		ctx.Send(message.Text("雪儿正在寻找中..."))
		for _, pic := range pics {
			//saucenao
			//	find := false
			if saucenaocli != nil || func() bool {
				saucenaokey, _ := rosm.PluginDB.FindString(en.Name, "0")
				if saucenaokey != "" {
					saucenaocli, _ = gophersauce.NewClient(&gophersauce.Settings{
						MaxResults: 1,
						APIKey:     saucenaokey,
					})
					return true
				}
				return false
			}() {
				resp, err := saucenaocli.FromURL(pic)
				if err == nil && resp.Count() > 0 {
					result := resp.First()
					s, err := strconv.ParseFloat(result.Header.Similarity, 64)
					if err == nil {
						//	find = s > 80.0
						images, err := web.GetData(result.Header.Thumbnail, web.RandUA())
						if err != nil {
							logrus.Info("[pixiv]下载预览图失败,ERROR: ", err)
						}
						ctx.Send(message.Message{ob11.FakeSenderForwardNode(ctx,
							message.Text("saucenao搜图结果: ", "\n匹配度: ", s, "%", "\n图源: ",
								result.Header.IndexName, "\n",
								result.Data.Source, "\n",
								strings.Join(result.Data.ExternalURLs, "\n"))),
							ob11.FakeSenderForwardNode(ctx, message.ImageByte(images))}...)
					}
				}

			} else {
				ctx.Send(message.Text("请私聊发送 设置 saucenao api key [apikey] 以启用 saucenao 搜图 (方括号不需要输入), key 请前往 https://saucenao.com/user.php?page=search-api 获取"))
			}
			//未找到时再调用ascii2d 搜索
			// if !find{}
		}
	})
	en.OnRex(`^\/设置\s?saucenao\s?api\s?key\s?([0-9a-f]{40})$`).SetRule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		var err error
		saucenaocli, err = gophersauce.NewClient(&gophersauce.Settings{
			MaxResults: 1,
			APIKey:     ctx.Being.ResultWord[1],
		})
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		rosm.PluginDB.InsertString(en.Name, "0", ctx.Being.ResultWord[1])
		ctx.Send(message.Text("设置成功!"))
	})
}

// 想办法获取一张图片,仅限ob11
func GetMustPic(ctx *rosm.Ctx) []string {
	var urls = GetPicFormCtx(ctx)
	if len(urls) > 0 {
		return urls
	}
	//额外请求一次图片
	ctx.Send(message.Text("请给雪儿一张图片..."))
	next, close := ctx.GetNext(rosm.AllMessage, false, rosm.OnlyTheUser(ctx.Being.User.ID))
	defer close()
	for {
		select {
		case <-time.After(time.Second * 120):
			return urls
		case newCtx := <-next:
			if us := GetPicFormCtx(newCtx); len(us) > 0 {
				return us
			}
		}
	}

}

// 获取这次ctx内容的图片
func GetPicFormCtx(ctx *rosm.Ctx) []string {
	var urls = []string{}
	e, _ := ctx.State["event"].(*ob11.Event)
	for _, elem := range e.Message {
		if elem.Type == "image" {
			if elem.Data["url"] != "" {
				urls = append(urls, elem.Data["url"])
			}
		}
	}
	return urls
}
