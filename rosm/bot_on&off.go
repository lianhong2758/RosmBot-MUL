// 用于管理插件的启用/禁用
package rosm

// 仅作为存储开关数据使用,没有Matcher
var boten = Register(&PluginData{
	Name: "响应管理",
	Help: "- @bot/早安\n" +
		"- @bot/晚安",
	//借用插件管理的存储器
	DataFolder: "regulate",
})

func GetBoten() *PluginData {
	return boten
}
