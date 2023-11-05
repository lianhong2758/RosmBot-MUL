package main

import (
	_ "github.com/lianhong2758/RosmBot-MUL/kanban"
	"github.com/lianhong2758/RosmBot-MUL/rosm"

	//"github.com/lianhong2758/RosmBot-MUL/server/mys"
	_ "github.com/lianhong2758/RosmBot-MUL/server/mys/init"
	//_ "github.com/lianhong2758/RosmBot-MUL/server/qq/init"

	//_ "github.com/lianhong2758/RosmBot-MUL/plugins/myplugin"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/myplugin/chatgpt"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/myplugin/rsshub"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/test"
)

func main() {
	/*
		var runner rosm.Boter = mys.NewConfig("config/mys.json")
		go runner.Run()
	*/
	rosm.Listen()
}
