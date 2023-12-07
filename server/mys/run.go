package mys

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

func (c *Config) Run() {
	switch c.Protocol {
	case "http":
		c.RunHTTP()
	case "ws":
		c.RunWS()
	default:
		log.Error("[mys]协议错误")
	}
}

func (c *Config) mul() {
	rosm.MULChan <- rosm.MUL{Types: "mys", Name: c.BotToken.BotName, BotID: c.BotToken.BotID}
}
func (c *Config) RunHTTP() {
	c.mul()
	botMap[c.BotToken.BotID] = c
	gin.SetMode(gin.ReleaseMode)
	r := gin.New() //初始化
	log.Infoln("[mys-http]bot开始监听消息")
	r.POST(c.EventPath, c.MessReceive())
	r.Run(c.Port)
}

func (c *Config) RunWS() {
	c.mul()
	botMap[c.BotToken.BotID] = c
	for {
		if err := c.GetWebsocketUrl(); err != nil || c.wr.Retcode != 0 {
			log.Errorln("[mys-ws]获取WebsocketUrl失败,ERROR:", err, "Message:", c.wr.Message)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
	c.Login()
	c.Listen()
}
