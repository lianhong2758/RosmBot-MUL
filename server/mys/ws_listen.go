package mys

import (
	"time"

	vila_bot "github.com/lianhong2758/RosmBot-MUL/server/mys/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
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
		log.Debugln("[mys-ws]接收到第", counter.Load(), "个事件", "类型:", msg.Type)
		//接收的类型进行选择
		switch vila_bot.Command(msg.Type) {
		case vila_bot.Command_P_HEARTBEAT:
			//心跳
			t := &vila_bot.PHeartBeatReply{}
			_ = proto.Unmarshal(msg.Body, t)
			log.Debugln("[mys-ws]收到服务端推送心跳", t)
		case vila_bot.Command_P_LOGOUT:
			//长连接登出协议(ProtoBuf)
			t := &vila_bot.PLogoutReply{}
			_ = proto.Unmarshal(msg.Body, t)
			log.Debugln("[mys-ws]登出返回", t)
			if t.Code != 0 {
				log.Warnln("[mys-ws]登出请求失败:", t.Msg)
				continue
			}
			log.Info("[mys-ws]已请求登出...")
			return
		case vila_bot.Command_SHUTDOWN:
			// 	shutdown(重启)
			go c.RunWS()
			return
		case vila_bot.Command_P_KICK_OFF:
			// 	踢下线协议（ProtoBuf）
			t := &vila_bot.PKickOff{}
			_ = proto.Unmarshal(msg.Body, t)
			log.Debugln("[mys-ws]PKickOff", t)
			log.Infoln("[mys-ws]Code: ", t.Code, "Reason", t.Reason)
			return
		case vila_bot.Command(30001): // Receive/Reply
			// 	机器人开放平台事件（protoBuf）
			t := &vila_bot.RobotEvent{}
			_ = proto.Unmarshal(msg.Body, t)
			log.Debugln("[mys-ws]RobotEvent", t)
			c.process(t)
		default:
			log.Warnln("[mys-ws]未知事件, 序号:", counter.Load(), "类型:", msg.Type, ", 数据:", msg.Body)
		}
	}
}
