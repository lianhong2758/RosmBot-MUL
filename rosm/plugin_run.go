// 事件的匹配实现
package rosm

import (
	"runtime/debug"

	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

// 匹配事件
func (ctx *Ctx) RunEvent(types int) (block bool) {
	defer func() {
		if pa := recover(); pa != nil {
			log.Errorf("[rosm] RunEvent Err: %v\n%v", pa, tool.BytesToString(debug.Stack()))
		}
	}()
	log.Debug("[Event]开始匹配事件type", types)
	if ctx.sendNext(types) {
		return true
	}
	ctx.on = PluginIsOn(boten)(ctx)
	for _, m := range EventMatch[types] {
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
	defer func() {
		if pa := recover(); pa != nil {
			log.Errorf("[rosm] RunEvent Err: %v\n%v", pa, tool.BytesToString(debug.Stack()))
		}
	}()
	ctx.Being.Word = word
	//全匹配
	if ctx.RunEvent(AllMessage) {
		return
	}
	ctx.on = PluginIsOn(boten)(ctx)
	//关键词触发
	if m, ok := WordMatch[word]; ok {
		if m.RulePass(ctx) {
			m.handler(ctx)
			log.Debugf("调用插件: %s - 匹配关键词: %s", m.PluginNode.Name, word)
		}
		return
	}
	//正则匹配
	for _, m := range RegexpMatch {
		regex := m.Rex
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
