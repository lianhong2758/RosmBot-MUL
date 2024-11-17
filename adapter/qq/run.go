package qq

import (
	"github.com/lianhong2758/RosmBot-MUL/adapter"
	log "github.com/sirupsen/logrus"
)

// 运行一个bot实例
func (c *Config) Run() {
	c.setInit()
	err := c.getinitinfo()
	if err != nil {
		log.Errorln("QQ-Run", err, "Name: ", c.BotName)
	}
	c.Connect()
	c.BotID= c.BotToken.AppId
	adapter.AddNewBot(c)
	c.Listen()
}
func (c *Config) setInit() {
	c.getAccessToken()
	for _, v := range c.IntentsNum {
		c.Intents = c.Intents | 1<<v
	}
}