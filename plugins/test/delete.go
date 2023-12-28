package test

import (
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
	"github.com/sirupsen/logrus"
)

func init() {
	en := rosm.Register(rosm.NewRegist("踢出别野", "- @机器人 踢出别野 @everyone", ""))
	en.AddRex(`踢出别野(.*)`).SetBlock(true).MUL("mys").Rule(rosm.OnlyOverHost()).Handle(func(ctx *rosm.Ctx) {
		list := ctx.Being.ATList
		if len(list) != 2 {
			return
		}
		x := list[1]
		logrus.Infof("[delete]别野%v 删除用户%v ", ctx.Being.RoomID2, x)
		err := mys.DeleteUser(ctx, x)
		ctx.Send(message.Text(err))
	})
}
