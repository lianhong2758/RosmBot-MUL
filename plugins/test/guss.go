package test

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func init() {
	//插件注册
	en := rosm.Register(&rosm.PluginData{
		Name: "猜数字",
		Help: "- 开始猜数字",
	})
	en.AddWord("开始猜数字").Handle(func(ctx *rosm.CTX) {
		num := strconv.Itoa(rand.Intn(9) + 1)
		next, stop := ctx.GetNextUserMess()
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
}
