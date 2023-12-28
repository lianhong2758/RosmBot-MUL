// Rule的实现,可以在这里增加更多Rule,也可以在server包增加独属于自己平台的rule
package rosm

// 判断rule
func (m *Matcher) RulePass(ctx *Ctx) bool {
	return rulePass(ctx, m.rules...)
}

func rulePass(ctx *Ctx, rs ...Rule) bool {
	for _, v := range rs {
		if !v(ctx) {
			return false
		}
	}
	return true
}

func OnlyAtMe() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Being.AtMe
	}
}

func OnlyMaster() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Bot.OnlyMaster(ctx)
	}
}

// 用于getnext
func OnlyTheUser(id string) Rule {
	return func(ctx *Ctx) bool {
		return id == ctx.Being.User.ID
	}
}

func OnlyReply() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Bot.OnlyReply(ctx)
	}
}

// 大于等于群主等权限
func OnlyOverHost() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Bot.OnlyOverHost(ctx)
	}
}

// 大于等于管理员等权限
func OnlyOverAdministrator() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Bot.OnlyOverAdministrator(ctx)
	}
}
