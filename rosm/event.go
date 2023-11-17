package rosm

import ()

type EventType = int

const (
	//消息有关类型1-2,事件3+
	AllMessage     EventType = iota //全消息匹配
	SurplusMessage                  //无其他插件匹配
	Join                            //入群
	Out                             //退群
	Create                          //添加Bot
	Delete                          //删除Bot
	Quick                           //表态-mys
	Audit                           //审核
	Click                           //点击事件
)
