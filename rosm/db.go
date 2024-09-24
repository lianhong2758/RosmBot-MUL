//db用于存储插件启用信息
//Other字段可以用于存储插件的文本数据,不用独自建立数据库
//复杂数据需要插件独自建立库

package rosm

import (
	"os"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

var PluginDB = &model{sql: &sql.Sqlite{}}

func init() {
	PluginDB.sql.DBPath = en.DataFolder + "regulate.db"
	err := PluginDB.sql.Open(time.Hour * 24)
	if err != nil {
		log.Error("[regulate]初始化数据库失败!")
		os.Exit(1)
	}
}

type model struct {
	sql *sql.Sqlite
	sync.RWMutex
}

type mode struct {
	RoomID string `db:"RoomID"`
	Off    bool   `db:"Off"`    //false为开启插件,true为关闭插件
	String string `db:"String"` //String类型数据
	Int    int    `db:"Int"`    //int类型数据
}

// 初始化插件表
func plugindbinit() {
	for _, plugin := range GetPlugins() {
		_ = PluginDB.sql.Create(plugin.Name, &mode{})
	}
}
func (db *model) InsertOff(pluginname, roomid string, off bool) (err error) {
	oldm, _ := db.Find(pluginname, roomid)
	m := mode{
		RoomID: roomid,
		Off:    off,
		String: oldm.String,
		Int:    oldm.Int,
	}
	return db.sql.Insert(pluginname, &m)
}

func (db *model) FindOff(pluginname, roomid string) (off bool, err error) {
	m, err := db.Find(pluginname, roomid)
	return m.Off, err

}

func (db *model) Delete(pluginname, roomid string) (err error) {
	db.Lock()
	defer db.Unlock()
	return db.sql.Del(pluginname, "where RoomID = '"+roomid+"'")
}

func (db *model) Find(pluginname, roomid string) (m mode, err error) {
	db.Lock()
	defer db.Unlock()
	err = db.sql.Find(pluginname, &m, "where RoomID = '"+roomid+"'")
	return
}

func (db *model) FindString(pluginname, roomid string) (s string, err error) {
	m, err := db.Find(pluginname, roomid)
	return m.String, err
}

func (db *model) InsertString(pluginname, roomid string, s string) (err error) {
	oldm, _ := db.Find(pluginname, roomid)
	m := mode{
		RoomID: roomid,
		Off:    oldm.Off,
		String: s,
		Int:    oldm.Int,
	}
	return db.sql.Insert(pluginname, &m)
}
func (db *model) FindInt(pluginname, roomid string) (i int, err error) {
	m, err := db.Find(pluginname, roomid)
	return m.Int, err
}

func (db *model) InsertInt(pluginname, roomid string, i int) (err error) {
	oldm, _ := db.Find(pluginname, roomid)
	m := mode{
		RoomID: roomid,
		Off:    oldm.Off,
		String: oldm.String,
		Int:     i,
	}
	return db.sql.Insert(pluginname, &m)
}
// 查询是否开启插件
func MatcherIsOn(m *Matcher) func(ctx *Ctx) bool {
	return func(ctx *Ctx) bool {
		off, err := PluginDB.FindOff(m.PluginNode.Name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2))
		log.Debugln("[db]PluginIsOn 插件:", m.PluginNode.Name, "Off: ", off, "err: ", err)
		return (!off && err == nil) || (!m.PluginNode.DefaultOff && err == sql.ErrNullResult)
	}
}

// 查询是否开启插件,传入Plugin
func PluginIsOn(PluginNode *PluginData) func(ctx *Ctx) bool {
	return func(ctx *Ctx) bool {
		off, err := PluginDB.FindOff(PluginNode.Name, tool.MergePadString(ctx.Being.RoomID, ctx.Being.RoomID2))
		log.Debugln("[db]PluginIsOn 插件:", PluginNode.Name, "Off: ", off, "err: ", err)
		return (!off && err == nil) || (!PluginNode.DefaultOff && err == sql.ErrNullResult)
	}
}
