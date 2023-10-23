package qq

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
)

var botMap map[string]*Config

type Token struct {
	AppId     string `json:"AppId"`
	AppSecret string `json:"AppSecret"`
	Token     string `json:"Token,omitempty"`
}
type Config struct {
	BotToken   Token    `json:"token"`
	Master     []string `json:"master_id"`
	BotName    string   `json:"bot_name"`
	Intents    uint32   `json:"-"`           // Intents 欲接收的事件
	IntentsNum []uint32 `json:"intents"`     //用户输入的
	ShardIndex uint16   `json:"shard_index"` //分片序号

	access    string
	shard     [2]byte         // shard 分片
	gateway   string          // gateway 获得的网关
	seq       uint32          // seq 最新的 s
	heartbeat uint32          // heartbeat 心跳周期, 单位毫秒
	mu        sync.Mutex      // 写锁
	conn      *websocket.Conn // conn 目前的 wss 连接
	hbonce    sync.Once       // hbonce 保证仅执行一次 heartbeat
	Ready     EventReady
}

// 下发的bot配置
type EventReady struct {
	Version   int     `json:"version"`
	SessionID string  `json:"session_id"`
	User      *User   `json:"user"`
	Shard     [2]byte `json:"shard"`
}

// User 用户对象
//
// https://bot.q.qq.com/wiki/develop/api/openapi/user/model.html
type User struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	Avatar           string `json:"avatar"`
	Bot              bool   `json:"bot"`
	UnionOpenid      string `json:"union_openid"`
	UnionUserAccount string `json:"union_user_account"`
}

func (c *Config) Name() string {
	return c.BotName
}

func NewConfig(path string) (c *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		c = new(Config)
		c.Master = []string{"123456"}
		c.BotName = "雪儿"
		c.IntentsNum = []uint32{0, 9, 30, 25}
		c.ShardIndex = 0
		err = os.MkdirAll("config", os.ModePerm)
		if err != nil {
			log.Fatalln("[qq]无法创建 config 目录: ", err)
		}
		data, _ = json.MarshalIndent(c, "", "  ")
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			log.Fatalln("[qq]创建config失败: ", err)
		}
		log.Fatalln("[qq]创建初始配置完成,请填写config中的配置文件后再启动本程序")
	}
	c = new(Config)
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Fatalln(err)
	}
	if c.BotToken.AppId == "" || c.BotToken.Token == "" {
		log.Fatalln("[qq]未设置bot信息")
	}
	return c
}
