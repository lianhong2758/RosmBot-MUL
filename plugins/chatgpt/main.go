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
)

type sessionKey struct {
	group string
	user  string
}

var (
	apiKey = ""
	cache  = ttl.NewCache[sessionKey, []chatMessage](time.Minute * 15)
	// preinstall = []chatMessage{{
	// 	Role:    "system",
	// 	Content: "", //这里写预设
	// }}
)
var config = struct {
	Apikey   string
	ProxyURL string
	Mode     string
}{Apikey: "", ProxyURL: "https://api.alioth.center/akasha-whisper/v1/", Mode: "gpt-4o-mini"}

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "chatgpt",
		Help: "- // [对话内容]\n" +
			"主人权限:\n" +
			"- 设置OpenAI apikey xxx\n" +
			"- (添加|删除)预设 x xx\n" +
			"- 设置预设 x\n" +
			"- 查看预设 xx",
		DataFolder: "chatgpt",
		Config:     &config,
	})

	//预设存储文件夹
	if file.IsNotExist(en.DataFolder + "preset") {
		_ = os.MkdirAll(en.DataFolder+"preset", 0755)
	}

	getPrese := func(presetName string) (prese string, err error) {
		if presetName == "" {
			return "", nil
		}
		if file.IsExist(en.DataFolder + "preset" + "/" + presetName + ".txt") {
			contentb, err := os.ReadFile(en.DataFolder + "preset" + "/" + presetName + ".txt")
			if err != nil {
				return "", err
			}
			return tool.BytesToString(contentb), nil
		}
		return "", nil
	}

	en.OnRex(`^(?:chatgpt|//)([\s\S]*)$`).Handle(func(ctx *rosm.Ctx) {
		var messages []chatMessage
		args := ctx.Being.ResultWord[1]
		key := sessionKey{
			group: tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID),
			user:  ctx.Being.User.ID,
		}
		if args == "reset" || args == "重置记忆" {
			cache.Delete(key)
			ctx.Send(message.Text("已清除上下文！"))
			return
		}
		//未获取到
		if len(cache.Get(key)) == 0 {
			presetName, _ := rosm.PluginDB.FindString(en.Name, tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID))
			prese, err := getPrese(presetName)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			messages =
				[]chatMessage{{
					Role:    "system",
					Content: prese, //这里写预设
				}, {
					Role:    "user",
					Content: args,
				}}
		} else {
			messages = append(cache.Get(key), chatMessage{
				Role:    "user",
				Content: args,
			})
		}
		resp, err := completions(messages, apiKey, config.Mode)
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

	en.OnRex(`^设置\s*OpenAI\s*apikey\s*(.*)$`).SetRule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		apiKey = ctx.Being.ResultWord[1]
		config.Apikey = apiKey
		en.SaveConfig()
		ctx.Send(message.Text("设置成功"))
	})
	en.OnRex(`^(删除|添加)预设\s*(\S+)\s+(.*)$`).SetRule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		modename := ctx.Being.ResultWord[2]
		content := ctx.Being.ResultWord[3]
		if ctx.Being.ResultWord[1] == "添加" {
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
			ctx.Send(message.Text("添加预设`", modename, "`成功"))
			return
		}
		err := os.Remove(en.DataFolder + "preset" + "/" + modename + ".txt")
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text("删除预设`", modename, "`成功"))
	})
	en.OnRex(`^设置预设\s*(\S+)`).SetRule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		presetName := ctx.Being.ResultWord[1]
		contentfile := en.DataFolder + "preset" + "/" + ctx.Being.ResultWord[1] + ".txt"
		if file.IsExist(contentfile) {
			//记录预设名
			rosm.PluginDB.InsertString(en.Name, tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID), presetName)
			ctx.Send(message.Text("设置预设`", presetName, "`成功"))
		} else {
			ctx.Send(message.Text("设置预设`", presetName, "`失败: 预设不存在"))
		}
	})
	en.OnRex(`^删除本群预设$`).SetRule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		_ = rosm.PluginDB.InsertString(en.Name, tool.MergePadString(ctx.Being.GroupID, ctx.Being.GuildID), "")
		ctx.Send(message.Text(message.Text("删除预设成功")))
	})
	en.OnRex(`^查看预设\s*(\S+)$`).SetRule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		if ctx.Being.ResultWord[1] == "列表" {
			lists := []string{}
			files, _ := os.ReadDir(en.DataFolder + "preset")
			for _, file := range files {
				lists = append(lists, file.Name())
			}
			if len(lists) == 0 {
				ctx.Send(message.Reply(), message.Text("当前没有任何预设"))
				return
			}
			ctx.Send(message.Text("当前所有预设:\n", strings.Join(lists, "\n")))
			return
		}
		prese, err := getPrese(ctx.Being.ResultWord[1])
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text("预设`", ctx.Being.ResultWord[1], "`:\n", prese))
	})
}
