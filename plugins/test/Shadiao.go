package test

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

const (
	shadiaoURL = "https://api.shadiao.pro"
	chpURL     = shadiaoURL + "/chp"
	duURL      = shadiaoURL + "/du"
	pyqURL     = shadiaoURL + "/pyq"

	ganhaiURL = "https://api.lovelive.tools/api/SweetNothings/Web/1"
)

func init() {
	Engine := rosm.Register(&rosm.PluginData{
		Name: "沙雕插件",
		Help: "- /夸我 今天的你已经很棒了加油\n- /毒鸡汤 换个角度看问题世界从此不同\n- /朋友圈文案 向朋友们宣告你的新脑洞\n- /段子 感到不开心的话就笑一笑吧\n- /甘海 爱情需要一些甜蜜小短句",
	})
	Engine.OnWord("夸我").SetBlock(true).Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Text(GetShaDiaOApiText(chpURL)))
	})
	Engine.OnWord("毒鸡汤").SetBlock(true).Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Text(GetShaDiaOApiText(duURL)))
	})
	Engine.OnWord("朋友圈文案").SetBlock(true).Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Text(GetShaDiaOApiText(pyqURL)))
	})
	Engine.OnWord("甘海").SetBlock(true).Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Text(GetGanHaiApiText(ganhaiURL)))
	})
}

type DiaOStruct struct {
	Data struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"data"`
}

func GetShaDiaOApiText(Url string) string {
	get, err := http.Get(Url)
	if err != nil {
		return ""
	}
	all, err := io.ReadAll(get.Body)
	if err != nil {
		return ""
	}
	var diao DiaOStruct
	err = json.Unmarshal(all, &diao)
	if err != nil {
		return ""
	}
	return diao.Data.Text
}

type DuanZiStruct struct {
	Success string `json:"success"`
	Duanzi  string `json:"duanzi"`
	Qiafan  bool   `json:"qiafan"`
}

type GanHaiStruct struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	ReturnObj struct {
		Id           string `json:"id"`
		Content      string `json:"content"`
		LikeCount    int    `json:"likeCount"`
		DislikeCount int    `json:"dislikeCount"`
	} `json:"returnObj"`
}

func GetGanHaiApiText(Url string) string {
	get, err := http.Get(Url)
	if err != nil {
		return ""
	}
	all, err := io.ReadAll(get.Body)
	if err != nil {
		return ""
	}
	var diao GanHaiStruct
	err = json.Unmarshal(all, &diao)
	if err != nil {
		return ""
	}
	return diao.ReturnObj.Content
}
