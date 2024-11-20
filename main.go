package main

import (
	//必须的依赖
	_ "github.com/lianhong2758/RosmBot-MUL/kanban"
	"github.com/lianhong2758/RosmBot-MUL/rosm"

	//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	//Bot注册
	//_"github.com/lianhong2758/RosmBot-MUL/adapter/mys/init"
	//_ "github.com/lianhong2758/RosmBot-MUL/adapter/qq/init"
	//"github.com/lianhong2758/RosmBot-MUL/adapter/ob11"
	_ "github.com/lianhong2758/RosmBot-MUL/adapter/ob11/init"
	//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	//插件注册
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/chatgpt"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/fhl"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/gscore"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/lc"
	//_ "github.com/lianhong2758/RosmBot-MUL/plugins/liuyao"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/onplugin"
	//	_ "github.com/lianhong2758/RosmBot-MUL/plugins/phi"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/score"
	//_ "github.com/lianhong2758/RosmBot-MUL/plugins/tarot"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/test"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/time"
	//_ "github.com/lianhong2758/RosmBot-MUL/plugins/wife"
	_ "github.com/lianhong2758/RosmBot-MUL/plugins/yujn"
	//^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
)

func main() {
	rosm.RosmInit()
	//下面两行用于同平台多bot构建,不需要可以忽略
	// var runner rosm.Boter = ob11.NewConfig("config/ob11.json")
	// go runner.Run()

	//阻塞主进程
	rosm.Listen()
}
