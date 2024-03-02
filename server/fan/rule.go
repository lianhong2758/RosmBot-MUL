package fan

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
//	log "github.com/sirupsen/logrus"
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
 

func (c *Config) OnlyReply(ctx *rosm.Ctx) bool {
	return  true
}

// 主人权限
func (c *Config) OnlyMaster(ctx *rosm.Ctx) bool {
	return OnlyMaster(ctx)
}

// 房东权限
func (c *Config) OnlyOverHost(ctx *rosm.Ctx) bool {
	return  false
}

// 管理员 未实现
func (c *Config) OnlyOverAdministrator(ctx *rosm.Ctx) bool {
	return false
}
