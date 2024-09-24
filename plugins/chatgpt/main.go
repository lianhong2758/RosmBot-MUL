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

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name: "chatgpt",
		Help: "- // [对话内容]\n" +
			"主人权限:\n" +
			"- 设置OpenAI apikey xxx\n" +
			"- (添加|删除)预设 x xx\n" +
			"- 设置预设 x\n"+
			"- 查看预设 xx",
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
	//rosm里读取
	// //读取预设
	// presetfile := en.DataFolder + "preset.txt"
	// if file.IsExist(presetfile) {
	// 	presetb, err := os.ReadFile(presetfile)
	// 	if err != nil {
	// 		logrus.Warn("[chatgpt]读取预设名失败...")
	// 	} else {
	// 		presetName = string(presetb)
	// 	}
	// }
	// if presetName != "" && file.IsExist(en.DataFolder+"preset"+"/"+presetName+".txt") {
	// 	contentb, err := os.ReadFile(en.DataFolder + "preset" + "/" + presetName + ".txt")
	// 	if err != nil {
	// 		logrus.Warn("[chatgpt]读取预设失败...")
	// 	} else {
	// 		//设置当前预设
	// 		preinstall = []chatMessage{{
	// 			Role:    "system",
	// 			Content: string(contentb), //这里写预设
	// 		}}
	// 	}
	// } else {
	// 	logrus.Warn("[chatgpt]预设不存在...")
	// }
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
			presetName, _ := rosm.PluginDB.FindString(en.Name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2))
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
		modeid, _ := rosm.PluginDB.FindInt(en.Name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2))
		resp, err := completions(messages, apiKey, modelList[modeid])
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
	en.AddRex(`^(删除|添加)预设\s*(\S+)\s+(.*)$`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		modename := ctx.Being.Rex[2]
		content := ctx.Being.Rex[3]
		if ctx.Being.Rex[1] == "添加" {
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
	en.AddRex(`^设置预设\s*(\S+)`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		presetName := ctx.Being.Rex[1]
		contentfile := en.DataFolder + "preset" + "/" + ctx.Being.Rex[1] + ".txt"
		if file.IsExist(contentfile) {
			//记录预设名
			rosm.PluginDB.InsertString(en.Name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2), presetName)
			ctx.Send(message.Text("设置预设`", presetName, "`成功"))
		} else {
			ctx.Send(message.Text("设置预设`", presetName, "`失败: 预设不存在"))
		}
	})
	en.AddRex(`^删除本群预设$`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		_ = rosm.PluginDB.InsertString(en.Name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2), "")
		ctx.Send(message.Text(message.Text("删除预设成功")))
	})
	en.AddRex(`^查看预设\s*(\S+)$`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		if ctx.Being.Rex[1] == "列表" {
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
		prese, err := getPrese(ctx.Being.Rex[1])
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text("预设`", ctx.Being.Rex[1], "`:\n", prese))
	})
}
