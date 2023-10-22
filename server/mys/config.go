package mys

// 默认值
var MYSconfig = Config{
	BotToken: Token{
		Master: []string{"123456"},
	}}

type Token struct {
	Master         []string `json:"master_id"`
	BotID          string   `json:"bot_id"`
	BotSecret      string   `json:"bot_secret"`
	BotPubKey      string   `json:"bot_pub_key"`
	BotName        string   `json:"bot_name"`
	BotSecretConst string   `json:"-"`
}
type Config struct {
	BotToken  Token  `json:"token"`
	EventPath string `json:"eventpath,omitempty"`
	Port      string `json:"port,omitempty"`
}

func (c *Config) Name() string {
	return MYSconfig.BotToken.BotName
}
