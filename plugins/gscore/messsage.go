package gscore

import (
	"encoding/json"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

func MakeSendCoreMessage(ctx *rosm.Ctx) []byte {
	MessageReport := MessageReceive{
		Bot_id:      ctx.BotType,
		Bot_self_id: ctx.Bot.Card().BotID,
		Msg_id:      ctx.Being.MsgID,
		User_type:   "group",
		Group_id:    tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2),
		User_id:     ctx.Being.User.ID,
		User_pm: func() int {
			if rosm.OnlyMaster()(ctx) {
				return 1
			} else if rosm.OnlyOverAdministrator()(ctx) {
				return 2
			}
			return 3
		}(),
		Content: []WriteMessage{
			{
				Type: "text",
				Data: ctx.Being.Word,
			},
		},
		Sender: Dictionary{
			Avater:   ctx.Bot.GetPortraitURI(ctx),
			Nickname: ctx.Being.User.Name,
		},
	}
	cache.Set(ctx.Being.MsgID, ctx)
	marshal, err := json.Marshal(MessageReport)
	if err != nil {
		return nil
	}
	return marshal
}
