// 事件的匹配实现
package rosm

import (
	log "github.com/sirupsen/logrus"
)

// 匹配事件
func (ctx *Ctx) RunEvent(types int) (block bool) {
	log.Debug("[Event]开始匹配事件type", types)
	if ctx.sendNext(types) {
		return true
	}
	for _, m := range caseEvent[types] {
		if m.RulePass(ctx) {
			m.handler(ctx)
			log.Debugf("调用插件: %s - 类型: %d", m.PluginNode.Name, types)
			return m.block
		}
	}
	return false
}

// 匹配修剪好的触发词
func (ctx *Ctx) RunWord(word string) {
	ctx.Being.Word = word
	//全匹配
	if ctx.RunEvent(AllMessage) {
		return
	}
	//关键词触发
	if m, ok := caseAllWord[word]; ok {
		if m.RulePass(ctx) {
			m.handler(ctx)
			log.Debugf("调用插件: %s - 匹配关键词: %s", m.PluginNode.Name, word)
		}
		return
	}
	//正则匹配
	for regex, m := range caseRegexp {
		if match := regex.FindStringSubmatch(word); len(match) > 0 {
			if m.RulePass(ctx) {
				ctx.Being.Rex = match
				m.handler(ctx)
				log.Debugf("调用插件: %s - 匹配关键词: %s", m.PluginNode.Name, word)
				if m.block {
					return
				}
			}
		}
	}
	//未匹配时触发
	ctx.RunEvent(SurplusMessage)
}
