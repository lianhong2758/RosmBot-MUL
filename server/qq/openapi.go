package qq

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	log "github.com/sirupsen/logrus"
)

func makeHeard(ACCESS_TOKEN, APPID string) func(*http.Request) {
	return func(req *http.Request) {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", ACCESS_TOKEN)
		req.Header.Add("X-Union-Appid", APPID)
	}
}

// 通过链接获取data并写入结构体
func (c *Config) GetOpenAPI(shortUrl string, body, result any) (err error) {
	var data []byte
	if body != nil && body != "" && body != 0 {
		data, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	data, err = web.Web(clientConst, host+shortUrl, http.MethodGet, makeHeard(c.access, c.BotToken.AppId), bytes.NewReader(data))
	log.Debugln("[GetOpenAPI]", host+shortUrl, err)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}
