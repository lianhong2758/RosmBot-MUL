package gscore

import (
	"context"
	"encoding/json"
	"os"

	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
)

var Config *GsConfig
var path = "data/gscore/config.json"

func configInit() {
	data, err := os.ReadFile(path)
	if err != nil {
		Config = new(GsConfig)
		//初始配置
		Config.CoreUrl = "ws://127.0.0.1:8765/ws/rosmbot"
		Config.CommandPrefix = ""

		data, _ = json.MarshalIndent(Config, "", "  ")
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			log.Error("[gscore]创建config失败: ", err)
		}
		log.Info("[gscore]保存json成功")
	}
	Config = new(GsConfig)
	err = json.Unmarshal(data, Config)
	if err != nil {
		log.Error(err)
	}
}

type GsConfig struct {
	on     bool
	conn   *websocket.Conn
	cancel context.CancelFunc

	CoreUrl       string `description:"连接Core的Url"`
	CommandPrefix string `default:"/" description:"命令前缀"`
}
