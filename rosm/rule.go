package rosm

func (m *Matcher) RulePass(ctx *CTX) bool {
	for _, v := range m.rules {
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
