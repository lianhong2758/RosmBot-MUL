package mys

import (
	"time"

	vila_bot "github.com/lianhong2758/RosmBot-MUL/server/mys/proto"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

func (c *Config) Listen() {
	log.Infoln("[mys-ws]bot开始监听消息")
	pack := NewDataPack(c.conn)
	var msg *Message
	lastheartbeat := time.Now()
	for {
		var err error
		msg, err = pack.UnpackWs()
		if err != nil {
			// 发送失败
			c.conn = nil
			log.Warnln("[mys-ws]", c.wr.Data.AppId, "的网关连接断开, 尝试恢复:", err)
			time.Sleep(time.Second)
			//开始重连
			c.Resume()
			//log.Warnln("[ws]", c.wr.Data.AppId, "的网关连接恢复失败:", err)
			continue
		}
		counter.Inc()
		log.Debugln("[mys-ws]接收到第", counter.Load(), "个事件", "类型:", msg.Type, ", 数据:", tool.BytesToString(msg.Body))
		//接收的类型进行选择
		switch vila_bot.Command(msg.Type) {
		case vila_bot.Command_P_HEARTBEAT: // Send/Receive

			log.Debugln("[ws]收到服务端推送心跳, 间隔:", time.Since(lastheartbeat))
			lastheartbeat = time.Now()
		case vila_bot.Command_P_KICK_OFF:
			//PKickOff
		case vila_bot.Command(30001): // Receive/Reply
			//解析消息
		default:
			log.Warnln("[ws]未知事件, 序号:", counter.Load(), "类型:", msg.Type, ", 数据:", tool.BytesToString(msg.Body))
		}
	}
}
