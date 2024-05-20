package time

import (
	"time"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/send"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var (
	c          *cron.Cron
	entryIDMap = map[string]cron.EntryID{}
)

/*
快捷学习
* * * * *
- - - - -
| | | | |
| | | | ----- 星期 - 0-6(0表示星期天)
| | | ------- 月 - 1-12
| | --------- 日 - 1-31
| ----------- 小时 - 0-23
------------- 分钟 - 0-59
*/
// *: 匹配该字段所有值，如 0 * 1 1 *, 第三个字段为 * 表示（1 月 1 日）每小时。
// /: 表示范围增量，如 */12 * * * *  表示每 12 秒执行一次
// ,: 用来分隔同一组中的项目，如 * * 5,10,15 3,4 * * 表示每个三月或四月的 5， 10， 15 号（3.05， 3.10， 3.15， 4.05， 4.10，4.15）
// -: 表示范围，如 */5 * 10-12 * *  表示每天十点到十二点每五秒执行一次
// ?: 同 *
// 表示每隔多长时间时，你还可以使用预定义的 @every <duration> 如每隔十分钟就可以表示为 @every 10m
func init() {
	en := rosm.Register(&rosm.PluginData{
		Name:       "定时消息",
		Help:       "- /记录在* * * *的指令\n"+
		"- /删除指令xxx",
		DataFolder: "time",
	})
	en.AddRex(`^/记录在(.*)的指令`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		next, stop := ctx.GetNext(rosm.AllMessage, true, rosm.OnlyTheUser(ctx.Being.User.ID))
		ctx.Send(message.Text("发送想要记录的指令:"))
		var order string
		select {
		case <-time.After(time.Second * 60):
			ctx.Send(message.Text("时间太久了"))
			return
		case ctx2 := <-next:
			order = ctx2.Being.Word
		}
		//结束记录
		stop()
		// 允许往正在执行的 cron 中添加任务
		id, err := c.AddFunc(ctx.Being.Rex[1], func() { ctx.RunWord(order) })
		if err != nil {
			ctx.Send(message.Text("参数不合法,ERROR: ", err))
			return
		}
		entryIDMap[order] = id
		//开始存储数据库
		m := &mode{
			Key:     tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2) + order,
			String1: tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2),
			Types:   ctx.BotType,
			Word:    order,
			Time:    ctx.Being.Rex[1],
			UserID:  ctx.Being.User.ID,
			BotID:   ctx.Bot.Card().BotID,
		}
		if err = TimeDB.Insert(m); err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Text("记录指令成功!"))
	})
	en.AddRex(`^/删除指令(.*)`).Rule(rosm.OnlyMaster()).Handle(func(ctx *rosm.Ctx) {
		// 允许往正在执行的 cron 中添加任务
		if id, ok := entryIDMap[ctx.Being.Rex[1]]; ok {
			if err := TimeDB.Delete(tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2) + ctx.Being.Rex[1]); err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			c.Remove(id)
			delete(entryIDMap, ctx.Being.Rex[1])
		} else {
			ctx.Send(message.Text("未找到指令任务!"))
		}
		ctx.Send(message.Text("指令删除成功!"))
	})
	initDB(en.DataFolder + "time.db")
	go cronRun(TimeDB)
}

func cronRun(db *model) {
	tool.WaitInit()
	c = cron.New()
	db.Range(func(i int, m *mode) bool {
		timeCtx := send.CTXBuild(m.Types, m.BotID, m.String1)
		timeCtx.Being.Word = m.Word
		timeCtx.Being.User = &rosm.UserData{ID: m.UserID}
		id, _ := c.AddFunc(m.Time, func() {
			timeCtx.RunWord(m.Word)
		})
		entryIDMap[m.Word] = id
		return true
	})
	// 开始执行
	c.Start()
	logrus.Info("[time]cron任务开始执行...")
}
