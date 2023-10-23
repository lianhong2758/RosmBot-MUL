package mys

import (
	"github.com/gin-gonic/gin"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

func (c *Config) Run() {
	c.mul()
	botMap[c.BotToken.BotID] = c
	gin.SetMode(gin.ReleaseMode)
	r := gin.New() //初始化
	log.Infoln("[mys-http]bot开始监听消息")
	r.POST(c.EventPath, c.MessReceive())
	r.Run(c.Port)
}

func (c *Config) mul() {
	rosm.MULChan <- rosm.MUL{Types: "mys", Name: c.BotToken.BotName, BotID: c.BotToken.BotID}
}
