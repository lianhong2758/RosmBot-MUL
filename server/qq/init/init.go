package init

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/qq"
)

var filePath = "config/qq.json"

func init() { //main
	var runner rosm.Boter
	runner = qq.NewConfig(filePath)
	go runner.Run()
}
