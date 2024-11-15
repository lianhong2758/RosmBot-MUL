package myplugin

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"

	"github.com/FloatTech/floatbox/file"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

var (
	infoMap = make(map[string]cardInfo, 80)
	cardMap = make(map[string]card, 80)
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name:       "抽塔罗牌",
		Help:       "- /抽塔罗牌",
		DataFolder: "tarot",
	})
	{
		data, _ := os.ReadFile(en.DataFolder + "tarots.json")
		_ = json.Unmarshal(data, &cardMap)
		for _, card := range cardMap {
			infoMap[card.Name] = card.cardInfo
		}
		log.Infoln("[tarot]加载", len(cardMap), "张塔罗牌...")
	}
	en.OnWord("/抽塔罗牌").Handle(func(ctx *rosm.Ctx) {
		reasons := [...]string{"您抽到的是~\n", "锵锵锵，塔罗牌的预言是~\n", "诶，让我看看您抽到了~\n"}
		position := [...]string{"『正位』", "『逆位』"}
		i := rand.Intn(22)
		p := rand.Intn(2)
		card := cardMap[strconv.Itoa(i)]
		name := card.Name
		description := card.Description
		if p == 1 {
			description = card.ReverseDescription
			card.ImgURL = "D" + card.ImgURL
		}
		imgpath := file.BOTPATH + "/" + en.DataFolder + card.ImgURL
		ctx.Send(message.Reply(), message.Text(reasons[rand.Intn(len(reasons))], position[p], "的『", name, "』\n其释义为: ", description), message.Image("file://"+imgpath))
	})
}

type cardInfo struct {
	Description        string `json:"description"`
	ReverseDescription string `json:"reverseDescription"`
	ImgURL             string `json:"imgUrl"`
}
type card struct {
	Name     string `json:"name"`
	cardInfo `json:"info"`
}
