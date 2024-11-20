package adapter

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

var AdapterList = map[string]map[string]rosm.Boter{}

// 添加bot到适配器列表
func AddNewBot(bot rosm.Boter) {
	mul := bot.Card().MUL
	if _, ok := AdapterList[mul.BotType]; !ok {
		AdapterList[mul.BotType] = map[string]rosm.Boter{}
	}
	AdapterList[mul.BotType][mul.BotID] = bot
	rosm.MULChan <- mul
}

func Delete(bot rosm.Boter) {
	mul := bot.Card().MUL
	delete(AdapterList[mul.BotType], mul.BotID)
}

// 获取对应bot消息,如果未找到返回nil
func GetBoter(botid string) rosm.Boter {
	//一级遍历所有平台
	for _, v := range AdapterList {
		if boter, ok := v[botid]; ok {
			return boter
		}
	}
	return nil
}

// 获取bot消息,如果未找到返回相应平台的随机boter,如果对应平台没有bot,则返回nil
func FindBoter(bottype, botid string) rosm.Boter {
	botlist, ok := AdapterList[bottype]
	if !ok {
		return nil
	}
	if boter, ok := botlist[botid]; ok {
		return boter
	}
	return GetRandBot(bottype)
}

// 获取对应平台的随机boter,如果没有则返回nil
func GetRandBot(bottype string) rosm.Boter {
	for _, v := range AdapterList[bottype] {
		return v
	}
	return nil
}

// 新建上下文
func NewCtx(bottype, botid, group, guild string) *rosm.Ctx {
	return &rosm.Ctx{
		Bot: FindBoter(bottype, botid),
		Being: &rosm.Being{
			GroupID: group,
			GuildID: guild,
		},
	}
}
func NewCtxWithPad(bottype, botid, padString string) *rosm.Ctx {
	gr, gu := tool.SplitPadString(padString)
	return &rosm.Ctx{
		Bot: FindBoter(bottype, botid),
		Being: &rosm.Being{
			GroupID: gr,
			GuildID: gu,
		},
	}
}
func NewCtxWithTypeAndPad(botid, typeAndTypeString string) *rosm.Ctx {
	types, pad := tool.SplitTypeAndPadString(typeAndTypeString)
	gr, gu := tool.SplitPadString(pad)
	return &rosm.Ctx{
		Bot: FindBoter(types, botid),
		Being: &rosm.Being{
			GroupID: gr,
			GuildID: gu,
		},
	}
}
