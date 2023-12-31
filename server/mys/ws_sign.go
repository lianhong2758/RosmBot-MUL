package mys

import (
	"time"

	"github.com/RomiChan/websocket"
	vila_bot "github.com/lianhong2758/RosmBot-MUL/server/mys/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func (c *Config) Login() {
	log.Infoln("[mys-sign]开始尝试连接到网关:", c.wr.Data.WebsocketUrl, ",BotID:", c.BotToken.BotID)
	//准备数据结构
	req := &vila_bot.PLogin{
		Uid:      c.wr.Data.Uid,
		Token:    c.TestVilla + "." + c.BotToken.BotSecret + "." + c.BotToken.BotID, //机器人 Websocket 鉴权 token，格式为 {villa_id}.{secret}.{bot_id} 。机器人未上线时，villa_id 使用测试别野，上线后可传 0
		Platform: c.wr.Data.Platform,
		AppId:    c.wr.Data.AppId,
		DeviceId: c.wr.Data.DeviceId,
		Region:   "",

		Meta: nil,
	}
	data, _ := proto.Marshal(req)

	for {
		conn, resp, err := websocket.DefaultDialer.Dial(c.wr.Data.WebsocketUrl, nil)
		if err != nil {
			log.Warnf("[mys-sign]连接到网关 %v 时出现错误: %v", c.wr.Data.WebsocketUrl, err)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		c.conn = conn
		_ = resp.Body.Close()
		err = c.sendTextMsg(uint32(vila_bot.Command_P_LOGIN), data) // 发送登录请求
		if err != nil {
			log.Warnf("[mys-sign]尝试Login错误: %v", err)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		log.Infof("[mys-sign]登录中...")
		//接收下一条消息
		res, err := NewDataPack(c.conn).UnpackWs()
		if err != nil {
			log.Warnf("[mys-sign]获取Login回复错误: %v", err)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		loginReply := new(vila_bot.PLoginReply)
		_ = proto.Unmarshal(res.GetBody(), loginReply)
		if loginReply.GetCode() != 0 {
			log.Warn("[mys-sign]登录失败...,Code: ", loginReply.GetCode())
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		break
	}
	log.Infoln("[mys-sign]连接到网关成功, 用户:", c.wr.Data.AppId)
	c.hbonce.Do(func() {
		go c.repHeart()
	})
}

// 发送消息
func (c *Config) sendTextMsg(bizType uint32, data []byte) error {
	return NewDataPack(c.conn).Pack(NewRequestMsg(bizType, UniqMsgID(), c.wr.Data.AppId, data)).Send()
}

func (c *Config) repHeart() {
	for {
		beatRequest := vila_bot.PHeartBeat{ClientTimestamp: time.Now().String()}
		data, _ := proto.Marshal(&beatRequest)

		err := c.sendTextMsg(uint32(vila_bot.Command_P_HEARTBEAT), data)
		if err != nil {
			// handle err
			log.Warn("[mys-ws]发送心跳失败...", err)
		}
		time.Sleep(20 * time.Second)
	}
}

// 重连
func (c *Config) Resume() {
	c.Login()
}

// 登出
func (c *Config) LogOut() {
	req := &vila_bot.PLogout{
		Uid:      c.wr.Data.Uid,
		Platform: c.wr.Data.Platform,
		AppId:    c.wr.Data.AppId,
		DeviceId: c.wr.Data.DeviceId,
	}
	data, _ := proto.Marshal(req)
	err := c.sendTextMsg(uint32(vila_bot.Command_P_LOGOUT), data) // 发送登录请求
	if err != nil {
		log.Warnf("[mys-sign]尝试PLogout错误: %v", err)
	}
}
