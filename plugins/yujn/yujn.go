package yujn

// Package yujn 来源于 https://api.yujn.cn/ 的接口

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

const (
	yujnURL      = "https://api.yujn.cn"
	zzxjjURL     = yujnURL + "/api/zzxjj.php?type=video"
	baisisURL    = yujnURL + "/api/baisis.php?type=video"
	heisisURL    = yujnURL + "/api/heisis.php?type=video"
	xjjURL       = yujnURL + "/api/xjj.php?type=video"
	tianmeiURL   = yujnURL + "/api/tianmei.php?type=video"
	ndymURL      = yujnURL + "/api/ndym.php?type=video"
	sbklURL      = yujnURL + "/api/sbkl.php?type=video"
	nvgaoURL     = yujnURL + "/api/nvgao.php?type=video"
	luoliURL     = yujnURL + "/api/luoli.php?type=video"
	yuzuURL      = yujnURL + "/api/yuzu.php?type=video"
	xggURL       = yujnURL + "/api/xgg.php?type=video"
	rewuURL      = yujnURL + "/api/rewu.php?type=video"
	diaodaiURL   = yujnURL + "/api/diaodai.php?type=video"
	hanfuURL     = yujnURL + "/api/hanfu.php?type=video"
	jpyzURL      = yujnURL + "/api/jpmt.php?type=video"
	qingchunURL  = yujnURL + "/api/qingchun.php?type=video"
	ksbianzhuang = yujnURL + "/api/ksbianzhuang.php?type=video"
	dybianzhuang = yujnURL + "/api/bianzhuang.php?type=video"
	mengwaURL    = yujnURL + "/api/mengwa.php?type=video"
	chuandaURL   = yujnURL + "/api/chuanda.php?type=video"
)

var (
	engine = rosm.Register(&rosm.PluginData{
		DefaultOff: false,
		Name:       "遇见API",
		Help: "- /小姐姐视频 - /小姐姐视频2 - /黑丝视频 - /白丝视频\n" +
			"- /欲梦视频 - /甜妹视频 - /双倍快乐 - /纯情女高\n" +
			"- /萝莉视频 - /玉足视频 - /帅哥视频 - /热舞视频\n" +
			"- /吊带视频 - /汉服视频 - /极品狱卒 - /清纯视频\n" +
			"- /快手变装 - /抖音变装 - /萌娃视频 - /穿搭视频\n",
	})
	urlMap = map[string]string{
		"小姐姐视频":  zzxjjURL,
		"小姐姐视频2": xjjURL,
		"黑丝视频":   heisisURL,
		"白丝视频":   baisisURL,
		"欲梦视频":   ndymURL,
		"甜妹视频":   tianmeiURL,
		"双倍快乐":   sbklURL,
		"纯情女高":   nvgaoURL,
		"萝莉视频":   luoliURL,
		"玉足视频":   yuzuURL,
		"帅哥视频":   xggURL,
		"热舞视频":   rewuURL,
		"吊带视频":   diaodaiURL,
		"汉服视频":   hanfuURL,
		"极品狱卒":   jpyzURL,
		"清纯视频":   qingchunURL,
		"快手变装":   ksbianzhuang,
		"抖音变装":   dybianzhuang,
		"萌娃视频":   mengwaURL,
		"穿搭视频":   chuandaURL,
	}
)

func init() {
	// 这里是您的处理逻辑的switch case重构版本
	engine.OnWord("小姐姐视频", "小姐姐视频2", "黑丝视频", "白丝视频", "欲梦视频", "甜妹视频", "双倍快乐", "纯情女高", "萝莉视频", "玉足视频", "帅哥视频", "热舞视频", "吊带视频", "汉服视频", "极品狱卒", "清纯视频", "快手变装", "抖音变装", "萌娃视频", "穿搭视频").
		SetBlock(true).Limit(rosm.LimitByUser).Handle(func(ctx *rosm.Ctx) {
		videoType := ctx.Being.RawWord[1:]
		videoURL := urlMap[videoType]
		ctx.Send(message.Video(videoURL))
	})
}
