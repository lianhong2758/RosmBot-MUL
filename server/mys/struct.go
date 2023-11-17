package mys

import ()

type H = map[string]any

// 回调的请求结构
type InfoSTR struct {
	Event struct {
		Robot struct {
			Template tem `json:"template"` // 机器人模板信息
			VillaID  int `json:"villa_id"` // 事件所属的大别野 id
		} `json:"robot"`
		Type       int      `json:"type"`
		ExtendData struct { // 包含事件的具体数据
			EventData EventData `json:"EventData"`
		} `json:"extend_data"`
		CreatedAt int64  `json:"created_at"`
		ID        string `json:"id"`
		SendAt    int    `json:"send_at"`
	} `json:"event"`
}

// 所有事件
type EventData struct {
	SendMessage       sendmessage       `json:"SendMessage"`
	JoinVilla         joinVilla         `json:"JoinVilla"`
	CreateRobot       changeRobot       `json:"CreateRobot"`
	DeleteRobot       changeRobot       `json:"DeleteRobot"`
	AddQuickEmoticon  addQuickEmoticon  `json:"AddQuickEmoticon"`
	AuditCallback     auditCallback     `json:"AuditCallback"`
	ClickMsgComponent clickMsgComponent `json:"ClickMsgComponent"`
}

// 机器人相关信息
type tem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Icon     string `json:"icon"`
	Commands []struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	} `json:"commands"`
}

// 用户@机器人发送消息
type sendmessage struct {
	Content    string          `json:"content"`
	FromUserID int             `json:"from_user_id"` // 发送者 id
	SendAt     int64           `json:"send_at"`      // 发送时间的时间戳
	RoomID     int             `json:"room_id"`      // 房间 id
	ObjectName int             `json:"object_name"`  // 目前只支持文本类型消息
	Nickname   string          `json:"nickname"`     // 用户昵称
	MsgUID     string          `json:"msg_uid"`      // 消息 id
	BotMsgID   string          `json:"bot_msg_id"`   // 如果被回复的消息从属于机器人，则该字段不为空字符串
	VillaID    int             `json:"villa_id"`     // 大别野 id
	QuoteMsg   MessageForQuote `json:"quote_msg"`
}
type MessageForQuote struct { // 回调消息引用消息的基础信息
	Content          string `json:"content"`            // 消息摘要，如果是文本消息，则返回消息的文本内容。如果是图片消息，则返回"[图片]"
	MsgUID           string `json:"msg_uid"`            // 消息 id
	BotMsgID         string `json:"bot_msg_id"`         // 如果消息从属于机器人，则该字段不为空字符串
	SendAt           int    `json:"send_at"`            // 发送时间的时间戳
	MsgType          string `json:"msg_type"`           // 消息类型，包括"文本"，"图片"，"帖子卡片"等
	FromUserID       int    `json:"from_user_id"`       // 发送者 id（整型）
	FromUserNickname string `json:"from_user_nickname"` // 发送者昵称
	FromUserIDStr    string `json:"from_user_id_str"`   // 发送者 id（字符串）可携带机器人发送者的id
}

// 有新用户加入大别野
type joinVilla struct {
	JoinUID          int    `json:"join_uid"`
	JoinUserNickname string `json:"join_user_nickname"`
	JoinAt           int64  `json:"join_at"`
}

// 大别野添加机器人实例,大别野删除机器人实例
type changeRobot struct {
	VillaID int `json:"villa_id"`
}

// 用户使用表情回复消息表态
type addQuickEmoticon struct {
	VillaID    int    `json:"villa_id"`
	RoomID     int    `json:"room_id"`
	UID        int    `json:"uid"`
	EmoticonID int    `json:"emoticon_id"`
	Emoticon   string `json:"emoticon"`
	MsgUID     string `json:"msg_uid"`
	BotMsgID   string `json:"bot_msg_id"`
	IsCancel   bool   `json:"is_cancel"`
}

// 审核结果回调
type auditCallback struct {
	AuditID     string `json:"audit_id"`
	BotTplID    string `json:"bot_tpl_id"`
	VillaID     int    `json:"villa_id"`
	RoomID      int    `json:"room_id"`
	UserID      int    `json:"user_id"`
	PassThrough string `json:"pass_through"`
	AuditResult int    `json:"audit_result"`
}

// 按钮回溯
type clickMsgComponent struct {
	VillaID     int    `json:"villa_id"`
	RoomID      int    `json:"room_id"`
	UID         int    `json:"uid"`
	MsgUID      string `json:"msg_uid"`
	BotMsgID    string `json:"bot_msg_id"`
	ComponentID string `json:"component_id"`
	TemplateID  string `json:"template_id"`
	Extra       string `json:"extra"`
}

// 接收的原始消息,解析
type MessageContent struct {
	Trace struct {
		VisualRoomVersion string `json:"visual_room_version"`
		AppVersion        string `json:"app_version"`
		ActionType        int    `json:"action_type"`
		BotMsgID          string `json:"bot_msg_id"`
		Client            string `json:"client"`
		Env               string `json:"env"`
		RongSdkVersion    string `json:"rong_sdk_version"`
	} `json:"trace"`
	Quote struct {
		QuotedMessageSendTime   int64  `json:"quoted_message_send_time"`
		OriginalMessageID       string `json:"original_message_id"`
		OriginalMessageSendTime int64  `json:"original_message_send_time"`
		QuotedMessageID         string `json:"quoted_message_id"`
	} `json:"quote"`
	MentionedInfo struct {
		MentionedContent string   `json:"mentionedContent"`
		UserIDList       []string `json:"userIdList"`
		Type             int      `json:"type"`
	} `json:"mentionedInfo"`
	User    user    `json:"user"`
	Content content `json:"content"`
}

type content struct {
	Images   []any `json:"images"`
	Entities []struct {
		Offset int `json:"offset"`
		Length int `json:"length"`
		Entity struct {
			Type  string `json:"type"`
			BotID string `json:"bot_id"`
		} `json:"entity"`
	} `json:"entities"`
	Text string `json:"text"`
}
type user struct {
	PortraitURI string `json:"portraitUri"`
	Extra       string `json:"extra"`
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	ID          string `json:"id"`
	Portrait    string `json:"portrait"`
}

// 消息发送回调
type SendState struct {
	ApiCode
	Data struct {
		BotMsgID string `json:"bot_msg_id"`
	} `json:"data"`
}

// api返回
type ApiCode struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
}
