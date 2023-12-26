package rosm

// 用于获取下一事件的结构,消息也作为一种事件
//通过rule限制下一条消息

import (
	"time"

	"github.com/sirupsen/logrus"
)

var nextList = map[EventType]map[int]*Matcher{}

// 获取下一事件
func (ctx *Ctx) GetNext(types EventType, SetBlock bool, rs ...Rule) (chan *Ctx, func()) {
	next := make(chan *Ctx, 1)
	ids := int(0xfffffff & time.Now().Unix())
	m := &Matcher{block: SetBlock, rules: rs, nestchan: next}
	if nextList[types] != nil {
		nextList[types][ids] = m
	} else {
		nextList[types] = map[int]*Matcher{ids: m}
	}
	return next, func() {
		close(next)
		delete(nextList[types], ids)
	}
}

func (ctx *Ctx) sendNext(types EventType) (block bool) {
	if len(nextList) == 0 || nextList[types] == nil {
		return false
	}
	logrus.Debug("[next]匹配事件type", types)
	for _, v := range nextList[types] {
		if v.RulePass(ctx) {
			v.nestchan <- ctx
			if v.block {
				return true
			}
		}
	}
	return false
}
