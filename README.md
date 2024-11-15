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

- ~~[mys部署]( server/mys/README.md)~~
- [qq官方bot部署]( server/qq/README.md)
- [ob11协议部署]( server/ob11/README.md)

## 插件编写教程

1注册插件
```
func init() {
	en := c.Register("chat", &c.PluginData{//第一个参数是插件名,用于区分插件
		Name: "@回复",        			   //插件名,用于help
		Help: "- @机器人",				   //帮助信息,用于help
		DataFolder: "chat",				   //可选,创建插件的数据文件夹,不需要数据存储则不需要填写
	})
	//这里是匹配词------这里设置是否阻断继续匹配
	//还有.SetRule()设置指令初始化函数
    //MUL()设置插件专用的平台,一般在插件调用了对应平台的server包后填写
	en.AddWord("").SetBlock(true).SetRule(rosm.OnlyAtMe()).Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Text(ctx.Bot.Card().BotName, "不在呢~"))
	})
}
```
2获取触发时传送的数据
```
//ctx.Being里有所有需要的数据,结构如下
type Being struct {
	RoomID   string         //房间号
	RoomID2  string         //如果有需要,存放房间号上级号码
	RoomName string         //房间名称,存在上级则存放上级名称
	User     *UserData      //触发事件者信息
	ATList   []string       //at的id列表
	MsgID    []string       //用于reply,存放消息id,reply的其他需要参数写在第二位
	AtMe     bool           //是否是at机器人触发的事件
	Word     string         //接收的用户发送的信息,进行了首位的空格切割
	Rex      []string       //如果有正则匹配,这里存放匹配结果
	Def      map[string]any //自定义存储的参数
}
```
3发送消息
```
1)文本或者图片消息
ctx.Send(xxx)
xxx有很多,可以无限续接,逗号分开
message包里为通用的结构,可以在任意平台使用
其中文本消息用message.Text(any)
byte图片用message.Image(img []byte)
url图片用message.ImageUrl(url string)
at用message.AT(id , name string)
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
