package init

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
)

var filePath = "config/mys.json"

func init() { //main
	var runner rosm.Boter
	runner = mys.NewConfig(filePath)
	go runner.Run()
}
