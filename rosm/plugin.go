// 插件的具体结构
package rosm

import (
	"os"
	"regexp"
	"sort"

	"github.com/FloatTech/floatbox/file"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/tool/rate"
	log "github.com/sirupsen/logrus"
)

type (
	// Rule filter the event
	Rule = func(ctx *Ctx) bool
	// Handler 事件处理函数
	Handler = func(ctx *Ctx)
)

type PluginData struct {
	DefaultOff bool //是否默认禁用
	Priority   int  //优先级控制,仅对Rex有效,用于插件的匹配顺序,范围0-9,默认5,数字越小优先级越高
	Help       string
	Name       string
	DataFolder string //"data/xxx/"+
	Matchers   []*Matcher
}
type Matcher struct {
	Word       []string
	Rex        *regexp.Regexp
	rules      []Rule
	handler    Handler
	mul        []string
	block      bool        //阻断
	nextchan   chan *Ctx   //用于上下文
	PluginNode *PluginData //溯源
}

// 插件注册
var (
	plugins = map[string]*PluginData{} //所有的插件集合

	// 全匹配字典
	WordMatch = map[string]*Matcher{} //Word

	// 正则字典
	RegexpMatch = []*Matcher{} //Rex

	//事件触发
	EventMatch = map[int][]*Matcher{} //事件触发
)

// 注册插件
func Register(p *PluginData) *PluginData {
	pluginName := p.Name
	log.Debugln("插件注册:", pluginName)
	plugins[pluginName] = p
	if p.DataFolder != "" && file.IsNotExist("data/"+p.DataFolder) {
		_ = os.MkdirAll("data/"+p.DataFolder, 0755)
	}
	plugins[pluginName].DataFolder = "data/" + p.DataFolder + "/"
	return plugins[pluginName]
}

// 创建插件对象信息
func NewRegist(name, help, dataFolder string) *PluginData {
	return &PluginData{Name: name, Help: help, DataFolder: dataFolder}
}

// 完全词匹配
func (p *PluginData) AddWord(word ...string) *Matcher {
	m := new(Matcher)
	m.block = true
	m.Word = append(m.Word, word...)
	for _, v := range word {
		WordMatch[v] = m
	}
	p.Matchers = append(p.Matchers, m)
	m.PluginNode = p
	return m
}

// 正则匹配
func (p *PluginData) AddRex(rex string) *Matcher {
	m := new(Matcher)
	m.block = true
	m.Rex = regexp.MustCompile(rex)
	RegexpMatch = append(RegexpMatch, m)
	p.Matchers = append(p.Matchers, m)
	m.PluginNode = p
	return m
}

// 其他事件匹配器
func (p *PluginData) AddEvent(types int) *Matcher {
	m := new(Matcher)
	m.block = false
	EventMatch[types] = append(EventMatch[types], m)
	p.Matchers = append(p.Matchers, m)
	m.PluginNode = p
	return m
}
func (p *PluginData) GetFolder() string {
	return p.DataFolder
}

// 平台限制
func (m *Matcher) MUL(name ...string) *Matcher {
	m.mul = append(m.mul, name...)
	return m
}

// 注册Handle
func (m *Matcher) Handle(h Handler) {
	//加载默认的rule
	//全局bot启动＋插件单独启用
	m.rules = append([]func(ctx *Ctx) bool{GetBotIsOnInThis()}, m.rules...)
	m.rules = append(m.rules, m.mulPass(), MatcherIsOn(m))
	//执行hander
	m.handler = h
}

// 阻断器
func (m *Matcher) SetBlock(ok bool) *Matcher {
	m.block = ok
	return m
}

func (m *Matcher) Rule(r ...Rule) *Matcher {
	m.rules = append(m.rules, r...)
	return m
}

func (m *Matcher) mulPass() Rule {
	return func(ctx *Ctx) bool {
		if len(m.mul) == 0 {
			return true
		}
		for _, value := range m.mul {
			if value == ctx.BotType {
				return true
			}
		}
		return false
	}
}

// Limit 限速器
// postfn 当请求被拒绝时的操作
func (m *Matcher) Limit(limiterfn func(*Ctx) *rate.Limiter, postfn ...func(*Ctx)) *Matcher {
	m.rules = append(m.rules, func(ctx *Ctx) bool {
		if limiterfn(ctx).Acquire() {
			return true
		}
		if len(postfn) > 0 {
			for _, fn := range postfn {
				fn(ctx)
			}
		}
		return false
	})
	return m
}

// 快捷发送消息
func (ctx *Ctx) Send(m ...message.MessageSegment) H {
	return ctx.Bot.BotSend(ctx, m...)
}
func Display() {
	log.Println(WordMatch)
	log.Println(RegexpMatch)
}
func GetPlugins() map[string]*PluginData {
	return plugins
}

func GetBotIsOnInThis() func(*Ctx) bool {
	return func(ctx *Ctx) bool {
		return ctx.on
	}
}

// Rosm RegexpMatch排序
func RosmInit() {
	sort.Slice(RegexpMatch, func(i, j int) bool {
		return RegexpMatch[i].PluginNode.Priority < RegexpMatch[j].PluginNode.Priority
	})
}
