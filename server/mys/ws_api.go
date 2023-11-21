package mys

import (
	"encoding/json"
	"net/http"

	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

// 获取ws地址
func (c *Config) GetWebsocketUrl() error {
	data, err := web.Web(web.NewDefaultClient(), getWSurl, http.MethodGet, func(r *http.Request) {
		r.Header.Add("x-rpc-bot_id", c.BotToken.BotID)
		r.Header.Add("x-rpc-bot_secret", c.BotToken.BotSecret)
		r.Header.Add("x-rpc-bot_villa_id", "2077")
		r.Header.Add("x-rpc-bot_ts", "ts")
		r.Header.Add("x-rpc-bot_nonce", "")
		r.Header.Add("Content-Type", "application/json")
	}, nil)
	if err != nil {
		return err
	}
	wr := new(WebsocketInfoResp)
	c.wr = wr
	return json.Unmarshal(data, &wr)
}

type WebsocketInfoResp struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		WebsocketUrl string `json:"websocket_url"`
		Uid          uint64 `json:"uid,string"`
		AppId        int32  `json:"app_id"`
		Platform     int32  `json:"platform"`
		DeviceId     string `json:"device_id"`
	} `json:"data"`
}
