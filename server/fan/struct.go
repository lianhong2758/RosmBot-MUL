package fan

import "encoding/json"

type H = map[string]any

// 回调的请求结构

type User struct {
	Ok     bool `json:"ok"`
	Result struct {
		ID                      int64  `json:"id"`
		IsBot                   bool   `json:"is_bot"`
		FirstName               string `json:"first_name"`
		LastName                string `json:"last_name"`
		Username                string `json:"username"`
		Avatar                  string `json:"avatar"`
		UserToken               string `json:"user_token"`
		OwnerID                 int    `json:"owner_id"`
		CanJoinGroups           bool   `json:"can_join_groups"`
		CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
		SupportsInlineQueries   bool   `json:"supports_inline_queries"`
	} `json:"result"`
}

type Message struct {
	ID   int64 `json:"message_id"`
	User *User `json:"from"`
	Time int64 `json:"date"`
	Chat *Chat `json:"chat"`
}

type Chat struct {
	ID        int64  `json:"id"`                   //id
	Type      string `json:"type"`                 // 聊天类型，可以是 "private", "channel", "group", "supergroup"
	Title     string `json:"title,omitempty"`      // 聊天标题，可选字段
	Username  string `json:"username,omitempty"`   // Fanbook的ID，可选字段
	FirstName string `json:"first_name,omitempty"` // 当前填充的用户的昵称，可选字段
}

type WsReq struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type PushData struct {
	Content     string      `json:"content"`
	Time        int64       `json:"time"`
	UserID      string      `json:"user_id"`
	ChannelID   string      `json:"channel_id"`
	MessageID   string      `json:"message_id"`
	QuoteL1     interface{} `json:"quote_l1"`
	QuoteL2     interface{} `json:"quote_l2"`
	GuildID     string      `json:"guild_id"`
	ChannelType int         `json:"channel_type"`
	Status      int         `json:"status"`
	Nonce       string      `json:"nonce"`
	Ctype       int         `json:"ctype"`
	Member      struct {
		Nick        interface{} `json:"nick"`
		Roles       []string    `json:"roles"`
		GuildCard   interface{} `json:"guild_card"`
		AssistLevel int         `json:"assist_level"`
	} `json:"member"`
	ResourceType string `json:"resource_type"`
	Platform     string `json:"platform"`
	Author       struct {
		Nickname  string      `json:"nickname"`
		Username  string      `json:"username"`
		Avatar    string      `json:"avatar"`
		AvatarNft interface{} `json:"avatar_nft"`
		Bot       bool        `json:"bot"`
	} `json:"author"`
	Desc string `json:"desc"`
}
