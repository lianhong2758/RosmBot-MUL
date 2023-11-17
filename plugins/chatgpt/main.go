package chatgpt

import (
	"os"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/ttl"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

type sessionKey struct {
	group string
	user  string
}

var (
	apiKey = ""
	cache  = ttl.NewCache[sessionKey, []chatMessage](time.Minute * 15)
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name:       "chatgpt",
		Help:       "- // [对话内容]\n",
		DataFolder: "chatgpt",
	})
	apikeyfile := en.DataFolder + "apikey.txt"
	if file.IsExist(apikeyfile) {
		apikey, err := os.ReadFile(apikeyfile)
		if err != nil {
			panic(err)
		}
		apiKey = string(apikey)
	}
	en.AddRex(`^(?:chatgpt|//)([\s\S]*)$`).Handle(func(ctx *rosm.CTX) {
		var messages []chatMessage
		args := ctx.Being.Rex[1]
		key := sessionKey{
			group: ctx.Being.RoomID2,
			user:  ctx.Being.User.ID,
		}
		if args == "reset" || args == "重置记忆" {
			cache.Delete(key)
			ctx.Send(message.Text("已清除上下文！"))
			return
		}
		messages = cache.Get(key)
		messages = append(messages, chatMessage{
			Role:    "user",
			Content: args,
		})
		resp, err := completions(messages, apiKey)
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

	en.AddRex(`^设置\s*OpenAI\s*apikey\s*(.*)$`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.CTX) {
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
}
