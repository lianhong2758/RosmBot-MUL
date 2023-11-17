package rosm

import (
	"os"
	"regexp"

	"github.com/FloatTech/floatbox/file"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/tool/rate"
	log "github.com/sirupsen/logrus"
)

type (
	// Rule filter the event
	Rule = func(ctx *CTX) bool
	// Handler 事件处理函数
	Handler = func(ctx *CTX)
)

type PluginData struct {
	Help       string
	Name       string
	DataFolder string //"data/xxx/"+
	Matchers   []*Matcher
}
type Matcher struct {
	Word       []string
	Rex        []*regexp.Regexp
	rules      []Rule
	handler    Handler
	mul        []string
	block      bool        //阻断
	nestchan   chan *CTX   //用于上下文
	PluginNode *PluginData //溯源
}

// 插件注册
var (
	plugins = map[string]*PluginData{}

	// 全匹配字典
	caseAllWord = map[string]*Matcher{}

	// 正则字典
	caseRegexp = map[*regexp.Regexp]*Matcher{}

	//事件触发
	caseEvent = map[int][]*Matcher{}
)

// 注册插件
func Register(p *PluginData) *PluginData {
	pluginName := p.Name
	log.Debugln("插件注册:", pluginName)
	plugins[pluginName] = p
	if p.DataFolder != "" && file.IsNotExist(p.DataFolder) {
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
		caseAllWord[v] = m
	}
	p.Matchers = append(p.Matchers, m)
	m.PluginNode = p
	return m
}

// 正则匹配
func (p *PluginData) AddRex(rex string) *Matcher {
	m := new(Matcher)
	m.block = true
	r := regexp.MustCompile(rex)
	m.Rex = append(m.Rex, r)
	caseRegexp[r] = m
	p.Matchers = append(p.Matchers, m)
	m.PluginNode = p
	return m
}

// 其他事件匹配器
func (p *PluginData) AddOther(types int) *Matcher {
	m := new(Matcher)
	m.block = false
	caseEvent[types] = append(caseEvent[types], m)
	p.Matchers = append(p.Matchers, m)
	m.PluginNode = p
	return m
}

// 平台限制
func (m *Matcher) MUL(name ...string) *Matcher {
	m.mul = append(m.mul, name...)
	return m
}

// 注册Handle
func (m *Matcher) Handle(h Handler) {
	m.handler = h
}

// 阻断器
func (m *Matcher) SetBlock(ok bool) *Matcher {
	m.block = ok
	return m
}

func (m *Matcher) Rule(r ...Rule) *Matcher {
	m.rules = append(append(m.rules, m.mulPass()), r...)
	return m
}

func (m *Matcher) mulPass() Rule {
	return func(ctx *CTX) bool {
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
func (m *Matcher) Limit(limiterfn func(*CTX) *rate.Limiter, postfn ...func(*CTX)) *Matcher {
	m.rules = append(m.rules, func(ctx *CTX) bool {
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
func (ctx *CTX) Send(m ...message.MessageSegment) any {
	return ctx.Bot.BotSend(ctx, m...)
}
func Display() {
	log.Println(caseAllWord)
	log.Println(caseRegexp)
}
func GetPlugins() map[string]*PluginData {
	return plugins
}
