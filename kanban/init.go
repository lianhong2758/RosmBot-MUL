package kanban

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

func init() {
	d := flag.Bool("d", false, "Enable debug level log and higher.")
	flag.Parse()
	if *d {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("IN DEBUG MODE")
}
