package rosm

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	log "github.com/sirupsen/logrus"
)

// Boter bot主接口,由server实现,Run函数由server自己调用运行
type Boter interface {
	//发送消息
	BotSend(*CTX, ...message.MessageSegment) any
	//运行,用于开启接收消息和调用插件
	Run()

	//Bot信息查询
	Card() *BotCard
}
type BotCard struct {
	BotName string   `json:"bot_name"`
	BotID   string   `json:"bot_id,omitempty"`
	Master  []string `json:"master_id"`
}

// 进行一个通道注册,同于接收平台注册消息进行统计,虽然不是强制性需要,但还是建议实现一下
var MULChan = make(chan MUL)

type MUL struct {
	Types string
	Name  string
	BotID string
}

func Listen() {
	plugindbinit()
	for mulData := range MULChan {
		log.Infof("[mul]新增注册,平台: %s,昵称: %s,BotID: %s", mulData.Types, mulData.Name, mulData.BotID)
	}
}
