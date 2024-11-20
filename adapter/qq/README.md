## QQ部署流程
1,在main.go的 `//Bot注册` 中

取消注释:
```
_ "github.com/lianhong2758/RosmBot-MUL/server/qq/init"
```
注释掉其余导入

2,在 `//插件注册` 中注释掉不需要的插件

3,运行run.bat,第一次会生成config文件[xx.json]到config文件夹中 

4,填写你bot的id,key等信息,设置bot自称的Name,主人等信息,填写订阅内容

订阅内容去[官方文档](https://bot.q.qq.com/wiki/develop/api/gateway/intents.html)查看

5,再次运行run.bat之后,在群/频道按照开放平台的要求添加bot进行测试即可

