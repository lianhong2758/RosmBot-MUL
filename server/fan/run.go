package fan

import (
	"time"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

func (c *Config) Run() {
	c.RunWS()
}

func (c *Config) mul() {
	rosm.MULChan <- rosm.MUL{Types: "fan", Name: c.BotName, BotID: c.BotID}
}
func (c *Config) RunWS() {
	for {
		if err := c.GetBotData(); err != nil {
			log.Errorln("[fan]获取长连接失败,ERROR: ", err)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
	c.setInit()
	c.mul()
	botMap[c.BotID] = c
	c.ListenWS()
}
func (c *Config) setInit() {
	c.BotID = tool.Int64ToString(c.user.Result.ID)
}
