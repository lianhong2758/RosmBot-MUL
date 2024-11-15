package rosm

// 上下文结构
type Ctx struct {
	Bot     Boter
	BotType string
	Message any            //解析后的原始消息
	Being   *Being         //常用消息解析,需实现
	State   map[string]any //框架中途产生的内容,不固定,即时有效
	on      bool           //插件是否在本群开启
}

// 常用数据
type Being struct {
	GroupID    string    //房间号,群号
	GuildID    string    //如果有需要,存放房间号上级号码
	GroupName  string    //房间名称,存在上级则存放上级名称
	User       *UserData //触发事件者信息
	ATList     []string  //at的id列表
	MsgID      string    //用于reply,存放消息id
	IsAtMe     bool      //是否是at机器人触发的事件
	RawWord    string    //原始消息
	KeyWord    string    //匹配词
	ResultWord []string  //匹配结果
}

// 触发者信息
type UserData struct {
	Name string
	ID   string
}
