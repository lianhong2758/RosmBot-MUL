package ob11

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

// 是否是主人权限
func OnlyMaster(ctx *rosm.Ctx) bool {
	for _, v := range ctx.Bot.(*Config).Master {
		if v == ctx.Being.User.ID {
			return true
		}
	}
	return false
}

// 群主权限以上
func OnlyOverHost(ctx *rosm.Ctx) bool {
	return OnlyMaster(ctx) || ctx.State["event"].(*Event).Sender.Role == "owner"
}

// 管理员权限以上
func OnlyOverAdministrator(ctx *rosm.Ctx) bool {
	return OnlyMaster(ctx) || ctx.State["event"].(*Event).Sender.Role == "owner" ||
		ctx.State["event"].(*Event).Sender.Role == "admin"
}

// 触发消息是否是回复消息
func OnlyReply(ctx *rosm.Ctx) bool {
	return ctx.State["reply"].(string) != ""
}

func (c *Config) OnlyReply(ctx *rosm.Ctx) bool {
	return OnlyReply(ctx)
}

// 主人权限
func (c *Config) OnlyMaster(ctx *rosm.Ctx) bool {
	return OnlyMaster(ctx)
}

// 房东权限
func (c *Config) OnlyOverHost(ctx *rosm.Ctx) bool {
	return OnlyOverHost(ctx)
}

// 管理员 未实现
func (c *Config) OnlyOverAdministrator(ctx *rosm.Ctx) bool {
	return OnlyOverAdministrator(ctx)
}
