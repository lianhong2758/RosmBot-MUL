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
		raw := new(RawPrivateMessage)
		err := json.Unmarshal(playload.D, raw)
		if err != nil {
			log.Errorln("[qq-err]", err)
			return
		}
		log.Infof("[qq] [↓]私聊消息%s:%s", raw.Author.UserOpenid, raw.Content)
		ctx := &rosm.Ctx{
			Bot:     c,
			BotType: "qq_group",
			Message: raw,
			Being: &rosm.Being{
				GroupID: raw.Author.UserOpenid,
				User: &rosm.UserData{
					ID: raw.Author.UserOpenid,
				},
				MsgID: raw.ID,
			},
			State: H{"type": playload.T, "id": raw.ID},
		}
		word := raw.Content
		//判断@
		ctx.Being.IsAtMe = true
		ctx.Being.RawWord = strings.TrimSpace(word)
		log.Debugf("[debug]关键词切割结果: %s", ctx.Being.RawWord)
		ctx.RunWord()
		//群聊
	case "GROUP_AT_MESSAGE_CREATE":
		raw := new(RawGroupMessage)
		err := json.Unmarshal(playload.D, raw)
		if err != nil {
			log.Errorln("[qq-err]", err)
			return
		}
		log.Infof("[qq] [↓]群聊消息[%s]%s:%s", raw.GroupID, raw.Author.ID, raw.Content)
		ctx := &rosm.Ctx{
			Bot:     c,
			BotType: "qq_group",
			Message: raw,
			Being: &rosm.Being{
				GroupID: raw.GroupID,
				User: &rosm.UserData{
					ID:   raw.Author.ID,
					Name: raw.Author.ID[len(raw.Author.ID)-8:],
				},
				MsgID: raw.ID,
			},
			State: H{"type": playload.T, "id": raw.ID},
		}
		word := raw.Content
		//判断@
		ctx.Being.IsAtMe = true
		ctx.Being.RawWord = strings.TrimSpace(word)
		log.Debugf("[debug]关键词切割结果: %s", ctx.Being.RawWord)
		ctx.RunWord()
		//频道私聊
		//文字子频道@机器人
		//文字子频道全量消息（私域）
	case "DIRECT_MESSAGE_CREATE", "AT_MESSAGE_CREATE", "MESSAGE_CREATE":
		raw := new(RawGuildMessage)
		err := json.Unmarshal(playload.D, raw)
		if err != nil {
			log.Errorln("[info]", err)
			return
		}
		log.Infof("[qq] [↓]频道消息[%s]%s:%s", raw.GuildID, raw.Author.Username, raw.Content)
		at := []string{}
		for _, v := range raw.Mentions {
			at = append(at, v.ID)
		}
		ctx := &rosm.Ctx{
			Bot:     c,
			BotType: "qq_gulid",
			Message: raw,
			Being: &rosm.Being{
				GroupID: raw.ChannelID,
				GuildID: raw.GuildID,
				ATList:  at,
				User: &rosm.UserData{
					Name: raw.Author.Username,
					ID:   raw.Author.ID,
				},
				MsgID: raw.ID,
			},
			State: H{"type": playload.T, "id": raw.ID},
		}
		word := raw.Content
		//判断@
		if strings.Contains(raw.Content, "<@!"+c.Ready.User.ID+">") {
			ctx.Being.IsAtMe = true
			word = strings.TrimSpace(strings.Replace(word, "<@!"+c.Ready.User.ID+">", "", 1))
		}
		ctx.Being.RawWord = word
		log.Debugf("[debug]关键词切割结果: %s", word)
		ctx.RunWord()
	}
}
