package qq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
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

// 获取频道用户信息
func GetGuildUser(ctx *rosm.CTX, uid string) (User *GuildUser, err error) {
	url := host + fmt.Sprintf(urlGuildGetUser, ctx.Being.RoomID2, uid)
	data, err := web.Web(clientConst, url, http.MethodGet, makeHeard(ctx.Bot.(*Config).access, ctx.Bot.(*Config).BotToken.AppId), nil)
	log.Debugln("[GetGuildUser][", url, "]", string(data))
	if err != nil {
		return nil, err
	}
	User = new(GuildUser)
	err = json.Unmarshal(data, &User)
	return
}

type GuildUser struct {
	User struct {
		ID               string `json:"id"`
		Username         string `json:"username"`
		Avatar           string `json:"avatar"`
		Bot              bool   `json:"bot"`
		UnionOpenid      string `json:"union_openid"`
		UnionUserAccount string `json:"union_user_account"`
	} `json:"user"`
	Nick     string    `json:"nick"`
	Roles    []string  `json:"roles"`
	JoinedAt time.Time `json:"joined_at"`
}
