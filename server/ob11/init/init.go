package init

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/ob11"
)

var filePath = "config/ob11.json"

func init() { //main
	var runner rosm.Boter = ob11.NewConfig(filePath)
	go runner.Run()
}
