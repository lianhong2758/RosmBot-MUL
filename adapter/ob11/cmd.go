package ob11

import (
	"encoding/json"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// DeleteMessage 撤回消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF
//
//nolint:interfacer
func DeleteMessage(ctx *rosm.Ctx, messageID string) {
	CallAction(ctx, "delete_msg", Params{
		"message_id": messageID,
	})
}

// GetMessage 获取消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_msg-%E8%8E%B7%E5%8F%96%E6%B6%88%E6%81%AF
//
//nolint:interfacer
func GetMessage(ctx *rosm.Ctx, messageID string) MessageContent {
	rsp := CallAction(ctx, "get_msg", Params{
		"message_id": messageID,
	}).Data
	m := MessageContent{
		Elements:    ParseMessage(tool.StringToBytes(rsp.Get("message").Raw)),
		MessageId:   rsp.Get("message_id").String(),
		MessageType: rsp.Get("message_type").String(),
		Sender:      &User{},
	}
	err := json.Unmarshal(tool.StringToBytes(rsp.Get("sender").Raw), m.Sender)
	if err != nil {
		return MessageContent{}
	}
	return m
}

// GetForwardMessage 获取合并转发消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_forward_msg-%E8%8E%B7%E5%8F%96%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E6%B6%88%E6%81%AF
func GetForwardMessage(ctx *rosm.Ctx, id string) gjson.Result {
	rsp := CallAction(ctx, "get_forward_msg", Params{
		"id": id,
	}).Data
	return rsp
}

// SendLike 发送好友赞
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_like-%E5%8F%91%E9%80%81%E5%A5%BD%E5%8F%8B%E8%B5%9E
func SendLike(ctx *rosm.Ctx, userID int64, times int) {
	CallAction(ctx, "send_like", Params{
		"user_id": userID,
		"times":   times,
	})
}

// SetGroupKick 群组踢人
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func SetGroupKick(ctx *rosm.Ctx, groupID, userID int64, rejectAddRequest bool) {
	CallAction(ctx, "set_group_kick", Params{
		"group_id":           groupID,
		"user_id":            userID,
		"reject_add_request": rejectAddRequest,
	})
}

// SetThisGroupKick 本群组踢人
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func SetThisGroupKick(ctx *rosm.Ctx, userID int64, rejectAddRequest bool) {
	SetGroupKick(ctx, tool.StringToInt64(ctx.Being.GroupID), userID, rejectAddRequest)
}

// SetGroupBan 群组单人禁言
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80
func SetGroupBan(ctx *rosm.Ctx, groupID, userID, duration int64) {
	CallAction(ctx, "set_group_ban", Params{
		"group_id": groupID,
		"user_id":  userID,
		"duration": duration,
	})
}

// SetThisGroupBan 本群组单人禁言
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80
func SetThisGroupBan(ctx *rosm.Ctx, userID, duration int64) {
	SetGroupBan(ctx, tool.StringToInt64(ctx.Being.GroupID), userID, duration)
}

// SetGroupWholeBan 群组全员禁言
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func SetGroupWholeBan(ctx *rosm.Ctx, groupID int64, enable bool) {
	CallAction(ctx, "set_group_whole_ban", Params{
		"group_id": groupID,
		"enable":   enable,
	})
}

// SetThisGroupWholeBan 本群组全员禁言
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func SetThisGroupWholeBan(ctx *rosm.Ctx, enable bool) {
	SetGroupWholeBan(ctx, tool.StringToInt64(ctx.Being.GroupID), enable)
}

// SetGroupAdmin 群组设置管理员
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func SetGroupAdmin(ctx *rosm.Ctx, groupID, userID int64, enable bool) {
	CallAction(ctx, "set_group_admin", Params{
		"group_id": groupID,
		"user_id":  userID,
		"enable":   enable,
	})
}

// SetThisGroupAdmin 本群组设置管理员
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func SetThisGroupAdmin(ctx *rosm.Ctx, userID int64, enable bool) {
	SetGroupAdmin(ctx, tool.StringToInt64(ctx.Being.GroupID), userID, enable)
}

// SetGroupAnonymous 群组匿名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func SetGroupAnonymous(ctx *rosm.Ctx, groupID int64, enable bool) {
	CallAction(ctx, "set_group_anonymous", Params{
		"group_id": groupID,
		"enable":   enable,
	})
}

// SetThisGroupAnonymous 群组匿名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func SetThisGroupAnonymous(ctx *rosm.Ctx, enable bool) {
	SetGroupAnonymous(ctx, tool.StringToInt64(ctx.Being.GroupID), enable)
}

// SetGroupCard 设置群名片（群备注）
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_card-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D%E7%89%87%E7%BE%A4%E5%A4%87%E6%B3%A8
func SetGroupCard(ctx *rosm.Ctx, groupID, userID int64, card string) {
	CallAction(ctx, "set_group_card", Params{
		"group_id": groupID,
		"user_id":  userID,
		"card":     card,
	})
}

// SetThisGroupCard 设置本群名片（群备注）
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_card-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D%E7%89%87%E7%BE%A4%E5%A4%87%E6%B3%A8
func SetThisGroupCard(ctx *rosm.Ctx, userID int64, card string) {
	SetGroupCard(ctx, tool.StringToInt64(ctx.Being.GroupID), userID, card)
}

// SetGroupName 设置群名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func SetGroupName(ctx *rosm.Ctx, groupID int64, groupName string) {
	CallAction(ctx, "set_group_name", Params{
		"group_id":   groupID,
		"group_name": groupName,
	})
}

// SetThisGroupName 设置本群名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func SetThisGroupName(ctx *rosm.Ctx, groupID int64, groupName string) {
	SetGroupName(ctx, tool.StringToInt64(ctx.Being.GroupID), groupName)
}

// SetGroupLeave 退出群组
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func SetGroupLeave(ctx *rosm.Ctx, groupID int64, isDismiss bool) {
	CallAction(ctx, "set_group_leave", Params{
		"group_id":   groupID,
		"is_dismiss": isDismiss,
	})
}

// SetThisGroupLeave 退出本群组
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func SetThisGroupLeave(ctx *rosm.Ctx, isDismiss bool) {
	SetGroupLeave(ctx, tool.StringToInt64(ctx.Being.GroupID), isDismiss)
}

// SetGroupSpecialTitle 设置群组专属头衔
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func SetGroupSpecialTitle(ctx *rosm.Ctx, groupID, userID int64, specialTitle string) {
	CallAction(ctx, "set_group_special_title", Params{
		"group_id":      groupID,
		"user_id":       userID,
		"special_title": specialTitle,
	})
}

// SetThisGroupSpecialTitle 设置本群组专属头衔
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func SetThisGroupSpecialTitle(ctx *rosm.Ctx, userID int64, specialTitle string) {
	SetGroupSpecialTitle(ctx, tool.StringToInt64(ctx.Being.GroupID), userID, specialTitle)
}

// SetFriendAddRequest 处理加好友请求
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_friend_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E5%A5%BD%E5%8F%8B%E8%AF%B7%E6%B1%82
func SetFriendAddRequest(ctx *rosm.Ctx, flag string, approve bool, remark string) {
	CallAction(ctx, "set_friend_add_request", Params{
		"flag":    flag,
		"approve": approve,
		"remark":  remark,
	})
}

// SetGroupAddRequest 处理加群请求／邀请
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E7%BE%A4%E8%AF%B7%E6%B1%82%E9%82%80%E8%AF%B7
func SetGroupAddRequest(ctx *rosm.Ctx, flag string, subType string, approve bool, reason string) {
	CallAction(ctx, "set_group_add_request", Params{
		"flag":     flag,
		"sub_type": subType,
		"approve":  approve,
		"reason":   reason,
	})
}

// GetLoginInfo 获取登录号信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_login_info-%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95%E5%8F%B7%E4%BF%A1%E6%81%AF
func GetLoginInfo(ctx *rosm.Ctx) gjson.Result {
	return CallAction(ctx, "get_login_info", Params{}).Data
}

// GetStrangerInfo 获取陌生人信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_stranger_info-%E8%8E%B7%E5%8F%96%E9%99%8C%E7%94%9F%E4%BA%BA%E4%BF%A1%E6%81%AF
func GetStrangerInfo(ctx *rosm.Ctx, userID int64, noCache bool) gjson.Result {
	return CallAction(ctx, "get_stranger_info", Params{
		"user_id":  userID,
		"no_cache": noCache,
	}).Data
}

// GetFriendList 获取好友列表
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_friend_list-%E8%8E%B7%E5%8F%96%E5%A5%BD%E5%8F%8B%E5%88%97%E8%A1%A8
func GetFriendList(ctx *rosm.Ctx) gjson.Result {
	return CallAction(ctx, "get_friend_list", Params{}).Data
}

// GetGroupInfo 获取群信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E4%BF%A1%E6%81%AF
func GetGroupInfo(ctx *rosm.Ctx, groupID int64, noCache bool) Group {
	rsp := CallAction(ctx, "get_group_info", Params{
		"group_id": groupID,
		"no_cache": noCache,
	}).Data
	group := Group{}
	_ = json.Unmarshal([]byte(rsp.Raw), &group)
	return group
}

// GetThisGroupInfo 获取本群信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E4%BF%A1%E6%81%AF
func GetThisGroupInfo(ctx *rosm.Ctx, noCache bool) Group {
	return GetGroupInfo(ctx, tool.StringToInt64(ctx.Being.GroupID), noCache)
}

// GetGroupList 获取群列表
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%88%97%E8%A1%A8
func GetGroupList(ctx *rosm.Ctx) gjson.Result {
	return CallAction(ctx, "get_group_list", Params{}).Data
}

// GetGroupMemberInfo 获取群成员信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF
func GetGroupMemberInfo(ctx *rosm.Ctx, groupID int64, userID int64, noCache bool) gjson.Result {
	return CallAction(ctx, "get_group_member_info", Params{
		"group_id": groupID,
		"user_id":  userID,
		"no_cache": noCache,
	}).Data
}

// GetThisGroupMemberInfo 获取本群成员信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF
func GetThisGroupMemberInfo(ctx *rosm.Ctx, userID int64, noCache bool) gjson.Result {
	return GetGroupMemberInfo(ctx, tool.StringToInt64(ctx.Being.GroupID), userID, noCache)
}

// GetGroupMemberList 获取群成员列表
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
func GetGroupMemberList(ctx *rosm.Ctx, groupID int64) gjson.Result {
	return CallAction(ctx, "get_group_member_list", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupMemberList 获取本群成员列表
func GetThisGroupMemberList(ctx *rosm.Ctx) gjson.Result {
	return GetGroupMemberList(ctx, tool.StringToInt64(ctx.Being.GroupID))
}

// GetGroupMemberListNoCache 无缓存获取群员列表
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
func GetGroupMemberListNoCache(ctx *rosm.Ctx, groupID int64) gjson.Result {
	return CallAction(ctx, "get_group_member_list", Params{
		"group_id": groupID,
		"no_cache": true,
	}).Data
}

// GetThisGroupMemberListNoCache 无缓存获取本群员列表
func GetThisGroupMemberListNoCache(ctx *rosm.Ctx) gjson.Result {
	return GetGroupMemberListNoCache(ctx, tool.StringToInt64(ctx.Being.GroupID))
}

// GetGroupHonorInfo 获取群荣誉信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func GetGroupHonorInfo(ctx *rosm.Ctx, groupID int64, hType string) gjson.Result {
	return CallAction(ctx, "get_group_honor_info", Params{
		"group_id": groupID,
		"type":     hType,
	}).Data
}

// GetThisGroupHonorInfo 获取本群荣誉信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func GetThisGroupHonorInfo(ctx *rosm.Ctx, hType string) gjson.Result {
	return GetGroupHonorInfo(ctx, tool.StringToInt64(ctx.Being.GroupID), hType)
}

// GetRecord 获取语音
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_record-%E8%8E%B7%E5%8F%96%E8%AF%AD%E9%9F%B3
func GetRecord(ctx *rosm.Ctx, file string, outFormat string) gjson.Result {
	return CallAction(ctx, "get_record", Params{
		"file":       file,
		"out_format": outFormat,
	}).Data
}

// GetImage 获取图片
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_image-%E8%8E%B7%E5%8F%96%E5%9B%BE%E7%89%87
func GetImage(ctx *rosm.Ctx, file string) gjson.Result {
	return CallAction(ctx, "get_image", Params{
		"file": file,
	}).Data
}

// GetVersionInfo 获取运行状态
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_status-%E8%8E%B7%E5%8F%96%E8%BF%90%E8%A1%8C%E7%8A%B6%E6%80%81
func GetVersionInfo(ctx *rosm.Ctx) gjson.Result {
	return CallAction(ctx, "get_version_info", Params{}).Data
}

// Expand API

// SetGroupPortrait 设置群头像
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%A4%B4%E5%83%8F
func SetGroupPortrait(ctx *rosm.Ctx, groupID int64, file string) {
	CallAction(ctx, "set_group_portrait", Params{
		"group_id": groupID,
		"file":     file,
	})
}

// SetThisGroupPortrait 设置本群头像
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%A4%B4%E5%83%8F
func SetThisGroupPortrait(ctx *rosm.Ctx, file string) {
	SetGroupPortrait(ctx, tool.StringToInt64(ctx.Being.GroupID), file)
}

// OCRImage 图片OCR
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E5%9B%BE%E7%89%87ocr
func OCRImage(ctx *rosm.Ctx, file string) gjson.Result {
	return CallAction(ctx, "ocr_image", Params{
		"image": file,
	}).Data
}

// GetGroupSystemMessage 获取群系统消息
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E7%B3%BB%E7%BB%9F%E6%B6%88%E6%81%AF
func GetGroupSystemMessage(ctx *rosm.Ctx) gjson.Result {
	return CallAction(ctx, "get_group_system_msg", Params{}).Data
}

// MarkMessageAsRead 标记消息已读
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E6%A0%87%E8%AE%B0%E6%B6%88%E6%81%AF%E5%B7%B2%E8%AF%BB
func MarkMessageAsRead(ctx *rosm.Ctx, messageID int64) APIResponse {
	return CallAction(ctx, "mark_msg_as_read", Params{
		"message_id": messageID,
	})
}

// MarkThisMessageAsRead 标记本消息已读
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E6%A0%87%E8%AE%B0%E6%B6%88%E6%81%AF%E5%B7%B2%E8%AF%BB
func MarkThisMessageAsRead(ctx *rosm.Ctx) APIResponse {
	return CallAction(ctx, "mark_msg_as_read", Params{
		"message_id": ctx.Being.MsgID[0],
	})
}

// GetOnlineClients 获取当前账号在线客户端列表
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E5%BD%93%E5%89%8D%E8%B4%A6%E5%8F%B7%E5%9C%A8%E7%BA%BF%E5%AE%A2%E6%88%B7%E7%AB%AF%E5%88%97%E8%A1%A8
func GetOnlineClients(ctx *rosm.Ctx, noCache bool) gjson.Result {
	return CallAction(ctx, "get_online_clients", Params{
		"no_cache": noCache,
	}).Data
}

// GetGroupAtAllRemain 获取群@全体成员剩余次数
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%85%A8%E4%BD%93%E6%88%90%E5%91%98%E5%89%A9%E4%BD%99%E6%AC%A1%E6%95%B0
func GetGroupAtAllRemain(ctx *rosm.Ctx, groupID int64) gjson.Result {
	return CallAction(ctx, "get_group_at_all_remain", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupAtAllRemain 获取本群@全体成员剩余次数
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%85%A8%E4%BD%93%E6%88%90%E5%91%98%E5%89%A9%E4%BD%99%E6%AC%A1%E6%95%B0
func GetThisGroupAtAllRemain(ctx *rosm.Ctx) gjson.Result {
	return GetGroupAtAllRemain(ctx, tool.StringToInt64(ctx.Being.GroupID))
}

// GetGroupMessageHistory 获取群消息历史记录
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%B6%88%E6%81%AF%E5%8E%86%E5%8F%B2%E8%AE%B0%E5%BD%95
//
//	messageID: 起始消息序号, 可通过 get_msg 获得
func GetGroupMessageHistory(ctx *rosm.Ctx, groupID, messageID int64) gjson.Result {
	return CallAction(ctx, "get_group_msg_history", Params{
		"group_id":    groupID,
		"message_seq": messageID,
	}).Data
}

// GettLatestGroupMessageHistory 获取最新群消息历史记录
func GetLatestGroupMessageHistory(ctx *rosm.Ctx, groupID int64) gjson.Result {
	return CallAction(ctx, "get_group_msg_history", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupMessageHistory 获取本群消息历史记录
//
//	messageID: 起始消息序号, 可通过 get_msg 获得
func GetThisGroupMessageHistory(ctx *rosm.Ctx, messageID int64) gjson.Result {
	return GetGroupMessageHistory(ctx, tool.StringToInt64(ctx.Being.GroupID), messageID)
}

// GettLatestThisGroupMessageHistory 获取最新本群消息历史记录
func GetLatestThisGroupMessageHistory(ctx *rosm.Ctx) gjson.Result {
	return GetLatestGroupMessageHistory(ctx, tool.StringToInt64(ctx.Being.GroupID))
}

// GetGroupEssenceMessageList 获取群精华消息列表
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%B2%BE%E5%8D%8E%E6%B6%88%E6%81%AF%E5%88%97%E8%A1%A8
func GetGroupEssenceMessageList(ctx *rosm.Ctx, groupID int64) gjson.Result {
	return CallAction(ctx, "get_essence_msg_list", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupEssenceMessageList 获取本群精华消息列表
func GetThisGroupEssenceMessageList(ctx *rosm.Ctx) gjson.Result {
	return GetGroupEssenceMessageList(ctx, tool.StringToInt64(ctx.Being.GroupID))
}

// SetGroupEssenceMessage 设置群精华消息
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%AE%BE%E7%BD%AE%E7%B2%BE%E5%8D%8E%E6%B6%88%E6%81%AF
func SetGroupEssenceMessage(ctx *rosm.Ctx, messageID int64) APIResponse {
	return CallAction(ctx, "set_essence_msg", Params{
		"message_id": messageID,
	})
}

// DeleteGroupEssenceMessage 移出群精华消息
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E7%A7%BB%E5%87%BA%E7%B2%BE%E5%8D%8E%E6%B6%88%E6%81%AF
func DeleteGroupEssenceMessage(ctx *rosm.Ctx, messageID int64) APIResponse {
	return CallAction(ctx, "delete_essence_msg", Params{
		"message_id": messageID,
	})
}

// GetWordSlices 获取中文分词
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E4%B8%AD%E6%96%87%E5%88%86%E8%AF%8D
func GetWordSlices(ctx *rosm.Ctx, content string) gjson.Result {
	return CallAction(ctx, ".get_word_slices", Params{
		"content": content,
	}).Data
}

// SendGuildChannelMessage 发送频道消息
func SendGuildChannelMessage(ctx *rosm.Ctx, guildID, channelID string, message interface{}) string {
	rsp := CallAction(ctx, "send_guild_channel_msg", Params{
		"guild_id":   guildID,
		"channel_id": channelID,
		"message":    message,
	}).Data.Get("message_id")
	if rsp.Exists() {
		log.Infof("[api] 发送频道消息(%v-%v): %v (id=%v)", guildID, channelID, formatMessage(message), rsp.Int())
		return rsp.String()
	}
	return "0" // 无法获取返回值
}

// UploadGroupFile 上传群文件
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E4%B8%8A%E4%BC%A0%E7%BE%A4%E6%96%87%E4%BB%B6
//
//	msg: FILE_NOT_FOUND FILE_SYSTEM_UPLOAD_API_ERROR ...
func UploadGroupFile(ctx *rosm.Ctx, groupID int64, file, name, folder string) APIResponse {
	return CallAction(ctx, "upload_group_file", Params{
		"group_id": groupID,
		"file":     file,
		"name":     name,
		"folder":   folder,
	})
}

// UploadThisGroupFile 上传本群文件
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E4%B8%8A%E4%BC%A0%E7%BE%A4%E6%96%87%E4%BB%B6
//
//	msg: FILE_NOT_FOUND FILE_SYSTEM_UPLOAD_API_ERROR ...
func UploadThisGroupFile(ctx *rosm.Ctx, file, name, folder string) APIResponse {
	return UploadGroupFile(ctx, tool.StringToInt64(ctx.Being.GroupID), file, name, folder)
} // SetMyAvatar 设置我的头像
// https://llonebot.github.io/zh-CN/develop/extends_api
func SetMyAvatar(ctx *rosm.Ctx, file string) APIResponse {
	return CallAction(ctx, "set_qq_avatar", Params{
		"file": file,
	})
}

// GetFile 下载收到的群文件或私聊文件
//
// https://llonebot.github.io/zh-CN/develop/extends_api
func GetFile(ctx *rosm.Ctx, fileID string) gjson.Result {
	return CallAction(ctx, "get_file", Params{
		"file_id": fileID,
	}).Data
}
