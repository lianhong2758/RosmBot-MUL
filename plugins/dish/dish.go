package myplugin

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	sql "github.com/FloatTech/sqlite"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

type dish struct {
	ID        uint32 `db:"id"`
	Name      string `db:"name"`
	Materials string `db:"materials"`
	Steps     string `db:"steps"`
}

var (
	db          = &sql.Sqlite{}
	initialized = false
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name:       "程序员做饭指南",
		Help:       "- /随机菜谱|随便做点菜",
		DataFolder: "Dish",
	})

	db.DBPath = en.DataFolder + "dishes.db"

	if err := db.Open(time.Hour); err != nil {
		logrus.Warnln("[dish]连接菜谱数据库失败")
	} else if err = db.Create("dish", &dish{}); err != nil {
		logrus.Warnln("[dish]同步菜谱数据表失败")
	} else if count, err := db.Count("dish"); err != nil {
		logrus.Warnln("[dish]统计菜谱数据失败")
	} else {
		logrus.Infoln("[dish]加载", count, "条菜谱")
		initialized = true
	}

	if !initialized {
		logrus.Warnln("[dish]插件未能成功初始化")
	}

	en.AddWord("/随机菜谱", "/随便做点菜").Handle(func(ctx *rosm.Ctx) {
		var d dish
		if err := db.Pick("dish", &d); err != nil {
			ctx.Send(message.Text("小店好像出错了，暂时端不出菜来惹"))
			logrus.Warnln("[dish]随机菜谱请求出错：" + err.Error())
			return
		}
		ctx.Send(message.Text("已为客官"), message.AT(ctx.Being.User.ID, ctx.Being.User.Name), message.Text(fmt.Sprintf("送上%s的做法：\n原材料：%s\n步骤：\n%s", d.Name, d.Materials, d.Steps)))
	})
}
