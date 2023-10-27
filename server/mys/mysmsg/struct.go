package mysmsg

import (
	. "github.com/lianhong2758/RosmBot-MUL/message"
)

type Entities struct {
	Entity H   `json:"entity,omitempty"`
	Length int `json:"length,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// 消息模板
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
