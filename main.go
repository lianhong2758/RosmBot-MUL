package main

import (
	_ "github.com/lianhong2758/RosmBot-MUL/kanban"
	"github.com/lianhong2758/RosmBot-MUL/rosm"

	//"github.com/lianhong2758/RosmBot-MUL/server/mys"
	_ "github.com/lianhong2758/RosmBot-MUL/server/mys/init"
	//_ "github.com/lianhong2758/RosmBot-MUL/server/qq/init"

	//_ "github.com/lianhong2758/RosmBot-MUL/plugins/myplugin"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/chatgpt"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/gscore"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/test"
)

func main() {
	/*
		var runner rosm.Boter = mys.NewConfig("config/qq2.json")
		go runner.Run()
	*/
	rosm.Listen()
}
