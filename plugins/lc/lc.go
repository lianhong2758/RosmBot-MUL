package lc

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "每日一题",
		Help: "- 每日一题",
	})
	en.AddRex(`^/?每日一题\s*(\d)?$`).Handle(func(ctx *rosm.Ctx) {
		//id,_:=strconv.Atoi(ctx.Being.Rex[1])
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
		onetopic, err := GetOneTopic(today)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		content := onetopic.Get("data.question.translatedContent").String()
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
		var text bytes.Buffer
		doc.Contents().Each(func(i int, s *goquery.Selection) {
			text.WriteString(s.Text())
		})
		sptext := strings.Split(text.String(), "\n")
		text.Reset()
		for _, v := range sptext {
			if t := strings.TrimSpace(v); t != "" {
				text.WriteString(t)
				text.WriteByte('\n')
			}
		}
		// 输出处理后的文本内容
		ctx.Send(message.Text(fmt.Sprintf("ID:%d Title:%s 难度:%s\n题目:%s",
			onetopic.Get("data.question.questionFrontendId").Int(),
			onetopic.Get("data.question.translatedTitle").String(),
			onetopic.Get("data.question.difficulty").String(),
			text.String())))
	})
}
