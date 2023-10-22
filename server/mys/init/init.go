package init

import (
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/server/mys"
)

func init() { //main
	var runner rosm.Boter
	runner = &mys.MYSconfig
	go runner.Run()
}
