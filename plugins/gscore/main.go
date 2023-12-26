package gscore

import (
	"context"
	"time"

	"github.com/FloatTech/ttl"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

var cache = ttl.NewCache[string, *rosm.Ctx](time.Minute * 3)

func init() {
	en := rosm.Register(rosm.NewRegist("gscore", "参考score帮助\nGenshinUID帮助页\nhttps://www.kdocs.cn/l/ccpc6z0bZx6u", "gscore"))
	en.AddWord("启动gscore").Handle(func(ctx *rosm.Ctx) { //.Rule(rosm.OnlyMaster())
		//创建ws
		Config.NewWebSocket()
		//启动ws接收
		var ctxback context.Context
		ctxback, Config.cancel = context.WithCancel(context.Background())
		go ReadAndSendMessage(ctxback, Config.conn)
		Config.on = true
		ctx.Send(message.Text("gscore已启动"))
	})
	en.AddWord("关闭gscore").Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		Config.on = false
		//启动ws接收
		Config.cancel()
		ctx.Send(message.Text("gscore已关闭"))
	})
	en.AddOther(rosm.AllMessage).Rule(func(ctx *rosm.Ctx) bool { return Config.on }).Handle(func(ctx *rosm.Ctx) {
		for {
			SendErr := SendWsMessage(MakeSendCoreMessage(ctx), Config.conn)
			if SendErr != nil {
				time.Sleep(time.Second * 5)
				log.Error("[gscore]SendErr", SendErr)
				Config.RecoveWebScoket()
				continue
			}
			break
		}
	})
	//最后init
	configInit()
}
