package ob11

import (
	"encoding/json"
	"os"
	"time"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/ob11/driver"
	log "github.com/sirupsen/logrus"
)

// Config
type Config struct {
	*rosm.BotCard
	MaxProcessTime time.Duration      `json:"max_process_time"` // 事件最大处理时间 (默认4min)
	W              []*driver.WSClient `json:"ws"`
	S              []*driver.WSServer `json:"wss"`
	Driver         []driver.Driver    `json:"-"` // 通信驱动
}

func NewConfig(path string) (c *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		c = new(Config)
		c.BotCard = new(rosm.BotCard)
		c.Master = []string{"123456"}
		c.BotName = "雪儿"
		c.W = []*driver.WSClient{driver.NewWebSocketClient("ws://127.0.0.1:6700", "")}
		err = os.MkdirAll("config", os.ModePerm)
		if err != nil {
			log.Fatalln("[ob11]无法创建 config 目录: ", err)
		}
		data, _ = json.MarshalIndent(c, "", "  ")
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			log.Fatalln("[ob11]创建config失败: ", err)
		}
		log.Fatalln("[ob11]创建初始配置完成,请填写config中的配置文件后再启动本程序")
	}
	c = new(Config)
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

func (c *Config) Card() *rosm.BotCard { return c.BotCard }
