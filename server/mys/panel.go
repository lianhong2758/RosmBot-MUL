package mys 

import (
	"encoding/json"
	"fmt"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

func NewPanel() *InfoContent {
	return &InfoContent{
		Panel: &PanelStr{
			SmallComponentGroupList: [][]Component{},
			MidComponentGroupList:   [][]Component{},
			BigComponentGroupList:   [][]Component{},
		},
	}
}

func (i *InfoContent) Byte() (info []byte) {
	info, _ = json.Marshal(i)
	return
}

// 小型组件
func (i *InfoContent) Small(nextLine bool, c *Component) {
	add(&i.Panel.SmallComponentGroupList, c, 3, nextLine)
}

// 中型组件
func (i *InfoContent) Mid(nextLine bool, c *Component) {
	add(&i.Panel.MidComponentGroupList, c, 2, nextLine)
}

// 大型组件
func (i *InfoContent) Big(nextLine bool, c *Component) {
	add(&i.Panel.BigComponentGroupList, c, 1, nextLine)
}

// 模板id
func (i *InfoContent) Template(id int) { i.Panel.TemplateID = id }

// Text
func (i *InfoContent) TextBuild(ctx *rosm.CTX, m ...message.MessageSegment) {
	cif, _ := MakeMsgContent(ctx, m...)
	i.Content = (*cif.(*H))["content"].(*Content)
}

// 简易快捷的构建,和上面的build二选一
func (i *InfoContent) Title(title ...any) {
	if i.Content == nil {
		i.Content = &Content{
			Text: fmt.Sprint(title...),
		}
	} else {
		i.Content.Text = fmt.Sprint(title...)
	}
}
func add(arr *[][]Component, c *Component, maxLen int, nextLine bool) {
	// 获取最后一个子数组
	groupIndex := len(*arr) - 1

	// 如果当前组已满或者没有组，则创建新的组
	if groupIndex < 0 || len((*arr)[groupIndex]) == maxLen || nextLine {
		*arr = append(*arr, []Component{*c})
	} else {
		// 向当前组添加组件
		(*arr)[groupIndex] = append((*arr)[groupIndex], *c)
	}
}
