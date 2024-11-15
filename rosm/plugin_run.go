// 事件的匹配实现
package rosm

import (
	"runtime/debug"

	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

// 匹配事件
func (ctx *Ctx) RunEvent(types EventType) (block bool) {
	defer func() {
		if pa := recover(); pa != nil {
			log.Errorf("[rosm] RunEvent Err: %v\n%v", pa, tool.BytesToString(debug.Stack()))
		}
	}()
	//bot启用状态
	ctx.on = PluginIsOn(boten)(ctx)
	//next
	if ctx.sendNext(types) {
		return true
	}
	//event
	for _, m := range EventMatch[types] {
		if m.RulePass(ctx) {
			m.handler(ctx)
			log.Debugf("调用插件: %s - 类型: %s", m.PluginNode.Name, types)
			return m.block
		}
	}
	return false
}

// 匹配修剪好的触发词
func (ctx *Ctx) RunWord() {
	defer func() {
		if pa := recover(); pa != nil {
			log.Errorf("[rosm] RunEvent Err: %v\n%v", pa, tool.BytesToString(debug.Stack()))
		}
	}()
	//全匹配
	if ctx.RunEvent("all") {
		return
	}
	//关键词触发,不检查block
	if m, ok := WordMatch[ctx.Being.RawWord]; ok {
		if m.RulePass(ctx) {
			m.handler(ctx)
		}
		return
	}
	//正则匹配
	for _, m := range RegexpMatch {
		if m.RulePass(ctx) {
			if m.handler(ctx); m.block {
				return
			}
		}
	}
	//未匹配时触发
	ctx.RunEvent("surplus")
}
