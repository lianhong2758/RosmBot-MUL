package ob11

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func (config *Config) Run() {
	switch config.Types {
	case "WS":
		config.Driver = NewWebSocketClient(config.URL, config.Token)
	case "WSS":
		//wss有问题
		config.Driver = NewWebSocketServer(16, config.URL, config.Token)
	}
	config.Driver.Connect(config) //连接
	config.mul()
	config.Driver.Listen(config)
}

func (c *Config) mul() {
	rosm.MULChan <- rosm.MUL{Types: "ob11", Name: c.BotName, BotID: c.BotID}
}
