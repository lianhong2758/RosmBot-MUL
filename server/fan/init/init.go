package init

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/fan"
)

var filePath = "config/fan.json"

func init() { //main
	var runner rosm.Boter = fan.NewConfig(filePath)
	go runner.Run()
}
