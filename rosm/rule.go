package rosm

// 判断rule
func (m *Matcher) RulePass(ctx *CTX) bool {
	return rulePass(ctx, m.rules...)
}

func rulePass(ctx *CTX, rs ...Rule) bool {
	for _, v := range rs {
		if !v(ctx) {
			return false
		}
	}
	return true
}

func OnlyAtMe() Rule {
	return func(ctx *CTX) bool {
		return ctx.Being.AtMe
	}
}

func OnlyMaster() Rule {
	return func(ctx *CTX) bool {
		for _, v := range ctx.Bot.Card().Master {
			if v == ctx.Being.User.ID {
				return true
			}
		}
		return false
	}
}

// 用于getnext
func OnlyTheUser(id string) Rule {
	return func(ctx *CTX) bool {
		return id == ctx.Being.User.ID
	}
}
