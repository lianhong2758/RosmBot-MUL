package qq

import (
	"encoding/json"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

func (c *Config) process(playload *WebsocketPayload) {
	switch playload.T {
	//私聊
	case "C2C_MESSAGE_CREATE":
		//群聊
	case "GROUP_AT_MESSAGE_CREATE":
		//频道私聊
	case "DIRECT_MESSAGE_CREATE":
		//文字子频道@机器人
	case "AT_MESSAGE_CREATE":
		//文字子频道全量消息（私域）
	case "MESSAGE_CREATE":
		raw := new(RawMessage)
		err := json.Unmarshal(playload.D, raw)
		if err != nil {
			log.Errorln("[info]", err)
			return
		}
		log.Infof("[info]接收消息[%s]%s:%s", raw.GuildID, raw.Author.Username, raw.Content)
		ctx := &rosm.CTX{
			Bot:     c,
			BotType: "qq",
			Message: raw,
			Being: &rosm.Being{
				RoomID:  tool.Int64(raw.ChannelID),
				RoomID2: tool.Int64(raw.GuildID),
				ATList:  raw.Mentions,
				User: &rosm.UserData{
					Name: raw.Author.Username,
					ID:   tool.Int64(raw.Author.ID),
				},
				Def: H{"type": playload.T, "id": raw.ID},
			},
		}
		word := raw.Content
		//判断@
		if strings.Contains(raw.Content, "\u003c@!"+c.Ready.User.ID+"\u003e ") {
			ctx.Being.AtMe = true
			word = strings.Replace(word, "\u003c@!"+c.Ready.User.ID+"\u003e ", "", 1)
		}
		log.Debugf("[debug]关键词切割结果: %s", word)
		ctx.RunWord(word)
	}
}
