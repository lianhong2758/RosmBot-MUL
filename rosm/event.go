package rosm

// 事件类型的集合
//新的平台的事件类型都需要加入其中再实现
//无须在意顺序,但还是建议消息事件排最前,同一平台的事件放在一起
//可以把有相同性质的事件作为一个事件类型
//可以通过平台限制+事件类型,实现更具体的事件监听

type EventType = int

const (
	//消息有关类型1-2,事件3+
	AllMessage     EventType = iota //全消息匹配
	SurplusMessage                  //无其他插件匹配
	Join                            //入群
	Out                             //退群
	UpFile                          //上传文件-ob11
)
