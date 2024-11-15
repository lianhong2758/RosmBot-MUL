package fhl

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	"github.com/sirupsen/logrus"
)

const host = "http://106.54.63.95:8080"

var helpMap = map[string]string{
	"A": "游戏类型: 梦笔生花\n规则:题目为单字或二字词，玩家轮流说出带有该字或词的诗句。",
	"B": "题目为一句诗句（可选择 5—9 的字数），其中的字按顺序依次作为关键字，玩家轮流各自说出包含当前关键字的诗句。",
	"C": "题目为一组固定字词与一组可消去字词（可选择“1 词 + 10 词”或“3 词 + 16 词”），玩家轮流从两组字词中各选择一个，说出同时含有两者的诗句。每个消去词只能被选择一次。",
	"D": "题目为两组字词（可选择每组 5—10 词）。玩家轮流从两组字词中各选择一个，说出同时含有两者的诗句。所有词都只能被选择一次。",
}

func init() { // 插件主体
	en := rosm.Register(&rosm.PluginData{
		Name: "飞花令",
		Help: "- /梦笔生花\n" +
			"- /走马观花[5-9]\n" +
			"- /天女散花[13]\n" +
			"- /雾里看花[5-10]",
	})
	en.OnWord("梦笔生花").Handle(casefunc("A"))
	en.OnRex(`^/走马观花\s*(\d)`).Handle(casefunc("B"))
	en.OnRex(`^/天女散花\s*(\d)`).Handle(casefunc("C"))
	en.OnRex(`^/雾里看花\s*(\d)`).Handle(casefunc("D"))
}

type TopicResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ModType       string `json:"modtype"`
		Size          int    `json:"size"`
		ID            string `json:"id"`
		SubjectString string `json:"subjectstring"`
	} `json:"data"`
}

type Topic struct {
	ModType string `json:"modtype"`
	Size    int    `json:"size"`
	ID      string `json:"id"`
}

type Answer struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type AnswerResq struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		//题目
		ModType       string `json:"modtype"`
		Size          int    `json:"size"`
		ID            string `json:"id"`
		SubjectString string `json:"subjectstring"`

		Text        string   `json:"text"` //ans
		Update      string   `json:"update"`
		HistoryText []string `json:"history"`

		NextOne int    `json:"user"`   //下一个该谁回答,默认0开始,也可以根据history%2推算
		Reason  string `json:"reason"` //游戏消息
	} `json:"data"`
}

func answer(ctx *rosm.Ctx, ar *AnswerResq) error {
	data, err := json.Marshal(Answer{
		ID:   "RosmBot" + ctx.Being.GroupID + ctx.Being.GuildID,
		Text: ctx.Being.RawWord,
	})
	if err != nil {

		return err
	}
	data, err = web.Web(web.NewDefaultClient(), host+"/answer", http.MethodPost, func(r *http.Request) {}, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, ar)
}

func gettopic(ctx *rosm.Ctx, types string, size int, r *TopicResp) error {
	data, err := json.Marshal(Topic{
		ModType: types,
		Size:    size,
		ID:      "RosmBot" + ctx.Being.GroupID + ctx.Being.GuildID,
	})
	if err != nil {
		return err
	}
	data, err = web.Web(web.NewDefaultClient(), host+"/gettopic", http.MethodPost, func(r *http.Request) {}, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, r)
}

func casefunc(types string) func(ctx *rosm.Ctx) {
	return func(ctx *rosm.Ctx) {
		r := new(TopicResp)
		size := 0
		if types != "A" {
			size, _ = strconv.Atoi(ctx.Being.ResultWord[1])
		}
		err := gettopic(ctx, types, size, r)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		if r.Code != 200 {
			ctx.Send(message.Text("ERROR: ", r.Message))
			return
		}
		ctx.Send(message.Text(helpMap[types], "\n题目如下:", r.Data.SubjectString))
		next, close := ctx.GetNext(rosm.AllMessage, false, rosm.OnlyAtMe(), rosm.OnlyTheRoom(ctx.Being.GroupID, ctx.Being.GuildID))
		defer close()
		ar := new(AnswerResq)
		for {
			select {
			case <-time.After(time.Second * 900):
				ctx.Send(message.Text("时间太久了,不玩了喵..."))
				return
			case ctx2 := <-next:
				if ctx2.Being.RawWord == "/结束游戏" {
					ctx.Send(message.Text("飞花令游戏结束..."))
					return
				}
				if err = answer(ctx2, ar); err != nil {
					ctx.Send(message.Text("ERROR: ", err))
					return
				}
				logrus.Debug("fhl result:", *ar)
				switch ar.Code {
				case 202:
					//游戏结束
					ctx.Send(message.Text(ar.Data.Reason), message.AT(ctx2.Being.User.ID, ctx2.Being.User.Name), message.Text("获胜。\nhistory:", strings.Join(ar.Data.HistoryText, "\n")))
					return
				case 200:
					tup := ""
					if ar.Data.Update != "" {
						tup = "\nUpdate: " + ar.Data.Update
					}
					ctx.Send(message.Text("题目:", ar.Data.SubjectString, "\n历史:", strings.Join(ar.Data.HistoryText, "\n"), tup))
				case 201:
					//不切题
					ctx.Send(message.Text("不切题之:", ar.Data.Reason))
				case 203:
					ctx.Send(message.Text(ar.Message))
					return
				default:
					ctx.Send(message.Text("ERROR: ", ar.Message))
					return
				}
			}
		}
	}
}
