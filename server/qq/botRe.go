package qq

// Authorization 返回 Authorization Header value
func (c *Config) Authorization() string {
	return "Bot " + c.BotToken.AppId + "." + c.BotToken.Token
}

// AtMe 返回 "<@!"+bot.ready.User.ID+">"
func (c *Config) AtMe() string {
	return "<@!" + c.Ready.User.ID + ">"
}
