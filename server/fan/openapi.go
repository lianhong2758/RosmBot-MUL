package fan

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

const (
	Host = "https://a1.fanbook.cn"
	URI  = "/api/bot/%s/%s"
	//get
	UrlGetBot = "/api/bot/%s/getMe"
	//post
	UrlGetUpdates = "/api/bot/%s/getUpdates"
	//ws
	UrlWs="wss://gateway-bot.fanbook.mobi/websocket?id=%s&dId=%s&v=1.6.60&x-super-properties=%s"
)

// 获取ws地址
func (c *Config) GetBotData() error {
	data, err := web.Web(web.NewDefaultClient(), Host+fmt.Sprintf(UrlGetBot, c.BotToken), http.MethodGet, func(r *http.Request) {
		r.Header.Add("Content-Type", "application/json")
	}, nil)
	if err != nil {
		return err
	}
	bot := new(User)
	c.user = bot
	return json.Unmarshal(data, bot)
}
