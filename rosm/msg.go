package rosm

// 上下文结构
type Ctx struct {
	Bot     Boter
	BotType string
	Message any    //解析后的原始消息
	Being   *Being //常用消息解析,需实现
}

// 常用数据
type Being struct {
	RoomID   string         //房间号,群号
	RoomID2  string         //如果有需要,存放房间号上级号码
	RoomName string         //房间名称,存在上级则存放上级名称
	User     *UserData      //触发事件者信息
	ATList   []string       //at的id列表
	MsgID    []string       //用于reply,存放消息id,reply的其他需要参数写在第二位
	AtMe     bool           //是否是at机器人触发的事件
	Word     string         //接收的用户发送的信息,进行了首位的空格切割
	Rex      []string       //如果有正则匹配,这里存放匹配结果
	Def      map[string]any //自定义存储的参数
}

// 触发者信息
type UserData struct {
	Name        string
	ID          string
	PortraitURI string //如果直接回调没有可以不写
}
