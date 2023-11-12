package test

import (
	"os"
	"time"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
	"github.com/lianhong2758/RosmBot-MUL/server/mys/mysmsg"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "测试",
		Help: "- 测试\n" +
			"- 测试下标\n" +
			"- 测试预览\n" +
			"- 测试全体next\n" +
			"- 测试个人next\n" +
			"- 测试表情",
	})
	en.AddWord("测试").Handle(func(ctx *rosm.CTX) {
		ctx.Send(message.Text("你好"), message.AT(ctx.Being.User.ID, ctx.Being.User.Name), message.Link("www.baidu.com", false, "百度一下"), mysmsg.RoomLink(ctx.Being.RoomID2, ctx.Being.RoomID, "# 本房间"), message.Text("[爱心]"))
	})
	en.AddWord("测试下标跳转房间").Handle(func(ctx *rosm.CTX) {
		s := mysmsg.BadgeStr{
			Icon: "http://47.93.28.113/favicon.ico",
			Text: "10248",
			URL:  "https://dby.miyoushe.com/chat/463/10248",
		}
		ctx.Send(message.Text("大别野房间~"), mysmsg.Badge(s))
	})
	en.AddWord("测试下标").Handle(func(ctx *rosm.CTX) {
		s := mysmsg.BadgeStr{
			Icon: "http://47.93.28.113/favicon.ico",
			Text: "清雪官方",
			URL:  "http://47.93.28.113",
		}
		ctx.Send(message.Text("清雪官网~"), mysmsg.Badge(s))
	})
	en.AddWord("测试预览").Handle(func(ctx *rosm.CTX) {
		s := mysmsg.PreviewStr{
			Icon:       "http://47.93.28.113/favicon.ico",
			URL:        "http://47.93.28.113",
			ImageURL:   "http://47.93.28.113/ippic",
			IsIntLink:  true,
			SourceName: "我是喵喵喵~",
			Title:      "这是一个标题测试",
			Content:    "我是具体内容",
		}
		ctx.Send(message.Text("测试"), mysmsg.Preview(s))
	})
	en.AddWord("测试全体next").Handle(func(ctx *rosm.CTX) {
		next, stop := ctx.GetNext(rosm.AllMessage, true)
		defer stop()
		ctx.Send(message.Text("测试开始"))
		for i := 0; i < 3; i++ {
			select {
			case <-time.After(time.Second * 60):
				ctx.Send(message.Text("时间太久了"))
				return
			case ctx2 := <-next:
				ctx.Send(message.Text("这是全体下一句话:", ctx2.Being.Word))
			}
		}
	})
	en.AddWord("测试个人next").Handle(func(ctx *rosm.CTX) {
		next, stop := ctx.GetNext(rosm.AllMessage, true, rosm.OnlyTheUser(ctx.Being.User.ID))
		defer stop()
		ctx.Send(message.Text("测试开始"))
		for i := 0; i < 3; i++ {
			select {
			case <-time.After(time.Second * 60):
				ctx.Send(message.Text("时间太久了"))
				return
			case ctx2 := <-next:
				ctx.Send(message.Text("这是个人下一句话:", ctx2.Being.Word))
			}
		}
	})
	en.AddWord("测试表态").MUL("mys").Handle(func(ctx *rosm.CTX) {
		result := ctx.Send(message.Text("测试开始,表态此条消息"))
		next, stop := ctx.GetNext(rosm.Quick, true, func(ctx *rosm.CTX) bool {
			if result.(*mys.SendState).Data.BotMsgID == ctx.Message.(*mys.InfoSTR).Event.ExtendData.EventData.AddQuickEmoticon.BotMsgID {
				return true
			}
			return false
		})
		defer stop()
		for i := 0; i < 3; i++ {
			select {
			case <-time.After(time.Second * 60):
				ctx.Send(message.Text("时间太久了"))
				return
			case ctx2 := <-next:
				ctx.Send(message.Text("这是表态结果:\n", ctx2.Message.(*mys.InfoSTR).Event.ExtendData.EventData.AddQuickEmoticon))
			}
		}
	})
	en.AddWord("测试视频").MUL("mys").Handle(func(ctx *rosm.CTX) {
		ctx.Send(message.Text("测试开始"))
		s := mysmsg.PreviewStr{
			Icon:       "http://47.93.28.113/favicon.ico",
			URL:        "http://47.93.28.113/file?path=CSGO/1.mp4",
			IsIntLink:  true,
			SourceName: "清雪API",
			Title:      "测试视频",
			Content:    "CSGO精彩击杀,完美竞技平台",
		}
		ctx.Send(message.Text("视频测试"), mysmsg.Preview(s))
	})
	en.AddWord("测试组合").MUL("mys").Handle(func(ctx *rosm.CTX) {
		ctx.Send(message.Text("测试开始"))
		s := mysmsg.PreviewStr{
			Icon:       "http://47.93.28.113/favicon.ico",
			URL:        "http://47.93.28.113/file?path=CSGO/1.mp4",
			IsIntLink:  true,
			SourceName: "清雪API",
			Title:      "测试视频",
			Content:    "CSGO精彩击杀,完美竞技平台",
		}
		ss := mysmsg.BadgeStr{
			Icon: "http://47.93.28.113/favicon.ico",
			Text: "清雪官方",
			URL:  "http://47.93.28.113",
		}
		ctx.Send(message.Text("视频测试"), mysmsg.Preview(s), mysmsg.Badge(ss))
	})
	en.AddWord("测试上传图片").MUL("mys").Handle(func(ctx *rosm.CTX) {
		file, _ := os.ReadFile("data/public/测试.jpg")
		url, err := mys.UploadFile(ctx, file)
		if err != nil {
			ctx.Send(message.Text("发送失败,ERROR:", err))
			return
		}
		ctx.Send(message.ImageUrl(url, 0, 0, 0))
	})
}

//ctx有消息的全部信息,ctx.Being有简单的消息信息获取
//ctx.Send(...)是发送消息的基本格式
