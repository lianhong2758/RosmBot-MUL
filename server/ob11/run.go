package ob11

import (
	"strconv"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func (config *Config) Run() {
	switch config.Types {
	case "WS":
		config.Driver = driver.NewWebSocketClient(config.URL, config.Token)
	case "WSS":
		config.Driver = driver.NewWebSocketServer(16, config.URL, config.Token)
	}
	config.Driver.Connect() //连接
	zero.APICallers.Range(func(key int64, value zero.APICaller) bool {
		switch v := value.(type) {
		case *driver.WSClient:
			if v.Url == config.URL {
				config.BotID = strconv.FormatInt(key, 10)
				botMap[strconv.FormatInt(key, 10)] = config
				return false
			}
		//反向	case *driver.WSServer:
		default:
			logrus.Warn("RunError:未适配的WSS")
			return true
		}
		return true
	})
	config.mul()
	config.Driver.Listen(config.processEvent())
}

func (c *Config) mul() {
	rosm.MULChan <- rosm.MUL{Types: "ob11", Name: c.BotName, BotID: c.BotID}
}
