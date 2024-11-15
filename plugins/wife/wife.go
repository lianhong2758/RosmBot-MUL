package myplugin

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/file"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

var cards = []string{}

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name:       "每日抽老婆",
		Help:       "- /抽wife",
		DataFolder: "wife",
	})
	{
		data, _ := os.ReadFile("data/wife/wife.json")
		_ = json.Unmarshal(data, &cards)
		log.Infoln("[wife]加载", len(cards), "位wife...")
	}
	en.OnWord("/抽wife").Handle(func(ctx *rosm.Ctx) {
		id, err := strconv.Atoi(ctx.Being.User.ID)
		if err != nil {
			var lastThree string = ctx.Being.User.ID
			if len(ctx.Being.User.ID) >= 3 {
				lastThree = ctx.Being.User.ID[len(ctx.Being.User.ID)-3:]
			}
			for _, char := range lastThree {
				id += int(char)
			}
		}
		card := cards[tool.RandSenderPerDayN(int64(id), len(cards))]
		path := file.BOTPATH + "/" + en.DataFolder + "pic/" + card
		card = strings.Split(card, ".")[0]
		ctx.Send(message.AT(ctx.Being.User.ID, ctx.Being.User.Name), message.Text("今天的二次元老婆是~【", card, "】哒"), message.Image("file://"+path))
	})
}
