package ob11

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// CallAction 调用 cqhttp API
func CallAction(ctx *rosm.Ctx, action string, params zero.Params) zero.APIResponse {
	req := zero.APIRequest{
		Action: action,
		Params: params,
	}
	rsp, err := ctx.Being.Def["caller"].(zero.APICaller).CallApi(req)
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

// DeleteMessage 撤回消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF
//
//nolint:interfacer
func DeleteMessage(ctx *rosm.Ctx, messageID string) {
	CallAction(ctx, "delete_msg", zero.Params{
		"message_id": messageID,
	})
}

var base64Reg = regexp.MustCompile(`"type":"image","data":\{"file":"base64://[\w/\+=]+`)

// formatMessage 格式化消息数组
//
//	仅用在 log 打印
func formatMessage(msg interface{}) string {
	switch m := msg.(type) {
	case string:
		return m
	case message.CQCoder:
		return m.CQCode()
	case fmt.Stringer:
		return m.String()
	default:
		s, _ := json.Marshal(msg)
		return helper.BytesToString(base64Reg.ReplaceAllFunc(s, func(b []byte) []byte {
			buf := bytes.NewBuffer([]byte(`"type":"image","data":{"file":"`))
			b = b[40:]
			b, err := base64.StdEncoding.DecodeString(helper.BytesToString(b))
			if err != nil {
				buf.WriteString(err.Error())
			} else {
				m := md5.Sum(b)
				_, _ = hex.NewEncoder(buf).Write(m[:])
				buf.WriteString(".image")
			}
			return buf.Bytes()
		}))
	}
}
