package kanban
import log "github.com/sirupsen/logrus"

func init() {
	log.SetFormatter(&ColorNotFormatter{})
}
