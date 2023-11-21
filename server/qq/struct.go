package qq

import (
	"encoding/json"
	"time"
)

type H = map[string]any

// AccessTokenStr 接口凭证回调
type AccessTokenStr struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

// CodeMessageBase 各种消息都有的 code + message 基类
type CodeMessageBase struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ShardWSSGateway 带分片 WSS 接入点响应数据
//
// https://bot.q.qq.com/wiki/develop/api/openapi/wss/shard_url_get.html#%E8%BF%94%E5%9B%9E
type ShardWSSGateway struct {
	URL               string `json:"url"`
	Shards            int    `json:"shards"`
	SessionStartLimit struct {
		Total          int `json:"total"`
		Remaining      int `json:"remaining"`
		ResetAfter     int `json:"reset_after"`
		MaxConcurrency int `json:"max_concurrency"`
	} `json:"session_start_limit"`
}

// WebsocketPayload payload 指的是在 websocket 连接上传输的数据，网关的上下行消息采用的都是同一个结构
//
// https://bot.q.qq.com/wiki/develop/api/gateway/reference.html
type WebsocketPayload struct {
	Op OpCode          `json:"op"`
	D  json.RawMessage `json:"d,omitempty"`
	S  uint32          `json:"s,omitempty"`
	T  string          `json:"t,omitempty"`
}

// guild message
type RawGuildMessage struct {
	Author    *User  `json:"author"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
	GuildID   string `json:"guild_id"`
	ID        string `json:"id"`
	Member    struct {
		JoinedAt time.Time `json:"joined_at"`
		Nick     string    `json:"nick"`
		Roles    []string  `json:"roles"`
	} `json:"member"`
	Mentions     []Mention `json:"mentions"`
	Seq          int       `json:"seq"`
	SeqInChannel string    `json:"seq_in_channel"`
	Timestamp    time.Time `json:"timestamp"`
}

type Mention struct {
	Avatar   string `json:"avatar"`
	Bot      bool   `json:"bot"`
	ID       string `json:"id"`
	Username string `json:"username"`
}

type RawGroupMessage struct {
	Author struct {
		ID           string `json:"id"`
		MemberOpenid string `json:"member_openid"`
	} `json:"author"`
	Content     string    `json:"content"`
	GroupID     string    `json:"group_id"`
	GroupOpenid string    `json:"group_openid"`
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
}
