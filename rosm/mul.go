package rosm

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	log "github.com/sirupsen/logrus"
)

type H = message.H

// Boter bot主接口,由server实现,Run函数由server自己调用运行
type Boter interface {
	//运行,用于开启接收消息和调用插件
	Run()
	//Bot信息查询
	Card() *BotCard

	BotRuler
	BotCaller
}

// BotRuler rule的实现,如果不存在,适配器需要设置函数返回false
type BotRuler interface {
	//判断回复消息
	OnlyReply(*Ctx) bool
	//判断主人
	OnlyMaster(*Ctx) bool
	//判断群主等
	OnlyOverHost(*Ctx) bool
	//判断管理员
	OnlyOverAdministrator(*Ctx) bool
}

// 接口调用
type BotCaller interface {
	//发送消息
	BotSend(*Ctx, ...message.MessageSegment) H
	//特殊的发送格式
	BotSendCustom(*Ctx, any) H
	//头像获取
	GetPortraitURI(*Ctx) string
}

type BotCard struct {
	MUL
	Master []string `json:"master_id"`
}

// 进行一个通道注册,同于接收平台注册消息进行统计,虽然不是强制性需要,但还是建议实现一下
var MULChan = make(chan MUL)

type MUL struct {
	BotType string `json:"bot_type"`
	BotID   string `json:"bot_id,omitempty"`
}

func Listen() {
	plugindbinit()
	for mulData := range MULChan {
		log.Infof("[mul]新增注册,平台: %s,BotID: %s", mulData.BotType, mulData.BotID)
	}
}
