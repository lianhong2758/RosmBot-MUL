package qqmsg

type H = map[string]any

// 键值对data结构
type KV = struct {
	Key   string `json:"key"`
	Value any    `json:"values"`
}

// 消息模板,content, embed, ark, image/file_image, markdown 至少需要有一个字段，否则无法下发消息
type Content struct {
	//通用
	Text      string      `json:"content,omitempty"`
	Ark       *ArkS       `json:"ark,omitempty"`               //ark
	Reference *ReferenceS `json:"message_reference,omitempty"` //消息引用
	Image     string      `json:"image,omitempty"`             //图片URL
	MsgID     string      `json:"msg_id,omitempty"`            //前置收到的消息ID，用于发送被动消息
	EventID   string      `json:"event_id,omitempty"`          //前置收到的事件ID，用于发送被动消息
	MarkDown  *MarkDownS  `json:"markdown,omitempty"`          //markDown
	//私聊/群聊
	Types    int        `json:"msg_type,omitempty"` //0 是文本，1 图文混排 ，2 是 md, 3 ark，4 embed，5 at @人或@all
	Keyboard *KeyboardS `json:"keyboard,omitempty"` //消息按钮
	//群聊
	Timestamp int64 `json:"timestamp,omitempty"` //unix 秒级时间戳
	//子频道
	Embed *Embed `json:"embed,omitempty"` //一种特殊的 ark
}

// 消息引用
type ReferenceS struct {
	ID        string `json:"message_id,omitempty"`
	NeedError bool   `json:"ignore_get_message_error,omitempty"`
}

// embed
type Embed struct{}

// 消息按钮
type KeyboardS struct{}

// MarkDown Object
type MarkDownS struct {
	Content string `json:"content,omitempty"`
	ID      string `json:"custom_template_id,omitempty"` //模板ID
	Params  []KV   `json:"params,omitempty"`
}

// ark
type ArkS struct {
	ID string `json:"template_id"` //模板ID
	KV KV     `json:"kv"`
}

type SendState struct {
	MsgID string `json:"id"`
	Time  int64  `json:"timestamp"` //时间戳
}
