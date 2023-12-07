package test

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
)

func init() {
	//插件注册
	en := rosm.Register(&rosm.PluginData{
		Name: "猜数字/猜拳",
		Help: "- /开始猜数字/猜拳",
	})
	en.AddWord("/开始猜数字").Handle(func(ctx *rosm.CTX) {
		num := strconv.Itoa(rand.Intn(9) + 1)
		next, stop := ctx.GetNext(rosm.AllMessage, true, rosm.OnlyTheUser(ctx.Being.User.ID))
		defer stop()
		ctx.Send(message.Reply(), message.Text("开始猜数字,大小为1-10,你有3次机会"))
		for i := 0; i < 3; i++ {
			select {
			case <-time.After(time.Second * 180):
				ctx.Send(message.Text("时间太久了"))
				return
			case ctx2 := <-next:
				switch {
				case ctx2.Being.Word == num:
					ctx2.Send(message.Reply(), message.Text("恭喜你猜对了"))
					return
				case ctx2.Being.Word > num:
					ctx2.Send(message.Reply(), message.Text("你猜大了"))
				case ctx2.Being.Word < num:
					ctx2.Send(message.Reply(), message.Text("你猜小了"))
				}
			}
		}
		ctx.Send(message.Text("游戏失败"))
	})
	en.AddWord("/开始猜拳").MUL("mys").Handle(func(ctx *rosm.CTX) {
		p := mys.NewPanel()
		p.Small(false, &mys.Component{
			ID:    "t1",
			Text:  "剪刀",
			Type:  1,
			CType: 1,
			Extra: "剪刀",
		})
		p.Small(false, &mys.Component{
			ID:    "t2",
			Text:  "石头",
			Type:  1,
			CType: 1,
			Extra: "石头",
		})
		p.Small(false, &mys.Component{
			ID:    "t3",
			Text:  "布",
			Type:  1,
			CType: 1,
			Extra: "布",
		})
		p.Title("请双方选择出拳:")
		ctx.Send(message.Custom(p))
		var ctxs []*rosm.CTX
		next, stop := ctx.GetNext(rosm.Click, false, func(ctx2 *rosm.CTX) bool {
			return (ctx2.Being.Word == "石头" || ctx2.Being.Word == "剪刀" || ctx2.Being.Word == "布") && ctx2.Being.RoomID == ctx.Being.RoomID
		})
		defer stop()
		for {
			select {
			case <-time.After(time.Second * 180):
				ctx.Send(message.Text("时间太久了"))
				return
			case ctxt := <-next:
				switch len(ctxs) {
				case 0:
					ctxs = append(ctxs, ctxt)
				case 1:
					if ctxt.Being.User.ID != ctxs[0].Being.User.ID {
						ctxs = append(ctxs, ctxt)
						victory(ctxs)
						return
					}
				}

			}
		}
	})
}
func victory(ctxs []*rosm.CTX) {
	var msg = []message.MessageSegment{message.Text("出拳结果: \n"),
		message.Text(mys.GetUserName(ctxs[0], ctxs[0].Being.User.ID), ":", ctxs[0].Being.Word, "\n"),
		message.Text(mys.GetUserName(ctxs[1], ctxs[1].Being.User.ID), ":", ctxs[1].Being.Word, "\n"),
	}
	switch {
	case ctxs[0].Being.Word == ctxs[1].Being.Word:
		ctxs[0].Send(append(msg, message.Text("平局"))...)
	case ctxs[0].Being.Word == "石头" && ctxs[1].Being.Word == "剪刀" || ctxs[0].Being.Word == "剪刀" && ctxs[1].Being.Word == "布" || ctxs[0].Being.Word == "布" && ctxs[1].Being.Word == "石头":
		ctxs[0].Send(append(msg, message.Text(mys.GetUserName(ctxs[0], ctxs[0].Being.User.ID), "获胜"))...)
	default:
		ctxs[0].Send(append(msg, message.Text(mys.GetUserName(ctxs[1], ctxs[1].Being.User.ID), "获胜"))...)

	}
}
