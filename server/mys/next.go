package mys

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

var nextEmoticonList = map[string]chan *rosm.CTX{}

// 获取本消息-全体的下一表态
func GetNextAllEmoticon(ctx *rosm.CTX, botMsgID string) (chan *rosm.CTX, func()) {
	next := make(chan *rosm.CTX, 1)
	nextEmoticonList[botMsgID] = next
	return next, func() {
		close(next)
		delete(nextEmoticonList, botMsgID)
	}
}

func emoticonNext(ctx *rosm.CTX) {
	if len(nextEmoticonList) == 0 {
		return
	}
	if c, ok := nextEmoticonList[ctx.Message.(*InfoSTR).Event.ExtendData.EventData.AddQuickEmoticon.BotMsgID]; ok {
		c <- ctx
	}
}
