package qq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	log "github.com/sirupsen/logrus"
)

func makeHeard(ACCESS_TOKEN, APPID string) func(*http.Request) {
	return func(req *http.Request) {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", ACCESS_TOKEN)
		req.Header.Add("X-Union-Appid", APPID)
	}
}

// 通过链接获取data并写入结构体
func (c *Config) GetOpenAPI(shortUrl string, body, result any) (err error) {
	var data []byte
	if body != nil && body != "" && body != 0 {
		data, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	data, err = web.Web(clientConst, host+shortUrl, http.MethodGet, makeHeard(c.access, c.BotToken.AppId), bytes.NewReader(data))
	log.Debugln("[GetOpenAPI]", host+shortUrl, err)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

// 获取频道用户信息
func GetGuildUser(ctx *rosm.Ctx, uid string) (User *GuildUser, err error) {
	url := host + fmt.Sprintf(urlGuildGetUser, ctx.Being.RoomID2, uid)
	data, err := web.Web(clientConst, url, http.MethodGet, makeHeard(ctx.Bot.(*Config).access, ctx.Bot.(*Config).BotToken.AppId), nil)
	log.Debugln("[GetGuildUser][", url, "]", tool.BytesToString(data))
	if err != nil {
		return nil, err
	}
	User = new(GuildUser)
	err = json.Unmarshal(data, &User)
	return
}

// 上传文件获取file_info,媒体类型：1 图片，2 视频，3 语音，4 文件（暂不开放）
func UpFile(ctx *rosm.Ctx, url string, types int) (result *UpFileResult, err error) {
	var upurl string
	if ctx.Being.Def["type"].(string) == "GROUP_AT_MESSAGE_CREATE" {
		upurl = host + fmt.Sprintf(urlUPFileGroup, ctx.Being.RoomID)
	} else {
		upurl = host + fmt.Sprintf(urlUPFilePrivate, ctx.Being.User.ID)
	}
	data, _ := json.Marshal(H{"file_type": types, "url": url, "srv_send_msg": false, "file_data": nil})
	data, err = web.Web(clientConst, upurl, http.MethodPost, makeHeard(ctx.Bot.(*Config).access, ctx.Bot.(*Config).BotToken.AppId), bytes.NewReader(data))
	log.Debugln("[UpFile][", url, "]", tool.BytesToString(data))
	if err != nil {
		log.Infoln("[UpFile]上传文件失败,type:", types, "url:", url)
		return nil, err
	}
	result = new(UpFileResult)
	err = json.Unmarshal(data, &result)
	log.Info("[UpFile]", result.UUID)
	return
}

func NewDms(ctx *rosm.Ctx, userID, guildID string) (guild_id, channel_id string, err error) {
	data, _ := json.Marshal(H{"recipient_id": userID, "source_guild_id": guildID})
	data, err = web.Web(clientConst, host+urlDMS, http.MethodPost, makeHeard(ctx.Bot.(*Config).access, ctx.Bot.(*Config).BotToken.AppId), bytes.NewReader(data))
	log.Debugln("[NewDms][", urlDMS, "]", tool.BytesToString(data))
	if err != nil {
		log.Infoln("[NewDms]ERROR: ", err)
		return "", "", err
	}
	result := H{}
	err = json.Unmarshal(data, &result)
	log.Info("[NewDms]", result["guild_id"], result["channel_id"])
	return result["guild_id"].(string), result["channel_id"].(string), err

}

// 子频道撤回消息 hide需要false | true
func ChannelsDeleteMessage(ctx *rosm.Ctx, MessageID string, hide string) error {
	urlDeleteMessage := fmt.Sprintf(urlDeleteMessage, ctx.Being.RoomID, MessageID, hide)
	data, err := web.Web(clientConst, host+urlDeleteMessage, http.MethodDelete, makeHeard(ctx.Bot.(*Config).access, ctx.Bot.(*Config).BotToken.AppId), nil)
	log.Debugln("[DeleteMessage][", host+urlDeleteMessage, "]", tool.BytesToString(data))
	if err != nil {
		log.Infoln("[DeleteMessage]ERROR: ", err)
		return err
	}
	return nil
}

// 频道私信撤回消息 hide需要false | true
func DmsDeleteMessage(ctx *rosm.Ctx, GuildID, MessageID string, hide string) error {
	urlDeleteMessage := fmt.Sprintf(urlDmsDeleteMessage, GuildID, MessageID, hide)
	data, err := web.Web(clientConst, host+urlDeleteMessage, http.MethodDelete, makeHeard(ctx.Bot.(*Config).access, ctx.Bot.(*Config).BotToken.AppId), nil)
	log.Debugln("[DMSDeleteMessage][", host+urlDeleteMessage, "]", tool.BytesToString(data))
	if err != nil {
		log.Infoln("[DMSDeleteMessage]ERROR: ", err)
		return err
	}
	return nil
}

// 单聊
func UserDeleteMessage(ctx *rosm.Ctx, MessageID, openid string) error {
	urlDeleteMessage := fmt.Sprintf(urlDeleteUserMessage, openid, MessageID)
	data, err := web.Web(clientConst, host+urlDeleteMessage, http.MethodDelete, makeHeard(ctx.Bot.(*Config).access, ctx.Bot.(*Config).BotToken.AppId), nil)
	log.Debugln("[UserDeleteMessage][", host+urlDeleteMessage, "]", tool.BytesToString(data))
	if err != nil {
		log.Infoln("[UserDeleteMessage]ERROR: ", err)
		return err
	}
	return nil
}

// 群聊
func GroupDeleteMessage(ctx *rosm.Ctx, MessageID, group_openid string) error {
	urlDeleteMessage := fmt.Sprintf(urlDeleteGroupMessage, group_openid, MessageID)
	data, err := web.Web(clientConst, host+urlDeleteMessage, http.MethodDelete, makeHeard(ctx.Bot.(*Config).access, ctx.Bot.(*Config).BotToken.AppId), nil)
	log.Debugln("[GroupDeleteMessage][", host+urlDeleteMessage, "]", tool.BytesToString(data))
	if err != nil {
		log.Infoln("[GroupDeleteMessage]ERROR: ", err)
		return err
	}
	return nil
}

type GuildUser struct {
	User struct {
		ID               string `json:"id"`
		Username         string `json:"username"`
		Avatar           string `json:"avatar"`
		Bot              bool   `json:"bot"`
		UnionOpenid      string `json:"union_openid"`
		UnionUserAccount string `json:"union_user_account"`
	} `json:"user"`
	Nick     string    `json:"nick"`
	Roles    []string  `json:"roles"`
	JoinedAt time.Time `json:"joined_at"`
}
type UpFileResult struct {
	UUID     string `json:"file_uuid"` //文件 ID
	FileINFO string `json:"file_info"` //文件信息，用于发消息接口的 media 字段使用
	TTL      int    `json:"ttl"`       //有效期，表示剩余多少秒到期，到期后 file_info 失效，当等于 0 时，表示可长期使用
}
