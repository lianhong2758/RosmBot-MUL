package mys

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

// 是否是主人权限
func OnlyMaster(ctx *rosm.Ctx) bool {
	for _, v := range ctx.Bot.(*Config).BotToken.Master {
		if v == ctx.Being.User.ID {
			return true
		}
	}
	return false
}

// 别野房东权限以上
func OnlyOverHost(ctx *rosm.Ctx) bool {
	data, err := GetVillaData(ctx)
	if err != nil {
		log.Errorln("[ctx](", ctx.Being.RoomID2, ")获取别野信息失败:", err)
	}
	return ctx.Being.User.ID == data.Data.Villa.OwnerUID || OnlyMaster(ctx)
}

// 触发消息是否是回复消息
func OnlyReply(ctx *rosm.Ctx) bool {
	return ctx.Being.Def["Content"].(*MessageContent).Quote.OriginalMessageSendTime != 0 && ctx.Being.Def["Content"].(*MessageContent).Quote.QuotedMessageID != ""
}

func (c *Config) OnlyReply(ctx *rosm.Ctx) bool {
	return ctx.Being.Def["Content"].(*MessageContent).Quote.OriginalMessageSendTime != 0 && ctx.Being.Def["Content"].(*MessageContent).Quote.QuotedMessageID != ""
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
	return false
}
