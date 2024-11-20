// 插件的具体结构
package rosm

import (
	"os"
	"sort"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/ttl"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/tool/rate"
	log "github.com/sirupsen/logrus"
)

var (
	//用于映射发送的消息id到触发id
	MessagesMapCache = ttl.NewCache[string, []string](time.Minute * 5)
	MessagesMu       = sync.Mutex{}
)

type (
	// Rule filter the event
	Rule = func(ctx *Ctx) bool
	// Handler 事件处理函数
	Handler = func(ctx *Ctx)
	//meta_event,message,notice,request,all,surplus
	EventType = string
)

const (
	AllMessage     = "all"
	SurplusMessage = "surplus"
)

type PluginData struct {
	DefaultOff bool //是否默认禁用
	Help       string
	Name       string
	DataFolder string //"data/xxx/"+
	Config     any
	Matchers   []*Matcher //子插件
}
type Matcher struct {
	types      string //匹配事件类型
	rules      []Rule
	handler    Handler
	priority   int         //优先级控制,用于插件的匹配顺序,范围0-9,默认5,数字越小优先级越高,wordRule无效
	mul        []string    //可用平台
	block      bool        //阻断
	temp       bool        //是否是临时Matcher
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

	//事件触发,meta_event,message,notice,request,all,surplus
	EventMatch = map[EventType][]*Matcher{} //事件触发
)

// 注册插件
func Register(p *PluginData) *PluginData {
	p.Matchers = make([]*Matcher, 0)
	pluginName := p.Name
	log.Debugln("[debug]插件注册:", pluginName)
	plugins[pluginName] = p
	if p.DataFolder != "" && file.IsNotExist("data/"+p.DataFolder) {
		_ = os.MkdirAll("data/"+p.DataFolder, 0755)
	}
	plugins[pluginName].DataFolder = "data/" + p.DataFolder + "/"
	if p.DataFolder != "" && p.Config != nil {
		if err := LoadConfig(p.DataFolder+"config.json", p.Config); err != nil {
			log.Warnln("[warn]插件注册:", err)
		}
	}
	return plugins[pluginName]
}

// 创建插件对象信息
func NewRegist(name, help, dataFolder string, DefaultOff bool, config any) *PluginData {
	return &PluginData{Name: name, Help: help, DataFolder: dataFolder, Config: config, DefaultOff: DefaultOff}
}

// 元事件
func (p *PluginData) OnMetaEvent() *Matcher {
	return p.On("meta_event")
}

// message
func (p *PluginData) OnMessage() *Matcher {
	return p.On("message").SetBlock(true)
}

// notice
func (p *PluginData) OnNotice() *Matcher {
	return p.On("notice")
}

// request
func (p *PluginData) OnRequest() *Matcher {
	return p.On("request")
}

// on 注册一个基础事件响应器，可自定义类型
func (p *PluginData) On(types string) *Matcher {
	m := &Matcher{
		types:      types,
		block:      false,
		priority:   5,
		temp:       false,
		PluginNode: p,
		rules:      []Rule{GetBotIsOnInThis()},
	}
	EventMatch[types] = append(EventMatch[types], m)
	p.Matchers = append(p.Matchers, m)
	return m
}

// 注册一个消息事件响应器，匹配word消息
func (p *PluginData) OnWord(words ...string) *Matcher {
	//这里用map跳表匹配,所以不用添加匹配的Rule
	m := p.OnMessage() //p.OnMessage().SetRule(WordRule(words...))
	for _, word := range words {
		for _, prefix := range config.CmdStar {
			WordMatch[prefix+word] = m
		}
	}
	return m
}

// 注册一个消息事件响应器，匹配rex消息
func (p *PluginData) OnRex(rex string) *Matcher {
	m := p.OnMessage().SetRule(RexRule(rex))
	RegexpMatch = append(RegexpMatch, m)
	return m
}

// 注册一个消息事件响应器，匹配rex消息
func (p *PluginData) OnNoticeWithType(types ...string) *Matcher {
	m := p.OnNotice().SetRule(NoticeRule(types...))
	return m
}

// 框架独特的事件,匹配所有消息,但不阻断
func (p *PluginData) OnAllMessage() *Matcher {
	m := p.On("all").SetBlock(false)
	return m
}

// 框架独特的事件,匹配没有插件匹配的消息,类似最低优先级,不阻断
func (p *PluginData) OnSurplusMessage() *Matcher {
	m := p.On("surplus").SetBlock(false)
	return m
}

// 模板匹配消息
func (p *PluginData) OnTemplate(msg ...message.MessageSegment) *Matcher {
	m := p.OnAllMessage().SetRule(TemplateRule(msg...))
	return m
}

// "data/xxx/"+
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
	m.SetRule(m.mulRule(), MatcherIsOn(m))
	//执行hander
	m.handler = h
}

// 阻断器
func (m *Matcher) SetBlock(ok bool) *Matcher {
	m.block = ok
	return m
}

func (m *Matcher) SetRule(r ...Rule) *Matcher {
	m.rules = append(m.rules, r...)
	return m
}

func (m *Matcher) mulRule() Rule {
	return func(ctx *Ctx) bool {
		if len(m.mul) == 0 {
			return true
		}
		for _, value := range m.mul {
			if value == ctx.Bot.Card().BotType {
				return true
			}
		}
		return false
	}
}

// Limit 限速器
// postfn 当请求被拒绝时的操作
func (m *Matcher) Limit(limiterfn func(*Ctx) *rate.Limiter, postfn ...func(*Ctx)) *Matcher {
	m.SetRule(func(ctx *Ctx) bool {
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
	MessagesMu.Lock()
	defer MessagesMu.Unlock()
	h := ctx.Bot.BotSend(ctx, m...)
	if h["id"] != "" {
		MessagesMapCache.Set(ctx.Being.MsgID, append([]string{h["id"]}, MessagesMapCache.Get(ctx.Being.MsgID)...))
	}
	return h
}

// 通过记录的回复id查找触发id
func GetMessageIDFormMapCache(id string) []string {
	MessagesMu.Lock()
	defer MessagesMu.Unlock()
	return MessagesMapCache.Get(id)
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

// Rosm Match排序
func RosmInit() {
	//正则排序
	sort.Slice(RegexpMatch, func(i, j int) bool {
		return RegexpMatch[i].priority < RegexpMatch[j].priority
	})
	//event排序
	for _, mas := range EventMatch {
		sort.Slice(mas, func(i, j int) bool {
			return mas[i].priority < mas[j].priority
		})
	}
}

func (p *PluginData) GetConfig() any {
	return p.Config
}

func DeleteOffRule(p *PluginData) {
	for _, m := range p.Matchers {
		m.rules = m.rules[1 : len(m.rules)-1]
	}
}
