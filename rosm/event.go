package rosm

// 事件类型的集合,用于 OnNoticeWithType

const (
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
