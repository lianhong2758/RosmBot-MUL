package qq

import (
	"encoding/json"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

func (c *Config) process(playload *WebsocketPayload) {
	switch playload.T {
	//私聊
	case "C2C_MESSAGE_CREATE":
		//群聊
	case "GROUP_AT_MESSAGE_CREATE":
		raw := new(RawGroupMessage)
		err := json.Unmarshal(playload.D, raw)
		if err != nil {
			log.Errorln("[info]", err)
			return
		}
		log.Infof("[info]接收消息[%s]%s:%s", raw.GroupID, raw.Author.ID, raw.Content)
		ctx := &rosm.CTX{
			Bot:     c,
			BotType: "qqgroup",
			Message: raw,
			Being: &rosm.Being{
				RoomID: raw.GroupID,
				User: &rosm.UserData{
					ID: raw.Author.ID,
				},
				MsgID: []string{raw.ID},
				Def:   H{"type": playload.T, "id": raw.ID},
			},
		}
		word := raw.Content
		//判断@
		ctx.Being.AtMe = true
		word = strings.TrimSpace(word)
		log.Debugf("[debug]关键词切割结果: %s", word)
		ctx.RunWord(word)
		//频道私聊
	case "DIRECT_MESSAGE_CREATE":
		//文字子频道@机器人
		//文字子频道全量消息（私域）
	case "AT_MESSAGE_CREATE", "MESSAGE_CREATE":
		raw := new(RawGuildMessage)
		err := json.Unmarshal(playload.D, raw)
		if err != nil {
			log.Errorln("[info]", err)
			return
		}
		log.Infof("[info]接收消息[%s]%s:%s", raw.GuildID, raw.Author.Username, raw.Content)
		at := []string{}
		for _, v := range raw.Mentions {
			at = append(at, v.ID)
		}
		ctx := &rosm.CTX{
			Bot:     c,
			BotType: "qq",
			Message: raw,
			Being: &rosm.Being{
				RoomID:  raw.ChannelID,
				RoomID2: raw.GuildID,
				ATList:  at,
				User: &rosm.UserData{
					Name:        raw.Author.Username,
					ID:          raw.Author.ID,
					PortraitURI: raw.Author.Avatar,
				},
				MsgID: []string{raw.ID},
				Def:   H{"type": playload.T, "id": raw.ID},
			},
		}
		word := raw.Content
		//判断@
		if strings.Contains(raw.Content, "<@!"+c.Ready.User.ID+">") {
			ctx.Being.AtMe = true
			word = strings.TrimSpace(strings.Replace(word, "<@!"+c.Ready.User.ID+">", "", 1))
		}
		log.Debugf("[debug]关键词切割结果: %s", word)
		ctx.RunWord(word)
	}
}
