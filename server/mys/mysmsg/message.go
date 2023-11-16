package mysmsg

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
)

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

func MYContent(content any) message.MessageSegment {
	return message.MessageSegment{
		Type: "my",
		Data: H{
			"my": content,
		},
	}
}
