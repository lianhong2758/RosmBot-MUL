package ob11

import "github.com/lianhong2758/RosmBot-MUL/rosm"

// 判断回复消息
func (c *Config) OnlyReply(ctx *rosm.Ctx) bool { return false }

// 判断主人
func (c *Config) OnlyMaster(ctx *rosm.Ctx) bool { return false }

// 判断群主等
func (c *Config) OnlyOverHost(ctx *rosm.Ctx) bool { return false }

// 判断管理员
func (c *Config) OnlyOverAdministrator(ctx *rosm.Ctx) bool { return false }
