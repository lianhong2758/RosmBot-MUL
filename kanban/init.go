package kanban

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
)

var isDebug bool

func init() {
	//log_file
	path := flag.String("lf", "data/log/"+time.Now().Format("2006_01_02_15_04_05")+".log", "log file path")
	notSaveLogs := flag.Bool("nolf", false, "不保存日志文件")
	//debug
	d := flag.Bool("d", false, "Enable debug level log and higher.")
	flag.Parse()
	if !*notSaveLogs {
		SetLogWithNewFile(*path)
	}
	isDebug = *d
	log.Debug("IN DEBUG MODE")
	//kanban
	Kanban()
}

func GetArgIsDebug() bool {
	return isDebug
}
