package qq

import (
	"sync"

	"github.com/RomiChan/websocket"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

var botMap = map[string]*Config{}

type Token struct {
	AppId     string `json:"AppId"`
	AppSecret string `json:"AppSecret"`
	Token     string `json:"Token"`
}
type Config struct {
	*rosm.BotCard
	BotToken   Token    `json:"token"`
	Intents    uint32   `json:"-"`           // Intents 欲接收的事件
	IntentsNum []uint32 `json:"intents"`     //按照类型填写接收事件
	ShardIndex uint16   `json:"shard_index"` //分片序号

	access    string
	shard     [2]byte         // shard 分片
	gateway   string          // gateway 获得的网关
	seq       uint32          // seq 最新的 s
	heartbeat uint32          // heartbeat 心跳周期, 单位毫秒
	mu        sync.Mutex      // 写锁
	conn      *websocket.Conn // conn 目前的 wss 连接
	hbonce    sync.Once       // hbonce 保证仅执行一次 heartbeat
	Ready     EventReady      `json:"-"`
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

func (c *Config) Card() *rosm.BotCard {
	return c.BotCard
}
func (c *Config) GetReady() EventReady {
	return c.Ready
}

func NewConfig(path string) (c *Config) {
	c = &Config{
		BotCard: &rosm.BotCard{
			Master:  []string{"123456"},
			BotName: "雪儿",
			BotID:   "123456",
		},
		IntentsNum: []uint32{0, 1, 12, 30, 25},
		ShardIndex: 0,
	}
	rosm.LoadBotConfig(path, c)
	if c.BotToken.AppId == "" || c.BotToken.Token == "" {
		log.Fatalln("[qq]未设置bot信息")
	}
	return c
}
