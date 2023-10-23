package rosm

// 上下文结构
type CTX struct {
	Bot     Boter
	BotType string
	Message any    //解析后的原始消息
	Being   *Being //常用消息解析,需实现
}

// 常用数据
type Being struct {
	RoomID   string
	RoomID2  string
	RoomName string
	User     *UserData
	ATList   any //at列表
	AtMe     bool
	Word     string
	Rex      []string
	Def      map[string]any //自定义存储的参数
}

// 触发者信息
type UserData struct {
	Name        string
	ID          string
	PortraitURI string //如果直接回调没有可以不写
}
