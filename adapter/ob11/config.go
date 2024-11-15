package ob11

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

var botMap = map[string]*Config{}

type Config struct {
	*rosm.BotCard
	URL    string `json:"url"`
	Token  string `json:"access_token"`
	Types  string `json:"types"`
	Driver Driver `json:"-"` // 通信驱动
}

func (c *Config) Card() *rosm.BotCard {
	return c.BotCard
}

func NewConfig(path string) (c *Config) {
	c = &Config{
		BotCard: &rosm.BotCard{
			Master:  []string{"123456"},
			BotName: "雪儿",
			BotID:   "1",
		},
		URL:   "ws://127.0.0.1:6700",
		Types: "WS",
	}
	rosm.LoadBotConfig(path, c)
	return c
}
