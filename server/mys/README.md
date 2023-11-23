## MYS部署流程
```
1,在main.go中取消注释
_ "github.com/lianhong2758/RosmBot-MUL/server/mys/init"
注释掉	
_ "github.com/lianhong2758/RosmBot-MUL/server/qq/init"

2,配置你需要的插件,将不需要的注释掉

3,运行run.bat,选择连接协议,之后他会生成一个mys.config在config文件夹里面

4,填写你bot的IO,PubKey等米游社下发的数据,设置bot自称的Name,主人等信息

5,如果选择http协议,你需要在服务器控制平台开启10001端口
在米游社回调地址填入http://你的ip:10001/rosmbot

6,之后在米游社添加bot,@bot 进行测试,如果回复则说明部署成功
```

