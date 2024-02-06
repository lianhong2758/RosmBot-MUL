package qq

const (
	host                = "https://api.sgroup.qq.com"
	urlAccessToken      = "https://bots.qq.com/app/getAppAccessToken"
	urlSendPrivate      = "/v2/users/%v/messages"  //私聊
	urlSendGroup        = "/v2/groups/%v/messages" //群聊
	urlSendGuild        = "/channels/%v/messages"  //子频道
	urlSendGuildPrivate = "/dms/%v/messages"       //频道私聊
	urlUPFileGroup      = "/v2/groups/%v/files"    //群富文本
	urlUPFilePrivate    = "/v2/users/%v/files"     //私聊富文本
	urlGetway           = "/gateway"
	urlGetwayWss        = "/gateway/bot"
	urlDMS              = "/users/@me/dms"
	//频道
	urlGuildGetUser  = "/guilds/%v/members/%v"               //获取成员详情
	urlDeleteMessage = "/channels/%v/messages/%v?hidetip=%v" //撤回消息 DELETE,参数三为隐藏灰条
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
