package mysmsg

import (
	"fmt"

	"github.com/lianhong2758/RosmBot-MUL/message"
)

// 加粗
func BoldText(text ...any) message.MessageSegment {
	return message.MessageSegment{
		Type: "bold",
		Data: H{"text": fmt.Sprint(text...)},
	}
}

// 斜体
func ItalicText(text ...any) message.MessageSegment {
	return message.MessageSegment{
		Type: "italic",
		Data: H{"text": fmt.Sprint(text...)},
	}
}

// 删除线
func DeleteText(text ...any) message.MessageSegment {
	return message.MessageSegment{
		Type: "strikethrough",
		Data: H{"text": fmt.Sprint(text...)},
	}
}

// 下划线
func UnderlineText(text ...any) message.MessageSegment {
	return message.MessageSegment{
		Type: "underline",
		Data: H{"text": fmt.Sprint(text...)},
	}
}

// goto the room
func RoomLink(VillaID, RoomID string, RoomName string) message.MessageSegment {
	return message.MessageSegment{
		Type: "villa_room_link",
		Data: H{
			"text":  RoomName,
			"villa": VillaID,
			"room":  RoomID,
		},
	}
}

// 特殊结构
// 下标文字
func Badge(str BadgeStr) message.MessageSegment {
	return message.MessageSegment{
		Type: "badge",
		Data: H{
			"badge": str,
		},
	}
}

// 预览组件
func Preview(str PreviewStr) message.MessageSegment {
	return message.MessageSegment{
		Type: "view",
		Data: H{
			"view": str,
		},
	}
}

// 帖子,只能单独使用
func Post(postid string) message.MessageSegment {
	return message.MessageSegment{
		Type: "post",
		Data: H{
			"id": postid,
		},
	}
}
