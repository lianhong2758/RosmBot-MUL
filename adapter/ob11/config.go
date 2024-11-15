package ob11

import (
	"encoding/json"
	"os"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var botMap = map[string]*Config{}

type Config struct {
	*rosm.BotCard
	URL    string      `json:"url"`
	Token  string      `json:"access_token"`
	Types  string      `json:"types"`
	Driver zero.Driver `json:"-"` // 通信驱动
}

func (c *Config) Card() *rosm.BotCard {
	return c.BotCard
}

func NewConfig(path string) (c *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		c = new(Config)
		c.BotCard = new(rosm.BotCard)
		c.Master = []string{"123456"}
		c.BotName = "雪儿"
		c.BotCard.BotID = "1"
		c. URL = "ws://127.0.0.1:6700"
		c.Token=""
		c.Types="WS"
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
