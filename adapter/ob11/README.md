## ob11部署流程
1,在main.go的 `//Bot注册` 中

取消注释:
```
_ "github.com/lianhong2758/RosmBot-MUL/server/ob11/init"
```
注释掉其余导入

2,在 `//插件注册` 中注释掉不需要的插件

3,运行run.bat,第一次会生成config文件[xx.json]到config文件夹中 

4,(可选)设置bot自称的Name,主人等信息

5,之后再次运行run.bat就可以连接ob11实现了

## [ob11实现部署流程](https://llonebot.github.io/zh-CN)
