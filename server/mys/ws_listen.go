package mys

import (
	"encoding/json"
	"os"
	"time"

	vila_bot "github.com/lianhong2758/RosmBot-MUL/server/mys/proto"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

func (c *Config) Listen() {
	log.Infoln("[mys-ws]bot开始监听消息")
	pack := NewDataPack(c.conn)
	var msg *Message
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
			var t map[string]any
			json.Unmarshal(msg.Body, &t)
			log.Debugln("[ws]收到服务端推送心跳", t["server_timestamp"])
		case vila_bot.Command_P_LOGIN, vila_bot.Command_P_LOGOUT:
			// 	长连接登录协议(ProtoBuf) ,长连接登出协议(ProtoBuf)
		case vila_bot.Command_SHUTDOWN:
			// 	服务关闭协议（无字段）
		case vila_bot.Command_P_KICK_OFF:
			// 	踢下线协议（ProtoBuf）
			log.Error("已在另一设备登录")
			os.Exit(0)
		case vila_bot.Command(30001): // Receive/Reply
			// 	机器人开放平台事件（protoBuf）
			c.process(msg.Body)
		default:
			log.Warnln("[ws]未知事件, 序号:", counter.Load(), "类型:", msg.Type, ", 数据:", tool.BytesToString(msg.Body))
		}
	}
}
