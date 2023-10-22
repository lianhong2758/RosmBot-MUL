package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/lianhong2758/RosmBot-MUL/kanban"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	//"github.com/lianhong2758/RosmBot-MUL/server/mys"
	//_ "github.com/lianhong2758/RosmBot-MUL/server/mys/init"
	_ "github.com/lianhong2758/RosmBot-MUL/server/qq/init"

	_ "github.com/lianhong2758/RosmBot-MUL/plugins/test"
)

func main() {
	kanban.Kanban()
	/*
		var runner rosm.Boter
			runner = &mys.MYSconfig
			go runner.Run()
	*/
	for {
		if i, ok := <-rosm.MULChan; ok {
			log.Printf("新增注册,平台: %s,昵称: %s,BotID: %s", i.Types, i.Name, i.BotID)
		}
	}
}
