package ob11

import (
	"encoding/json"
	"strconv"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

//nolint:revive
type (
	NoCtxGetMsg  func(int64) message.Message
	NoCtxSendMsg func(message.Message) rosm.H
)

// Forward 合并转发
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91-
func Forward(id string) message.MessageSegment {
	return message.MessageSegment{
		Type: "forward",
		Data: rosm.H{
			"id": id,
		},
	}
}

// Node 合并转发节点
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E8%8A%82%E7%82%B9-
func Node(id int64) message.MessageSegment {
	return message.MessageSegment{
		Type: "node",
		Data: rosm.H{
			"id": strconv.FormatInt(id, 10),
		},
	}
}

// CustomNode 自定义合并转发节点
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E8%87%AA%E5%AE%9A%E4%B9%89%E8%8A%82%E7%82%B9
func CustomNode(nickname string, userID int64, content interface{}) message.MessageSegment {
	var str string
	switch c := content.(type) {
	case string:
		str = c
	case message.Message:
		str = RosmToZeroMessage(c).String()
	default:
		b, _ := json.Marshal(content)
		str = tool.BytesToString(b)
	}
	return message.MessageSegment{
		Type: "node",
		Data: map[string]string{
			"uin":     strconv.FormatInt(userID, 10),
			"name":    nickname,
			"content": str,
		},
	}
}

// SendToSelf ...
func SendToSelf(ctx *rosm.Ctx) NoCtxSendMsg {
	return func(msg message.Message) rosm.H {
		ctx.Being.GroupID = "-" + ctx.Bot.Card().BotID
		return ctx.Send(msg...)
	}
}

// FakeSenderForwardNode ...
func FakeSenderForwardNode(ctx *rosm.Ctx, msgs ...message.MessageSegment) message.MessageSegment {
	return CustomNode(
		CardOrNickName(ctx, tool.StringToInt64(ctx.Being.User.ID)),
		tool.StringToInt64(ctx.Being.User.ID),
		msgs)
}

// SendFakeForwardToGroup ...
func SendFakeForwardToGroup(ctx *rosm.Ctx, msgs ...message.MessageSegment) NoCtxSendMsg {
	return func(msg message.Message) rosm.H {
		return rosm.H{"id": SendGroupForwardMessage(ctx, tool.StringToInt64(ctx.Being.GroupID), message.Message{
			FakeSenderForwardNode(ctx, msg...),
			FakeSenderForwardNode(ctx, msgs...),
		}).Get("message_id").String()}
	}
}
