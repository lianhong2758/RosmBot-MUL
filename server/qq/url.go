package qq

const (
	host                = "https://api.sgroup.qq.com"
	urlAccessToken      = "https://bots.qq.com/app/getAppAccessToken"
	urlSendPrivate      = "/v2/users/%v/messages"  //私聊
	urlSendGroup        = "/v2/groups/%v/messages" //群聊
	urlSendGuild        = "/channels/%v/messages"  //子频道
	urlSendGuildPrivate = "/v2/dms/%v/messages"    //频道私聊
	urlSendFileGroup    = "/v2/groups/%v/files"    //群富文本
	urlSendFilePrivate  = "/v2/users/%v/files"     //私聊富文本
	urlGetway           = "/gateway"
	urlGetwayWss        = "/gateway/bot"
)

/*
var urlMap map[string]string = map[string]string{
	"C2C_MESSAGE_CREATE":      urlSendPrivate,      //私聊
	"GROUP_AT_MESSAGE_CREATE": urlSendGroup,        //群聊
	"AT_MESSAGE_CREATE":       urlSendGuild,        //子频道
	"MESSAGE_CREATE":          urlSendGuild,        //子频道全量消息
	"DIRECT_MESSAGE_CREATE":   urlSendGuildPrivate, //频道私聊
	"File":                    urlSendFileGroup,    //群富文本
	"file":                    urlSendFilePrivate,  //私聊富文本

}
*/
