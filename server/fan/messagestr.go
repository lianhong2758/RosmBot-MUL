package fan

import ()

type Entities struct {
	Entity H   `json:"entity,omitempty"`
	Length int `json:"length,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// 普通消息模板
type Content struct {
	//图片
	ImageStr
	//文本
	Text     string     `json:"text,omitempty"`
	Entities []Entities `json:"entities,omitempty"`
	Images   []ImageStr `json:"images,omitempty"`
	//链接预览
	Preview *PreviewStr `json:"preview_link,omitempty"`
	//下标
	Badge *BadgeStr `json:"badge,omitempty"`
}

type ImageStr struct {
	URL      string `json:"url,omitempty"`
	FileSize int    `json:"file_size,omitempty"`
	Size     *Size  `json:"size,omitempty"`
}
type Size struct {
	Height int `json:"height,omitempty"`
	Width  int `json:"width,omitempty"`
}

type MentionedInfoStr struct {
	Type       int      `json:"type"`
	UserIDList []string `json:"userIdList"`
}

// 预览图片的结构
type PreviewStr struct {
	URL        string `json:"url,omitempty"`              //跳转
	ImageURL   string `json:"image_url,omitempty"`        //大图
	IsIntLink  bool   `json:"is_internal_link,omitempty"` //link
	Title      string `json:"title,omitempty"`            //右上角大标题
	Content    string `json:"content,omitempty"`          //中间小字
	Icon       string `json:"icon_url,omitempty"`         //左下角小图
	SourceName string `json:"source_name,omitempty"`      //左下角图标名
}

// 左下标的结构
type BadgeStr struct {
	Icon string `json:"icon_url,omitempty"` //图标
	Text string `json:"text,omitempty"`     //文字
	URL  string `json:"url,omitempty"`      //跳转
}

// 组件模板
type PanelStr struct {
	TemplateID              int           `json:"template_id"`                //模板id，通过创建消息组件模板接口，可以提前将组件面板保存，使用 template_id来快捷发送消息
	SmallComponentGroupList [][]Component `json:"small_component_group_list"` //定义小型组件，即一行摆置3个组件，每个组件最多展示2个中文字符或4个英文字符
	MidComponentGroupList   [][]Component `json:"mid_component_group_list"`   //定义中型组件，即一行摆置2个组件，每个组件最多展示4个中文字符或8个英文字符
	BigComponentGroupList   [][]Component `json:"big_component_group_list"`   //定义大型组件，即一行摆置1个组件，每个组件最多展示10个中文字符或20个英文字符
}

// 组件单元
type Component struct {
	ID           string `json:"id"`            //组件id，由机器人自定义，不能为空字符串。面板内的id需要唯一
	Text         string `json:"text"`          //组件展示文本, 不能为空
	Type         int    `json:"type"`          //组件类型，目前支持 type=1 按钮组件，未来会扩展更多组件类型
	NeedCallback bool   `json:"need_callback"` //是否订阅该组件的回调事件
	Extra        string `json:"extra"`         //组件回调透传信息，由机器人自定义

	//按钮
	CType        int    `json:"c_type"`        //组件交互类型，包括：1回传型，2输入型，3跳转型
	InputContent string `json:"input_content"` //如果交互类型为输入型，则需要在该字段填充输入内容，不能为空
	Link         string `json:"link"`          //如果交互类型为跳转型，需要在该字段填充跳转链接，不能为空
	NeedToken    bool   `json:"need_token"`    //对于跳转链接来说，如果希望携带用户信息token，则need_token设置为true
}

// 发送消息的抽象结构
type InfoContent struct {
	Panel   *PanelStr `json:"panel,omitempty"`
	Content *Content  `json:"content"`
}
