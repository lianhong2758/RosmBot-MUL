package chatgpt

import (
	"os"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/ttl"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/sirupsen/logrus"
)

type sessionKey struct {
	group string
	user  string
}

var (
	apiKey     = ""
	presetName = "" //预设名
	cache      = ttl.NewCache[sessionKey, []chatMessage](time.Minute * 15)
	preinstall = []chatMessage{{
		Role:    "system",
		Content: "", //这里写预设
	}}
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "chatgpt",
		Help: "- // [对话内容]\n" +
			"- 设置OpenAI apikey xxx\n" +
			"- 添加预设 x xx\n" +
			"- 设置预设 x",
		DataFolder: "chatgpt",
	})

	//预设存储文件夹
	if file.IsNotExist(en.DataFolder + "preset") {
		_ = os.MkdirAll(en.DataFolder+"preset", 0755)
	}

	apikeyfile := en.DataFolder + "apikey.txt"
	if file.IsExist(apikeyfile) {
		apikey, err := os.ReadFile(apikeyfile)
		if err != nil {
			panic(err)
		} else {
			apiKey = string(apikey)
		}
	}
	//读取预设
	presetfile := en.DataFolder + "preset.txt"
	if file.IsExist(presetfile) {
		presetb, err := os.ReadFile(presetfile)
		if err != nil {
			logrus.Warn("[chatgpt]读取预设名失败...")
		} else {
			presetName = string(presetb)
		}
	}
	if presetName != "" && file.IsExist(en.DataFolder+"preset"+"/"+presetName+".txt") {
		contentb, err := os.ReadFile(en.DataFolder + "preset" + "/" + presetName + ".txt")
		if err != nil {
			logrus.Warn("[chatgpt]读取预设失败...")
		} else {
			//设置当前预设
			preinstall = []chatMessage{{
				Role:    "system",
				Content: string(contentb), //这里写预设
			}}
		}
	} else {
		logrus.Warn("[chatgpt]预设不存在...")
	}

	en.AddRex(`^(?:chatgpt|//)([\s\S]*)$`).Handle(func(ctx *rosm.Ctx) {
		var messages []chatMessage
		args := ctx.Being.Rex[1]
		key := sessionKey{
			group: tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2),
			user:  ctx.Being.User.ID,
		}
		if args == "reset" || args == "重置记忆" {
			cache.Delete(key)
			ctx.Send(message.Text("已清除上下文！"))
			return
		}
		//未获取到
		if len(cache.Get(key)) == 0 {
			messages = append(preinstall, chatMessage{
				Role:    "user",
				Content: args,
			})
		}else{
			messages = append(cache.Get(key),chatMessage{
				Role:    "user",
				Content: args,
			})
		}
		resp, err := completions(messages, apiKey,modelList["chatgpt4om"])
		if err != nil {
			ctx.Send(message.Text("请求ChatGPT失败: ", err))
			return
		}
		reply := resp.Choices[0].Message
		reply.Content = strings.TrimSpace(reply.Content)
		messages = append(messages, reply)
		cache.Set(key, messages)
		ctx.Send(message.Reply(), message.Text(reply.Content, "\n本次消耗token: ", resp.Usage.PromptTokens, "+", resp.Usage.CompletionTokens, "=", resp.Usage.TotalTokens))
	})

	en.AddRex(`^设置\s*OpenAI\s*apikey\s*(.*)$`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		apiKey = ctx.Being.Rex[1]
		f, err := os.Create(apikeyfile)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		defer f.Close()
		_, err = f.WriteString(apiKey)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text("设置成功"))
	})
	en.AddRex(`^添加预设\s*(\S+)\s+(.*)$`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		modename := ctx.Being.Rex[1]
		content := ctx.Being.Rex[2]
		f, err := os.Create(en.DataFolder + "preset" + "/" + modename + ".txt")
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		defer f.Close()
		_, err = f.WriteString(content)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text("添加预设成功"))
	})
	en.AddRex(`^设置预设\s*(\S+)`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		contentfile := en.DataFolder + "preset" + "/" + ctx.Being.Rex[1] + ".txt"
		if file.IsExist(contentfile) {
			contentb, err := os.ReadFile(contentfile)
			if err != nil {
				ctx.Send(message.Text("读取预设", ctx.Being.Rex[1], "失败,ERROR: ", err))
				return
			}
			//记录预设名
			presetName = ctx.Being.Rex[1]
			f, err := os.Create(en.DataFolder + "preset.txt")
			if err != nil {
				ctx.Send(message.Text("ERROR1: ", err))
				return
			}
			defer f.Close()
			_, err = f.WriteString(presetName)
			if err != nil {
				ctx.Send(message.Text("ERROR2: ", err))
				return
			}
			//设置当前预设
			preinstall = []chatMessage{{
				Role:    "system",
				Content: string(contentb), //这里写预设
			}}
			ctx.Send(message.Text("设置预设`", presetName, "`成功"))
		}
	})
}
