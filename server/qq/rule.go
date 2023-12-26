package qq

import "github.com/lianhong2758/RosmBot-MUL/rosm"

func (c *Config) OnlyReply(ctx *rosm.Ctx) bool {
	switch ctx.Being.Def["type"].(string) {
	case "C2C_MESSAGE_CREATE":
		return false
	case "GROUP_AT_MESSAGE_CREATE":
		return false
	case "DIRECT_MESSAGE_CREATE", "AT_MESSAGE_CREATE", "MESSAGE_CREATE":
		return ctx.Message.(*RawGuildMessage).MessageReference.MessageID != ""
	default:
		return false
	}
}

// 主人
func (c *Config) OnlyMaster(ctx *rosm.Ctx) bool {
	switch ctx.Being.Def["type"].(string) {
	case "C2C_MESSAGE_CREATE":
		return false
	case "GROUP_AT_MESSAGE_CREATE":
		return false
	case "DIRECT_MESSAGE_CREATE", "AT_MESSAGE_CREATE", "MESSAGE_CREATE":
		for _, v := range c.Master {
			if v == ctx.Being.User.ID {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// 群主权限 未实现
func (c *Config) OnlyOverHost(ctx *rosm.Ctx) bool {
	return false
}

// 管理员 未实现
func (c *Config) OnlyOverAdministrator(ctx *rosm.Ctx) bool {
	return false
}
