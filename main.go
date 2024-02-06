package main

import (
	//必须的依赖
	_ "github.com/lianhong2758/RosmBot-MUL/kanban"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	//Bot注册
	//"github.com/lianhong2758/RosmBot-MUL/server/mys"
	//_"github.com/lianhong2758/RosmBot-MUL/server/mys/init"
	_ "github.com/lianhong2758/RosmBot-MUL/server/qq/init"
	//_ "github.com/lianhong2758/RosmBot-MUL/server/ob11/init"
	//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	//插件注册
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/chatgpt"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/gscore"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/test"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/time"
	//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
)

func main() {
	//下面两行用于同平台多bot构建,不需要可以忽略
	// var runner rosm.Boter = mys.NewConfig("config/mys.json")
	// go runner.Run()

	//阻塞主进程
	rosm.Listen()
}
