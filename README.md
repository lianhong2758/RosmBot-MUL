## RosmBot_MUL(迷迭香Bot)
RosmBot-MUL是一个多平台bot,未来将实现一个插件多平台接入使用,本项目由golang编写


## 命令行参数
> `[]`代表是可选参数
```bash
RosmBot-Mul [-d] [-lf path] [-nolf]
```
- **-d**: debug模式
- **-lf data/log/test.txt**: 自定义本次的日志存储位置
- **-nolf**: 本次不保存日志

## 使用方法

直接运行
```
	进入main.go
	注释或取消注释掉需要的平台或插件
	运行run.bat
```
之后运行即可

## 部署方法

- ~~[mys部署]~~
- [qq官方bot部署]( server/qq/README.md)
- [ob11协议部署]( server/ob11/README.md)

## 插件编写教程

1注册插件
```
func init() {
	en := c.Register("chat", &c.PluginData{//第一个参数是插件名,用于区分插件
		DefaultOff: false,    //是否默认禁用本插件
		Name:       "@回复",    //插件名,用于help
		Help:       "- @机器人", //帮助信息,用于help
		//DataFolder: "chat",				   //可选,创建插件的数据文件夹,不需要数据存储则不需要填写
		//Config: map[string]string{"1": "1"},//可选,配置每个插件独立的config,传入带有默认值的结构体指针即可
	})
	//这里是匹配词------这里设置是否阻断继续匹配
	//还有.SetRule()设置指令初始化函数
    //MUL()设置插件专用的平台,一般在插件调用了对应平台的server包后填写
	//完全词匹配,因为前缀config默认为[]string{"/",""},所以这里可以匹配 "" 或者 "/"
	en.OnWord("").SetBlock(true).SetRule(rosm.OnlyAtMe()).Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Text(rosm.GetRandBotName(), "不在呢~"))
	})
	//除此之外还有别的匹配方法
	//匹配正则
	en.OnRex(`^/你好`).Handle(func(ctx *rosm.Ctx) {})
	//匹配事件rosm.FriendRecall
	en.OnNoticeWithType(rosm.FriendRecall).Handle(func(ctx *rosm.Ctx) {})

}
```
2获取触发时传送的数据
```
//ctx.Being里有所有需要的数据,结构如下
type Being struct {
	GroupID    string    //房间号,群号
	GuildID    string    //如果有需要,存放房间号上级号码
	GroupName  string    //房间名称,存在上级则存放上级名称
	User       *UserData //触发事件者信息
	ATList     []string  //at的id列表
	MsgID      string    //用于reply,存放消息id
	IsAtMe     bool      //是否是at机器人触发的事件
	RawWord    string    //原始消息
	KeyWord    string    //匹配词
	ResultWord []string  //匹配结果
}

//ctx中有上下文传递的信息
type Ctx struct {
	Bot     Boter
	Message message.Message //Message
	Being   *Being          //常用消息解析,需实现
	State   map[string]any  //框架中途产生的内容,不固定,即时有效
}
```
3发送消息
```
1)文本或者图片消息
ctx.Send(xxx)
xxx有很多,可以无限续接,逗号分开
message包里为通用的结构,可以在任意平台使用
其中文本消息用message.Text(any)
byte图片用message.ImageByte(img []byte)
at用message.AT(id string)
reply用message.reply()
其余看源码学习...
```
4更改发送房间
```
ctx.Being.GroupID/RoomID2用于发送消息的房间索引,可以直接修改这里的数据
```
5部分接口(可能存在没有及时更新,导致调用出错的情况,如有请反馈)
```
//接口存放在server/平台id/openapi里面,请导入对应的包进行使用
//举例
import	"github.com/lianhong2758/RosmBot-MUL/server/mys"
result, err := mys.GetRoomList(ctx)
```

6启用插件
```
如果编写的插件没有在plugins/test里面,请手动在main.go里面进行导入注册
```
## 特别鸣谢
[ZeroBot](https://github.com/wdvxdr1123/ZeroBot)提供部分代码借鉴
## 相关地址

- QQ交流群 : 678586912

- Github : https://github.com/lianhong2758

- [QQ官方Bot文档](https://bot.q.qq.com/wiki/develop/api-v2/)

- [QQ开放平台](https://q.qq.com)

- [NTQQ-llonebot](https://llonebot.github.io/zh-CN/)
