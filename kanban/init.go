package kanban

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	//log_file
	path := flag.String("lf", "data/log/"+time.Now().Format("2006_01_02_15_04_05")+".log", "log file path")
	notSaveLogs := flag.Bool("nolf", false, "不保存日志文件")
	if !*notSaveLogs {
		SetLogWithNewFile(*path)
	}
	//debug
	d := flag.Bool("d", false, "Enable debug level log and higher.")
	flag.Parse()
	if *d {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("IN DEBUG MODE")
	//kanban
	Kanban()
}
