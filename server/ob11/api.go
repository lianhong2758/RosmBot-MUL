package ob11

import (
	"errors"
	"strconv"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// CallAction 调用 cqhttp API
func CallAction(ctx *rosm.Ctx, action string, Params zero.Params) zero.APIResponse {
	req := zero.APIRequest{
		Action: action,
		Params: Params,
	}
	rsp, err := ctx.Bot.(*Config).Driver.(zero.APICaller).CallApi(req)
	if err != nil {
		log.Errorln("[ob11] [↑]调用", action, "时出现错误: ", err)
	}
	if err == nil && rsp.RetCode != 0 {
		log.Errorln("[ob11] [↑]调用", action, "时出现错误, 返回值:", rsp.RetCode, ", 信息:", rsp.Msg, "解释:", rsp.Wording)
	}
	return rsp
}

// SendGroupMessage 发送群消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_group_msg-%E5%8F%91%E9%80%81%E7%BE%A4%E6%B6%88%E6%81%AF
func SendGroupMessage(ctx *rosm.Ctx, groupID int64, message interface{}) int64 {
	rsp := CallAction(ctx, "send_group_msg", zero.Params{ // 调用并保存返回值
		"group_id": groupID,
		"message":  message,
	}).Data.Get("message_id")
	if rsp.Exists() {
		log.Infof("[ob11] [↑][群消息(%v)]: %v (id=%v)", groupID, formatMessage(message), rsp.Int())
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// SendPrivateMessage 发送私聊消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func SendPrivateMessage(ctx *rosm.Ctx, userID int64, message interface{}) int64 {
	rsp := CallAction(ctx, "send_private_msg", zero.Params{
		"user_id": userID,
		"message": message,
	}).Data.Get("message_id")
	if rsp.Exists() {
		log.Infof("[ob11] [↑][私聊消息(%v)]: %v (id=%v)", userID, formatMessage(message), rsp.Int())
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// CardOrNickName 从 uid 获取群名片，如果没有则获取昵称
func CardOrNickName(ctx *rosm.Ctx, uid int64) (name string) {
	name = GetGroupMemberInfo(ctx, tool.StringToInt64(ctx.Being.RoomID), uid, false).Get("card").String()
	if name == "" {
		name = GetStrangerInfo(ctx, uid, false).Get("nickname").String()
	}
	return
}

// SendGroupForwardMessage 发送合并转发(群)
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E5%9B%BE%E7%89%87ocr
func SendGroupForwardMessage(ctx *rosm.Ctx, groupID int64, message message.Message) gjson.Result {
	return CallAction(ctx, "send_group_forward_msg", zero.Params{
		"group_id": groupID,
		"messages": message,
	}).Data
}

// SendPrivateForwardMessage 发送合并转发(私聊)
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E5%9B%BE%E7%89%87ocr
func SendPrivateForwardMessage(ctx *rosm.Ctx, userID int64, message message.Message) gjson.Result {
	return CallAction(ctx, "send_private_forward_msg", zero.Params{
		"user_id":  userID,
		"messages": message,
	}).Data
}

// ForwardFriendSingleMessage 转发单条消息到好友
//
// https://llonebot.github.io/zh-CN/develop/extends_api
func ForwardFriendSingleMessage(ctx *rosm.Ctx, userID int64, messageID interface{}) zero.APIResponse {
	return CallAction(ctx, "forward_friend_single_msg", zero.Params{
		"user_id":    userID,
		"message_id": messageID,
	})
}

// ForwardGroupSingleMessage 转发单条消息到群
//
// https://llonebot.github.io/zh-CN/develop/extends_api
func ForwardGroupSingleMessage(ctx *rosm.Ctx, groupID int64, messageID interface{}) zero.APIResponse {
	return CallAction(ctx, "forward_group_single_msg", zero.Params{
		"group_id":   groupID,
		"message_id": messageID,
	})
}

// SetMessageEmojiLike 发送表情回应
//
// https://llonebot.github.io/zh-CN/develop/extends_api
//
// emoji_id 参考 https://bot.q.qq.com/wiki/develop/api-v2/openapi/emoji/model.html#EmojiType
func SetMessageEmojiLike(ctx *rosm.Ctx, messageID interface{}, emojiID rune) error {
	ret := CallAction(ctx, "set_msg_emoji_like", zero.Params{
		"message_id": messageID,
		"emoji_id":   strconv.Itoa(int(emojiID)),
	}).Data.Get("errMsg").Str
	if ret != "" {
		return errors.New(ret)
	}
	return nil
}
