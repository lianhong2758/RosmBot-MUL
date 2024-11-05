package rosm

// 事件类型的集合
//新的平台的事件类型都需要加入其中再实现
//无须在意顺序,但还是建议消息事件排最前,同一平台的事件放在一起
//可以把有相同性质的事件作为一个事件类型
//可以通过平台限制+事件类型,实现更具体的事件监听

type EventType string

const (
	//消息有关类型1-2
	AllMessage     = "AllMessage"     //全消息匹配
	SurplusMessage = "SurplusMessage" //无其他插件匹配
	//ob11标准 https://github.com/botuniverse/onebot-11/blob/master/event/message.md
	MemberJoin       = "group_increase" //入群
	MemberOut        = "group_decrease" //退群
	UpFile           = "group_upload"   //上传文件
	ChangeGroupAdmin = "group_admin"    //管理员变动
	GroupBan         = "group_ban"      //群禁言
	FriendAdd        = "friend_add"     //好友添加
	GroupRecall      = "group_recall"   //群消息撤回
	FriendRecall     = "friend_recall"  //好友消息撤回
	Poke             = "notify"         //群内戳一戳
	LuckyKing        = "lucky_king"     //群红包运气王
	Notify           = "notify"         //群成员荣誉变更
)
