package qq

import (
	"encoding/json"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/message"
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
		word := strings.TrimSpace(raw.Content)
		log.Infof("[qq] [↓]私聊消息%s:%s", raw.Author.UserOpenid, word)
		ctx := &rosm.Ctx{
			Bot:     c,
			Message: message.Message{},
			Being: &rosm.Being{
				GroupID: raw.Author.UserOpenid,
				User: &rosm.UserData{
					ID: raw.Author.UserOpenid,
				},
				MsgID: raw.ID,
			},
			State: H{"type": playload.T, "id": raw.ID, "event": raw},
		}
		//判断@
		ctx.Being.IsAtMe = true
		ctx.Being.RawWord = word
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
		word := strings.TrimSpace(raw.Content)
		log.Infof("[qq] [↓]群聊消息[%s]%s:%s", raw.GroupID, raw.Author.ID, word)
		ctx := &rosm.Ctx{
			Bot:     c,
			Message: message.Message{message.Text(word)},
			Being: &rosm.Being{
				GroupID: raw.GroupID,
				User: &rosm.UserData{
					ID:   raw.Author.ID,
					Name: raw.Author.ID[len(raw.Author.ID)-8:],
				},
				MsgID: raw.ID,
			},
			State: H{"type": playload.T, "id": raw.ID, "event": raw},
		}

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
		atme := false
		//判断@
		// if strings.Contains(raw.Content, "<@!"+c.Ready.User.ID+">") {
		// 	atme = true
		// 	raw.Content = strings.TrimSpace(strings.Replace(raw.Content, "<@!"+c.Ready.User.ID+">", "", 1))
		// 	raw.Mentions = raw.Mentions[1:]
		// }
		word := raw.Content
		log.Infof("[qq] [↓]频道消息[%s]%s:%s", raw.GuildID, raw.Author.Username, raw.Content)
		at := []string{}
		msg := message.Message{}
		for _, v := range raw.Mentions {
			at = append(at, v.ID)
			// 必定存在
			// raw.Content = strings.TrimSpace(strings.Replace(raw.Content, "<@!"+c.Ready.User.ID+">", "", 1))
			ct := strings.Split(word, "<@!"+c.Ready.User.ID+">")
			msg, word = append(msg, message.Text(ct[0]), message.AT(v.ID)), ct[1]
		}
		msg = append(msg, message.Text(word))
		word = msg.ExtractPlainText()
		ctx := &rosm.Ctx{
			Bot:     c,
			Message: msg,
			Being: &rosm.Being{
				GroupID: raw.ChannelID,
				GuildID: raw.GuildID,
				ATList:  at,
				User: &rosm.UserData{
					Name: raw.Author.Username,
					ID:   raw.Author.ID,
				},
				RawWord: word,
				IsAtMe:  atme,
				MsgID:   raw.ID,
			},
			State: H{"type": playload.T, "id": raw.ID, "event": raw},
		}
		log.Debugf("[debug]关键词切割结果: %s", word)
		ctx.RunWord()
	}
}