package rosm

import ()

const (
	AllMessage     int = iota //全消息匹配
	SurplusMessage            //无其他插件匹配
	Join                      //入群
	Create                    //添加Bot
	Delete                    //删除Bot
	Quick                     //表态-mys
)
