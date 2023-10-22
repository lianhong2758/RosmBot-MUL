package rosm

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
)

// bot主接口
type Boter interface {
	//获取config
	BotSend(*CTX, ...message.MessageSegment) any
	//运行
	Run()

	//Bot信息查询
	Name() string
}

// 进行一个通道注册,同于接收平台注册消息进行统计
var MULChan = make(chan MUL)

type MUL struct {
	Types string
	Name  string
	BotID string
}
