package lc

import (
	"fmt"
	"math/rand/v2"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "lc",
		Help: "- 每日一题\n"+
		"- rlc e/m/h",
	})
	en.AddWord("每日一题", "/每日一题").Handle(func(ctx *rosm.Ctx) {
		err := GetCsrftoken()
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		today, err := GetTodayTopic()
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		onetopic, err := GetOneTopic(today.Get("data.todayRecord.0.question.titleSlug").String())
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text(fmt.Sprintf("ID:%d Title:%s 难度:%s\n题目:%s",
			onetopic.Get("data.question.questionFrontendId").Int(),
			onetopic.Get("data.question.translatedTitle").String(),
			onetopic.Get("data.question.difficulty").String(),
			ProcessContent(onetopic))))
	})
	en.AddRex(`^/?(?:rlc|随机lc)\s*([EeMmHh])?$`).Handle(func(ctx *rosm.Ctx) {
		/*
			大概估计题库有3500+
			其中会员题700+
			EASY 900+
			MEDIUM 1900+
			HARD 800+
		*/
		difficulty := ctx.Being.Rex[1]
		err := GetCsrftoken()
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		//先rand一个起始id
		var id int
		switch difficulty {
		case "":
			id = rand.IntN(3500)
		case "E", "e":
			id = rand.IntN(900)
			difficulty ="EASY"
		case "M", "m":
			id = rand.IntN(1900)
				difficulty ="MEDIUM"
		case "H", "h":
			id = rand.IntN(800)
				difficulty ="HARD"
		}
		r, err := GetTopicList(id, difficulty)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		var title string
		for _, value := range r.Get("data").Get("problemsetQuestionList").Get("questions").Array() {
			v := value.Map()
			// 判断是否是要付费
			if !v["paidOnly"].Bool() {
				title = v["titleSlug"].String()
				break
			}
		}
		onetopic, err := GetOneTopic(title)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text(fmt.Sprintf("ID:%d Title:%s 难度:%s\n题目:%s",
			onetopic.Get("data.question.questionFrontendId").Int(),
			onetopic.Get("data.question.translatedTitle").String(),
			onetopic.Get("data.question.difficulty").String(),
			ProcessContent(onetopic))))
	})
}
