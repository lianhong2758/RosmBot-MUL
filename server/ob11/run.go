package ob11

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func (config *Config) Run() {
	config.mul()
	switch config.Types {
	case "WS":
		config.Driver = driver.NewWebSocketClient(config.URL, config.Token)
	case "WSS":
		config.Driver = driver.NewWebSocketServer(16, config.URL, config.Token)
	}
	config.Driver.Connect()
	config.Driver.Listen(config.processEvent())
}

func (c *Config) mul() {
	rosm.MULChan <- rosm.MUL{Types: "ob11", Name: c.BotName, BotID: c.BotID}
}
