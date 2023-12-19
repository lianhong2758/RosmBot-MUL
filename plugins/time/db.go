package time

import (
	"os"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var TimeDB = &model{sql: &sql.Sqlite{}}

type model struct {
	sql *sql.Sqlite
	sync.RWMutex
}

type mode struct {
	Key     string `db:"Key"`     //用string1+word拼接作为key,避免重复
	String1 string `db:"String1"` //tool-String1
	Types   string `db:"Types"`   //平台
	Word    string `db:"Word"`    //记录的指令
	Time    string `db:"Time"`    //定时
	UserID  string `db:"UserID"`  //记录者ID
	BotID   string `db:"BotID"`   //优先匹配的botid
}

// 初始化
func initDB(path string) {
	TimeDB.sql.DBPath = path
	err := TimeDB.sql.Open(time.Hour * 24)
	if err != nil {
		log.Error("[time]初始化数据库失败!")
		os.Exit(1)
	}
	err = TimeDB.sql.Create("time", &mode{})
	if err != nil {
		logrus.Error("[time]初始化表(time)失败: ", err)
		os.Exit(1)
	}
}

func (db *model) Insert(m *mode) (err error) {
	db.Lock()
	defer db.Unlock()
	return db.sql.Insert("time", m)
}

func (db *model) Delete(key string) (err error) {
	db.Lock()
	defer db.Unlock()
	return db.sql.Del("time", "where Key = '"+key+"'")
}

func (db *model) Find(key string) (m *mode, err error) {
	m = new(mode)
	db.Lock()
	defer db.Unlock()
	err = db.sql.Find("time", &m, "where Key = '"+key+"'")
	return
}

func (db *model) Range(f func(i int, m *mode) bool) {
	// 执行查询语句
	rows, err := db.sql.DB.Query("SELECT * FROM " + "time")
	if err != nil {
		return
	}
	defer rows.Close()

	var i int = 0
	// 遍历结果
	for rows.Next() {
		var m mode
		err := rows.Scan(&m.Key, &m.String1, &m.Types, &m.Word, &m.Time, &m.UserID, &m.BotID)
		if err != nil {
			return
		}
		if !f(i, &m) {
			return
		}
		i++
	}
}
