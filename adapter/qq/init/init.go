package init

import (
	"github.com/lianhong2758/RosmBot-MUL/adapter/qq"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
)

var filePath = "qq.json"

func init() { //main
	var runner rosm.Boter = qq.NewConfig(filePath)
	go runner.Run()
}
