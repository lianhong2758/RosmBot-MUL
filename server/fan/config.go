package fan

import (
	"encoding/json"
	"os"

	"github.com/RomiChan/websocket"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

var botMap = map[string]*Config{}

type Config struct {
	BotToken string `json:"token"`
	*rosm.BotCard
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
	Timeout int    `json:"timeout"`
	MsgType string `json:"msg_type"`

	conn *websocket.Conn 
	user *User
}

func (c *Config) Card() *rosm.BotCard {
	return c.BotCard
}
func NewConfig(path string) (c *Config) {
	c = new(Config)
	data, err := os.ReadFile(path)
	if err != nil {
		c.BotCard = new(rosm.BotCard)
		//初始配置
		c.Master = []string{"123456"}
		c.BotName = "雪儿"
		c.BotToken = "token填写"
		err = os.MkdirAll("config", os.ModePerm)
		if err != nil {
			log.Fatalln("[fan]无法创建 config 目录: ", err)
		}
		data, _ = json.MarshalIndent(c, "", "  ")
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			log.Fatalln("[fan]创建config失败: ", err)
		}
		log.Infoln("创建初始化配置完成\n请填写config/fan.json文件后重新运行本程序")
		os.Exit(0)
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Fatalln(err)
	}
	if c.BotToken == "" {
		log.Fatalln("[fan]未设置bot信息")
	}
	return
}
