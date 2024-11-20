// Rule的实现,可以在这里增加更多Rule,也可以在server包增加独属于自己平台的rule
package rosm

import (
	"regexp"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/message"
)

// 判断rule
func (m *Matcher) RulePass(ctx *Ctx) bool {
	return rulePass(ctx, m.rules...)
}

func rulePass(ctx *Ctx, rs ...Rule) bool {
	for _, v := range rs {
		if !v(ctx) {
			return false
		}
	}
	return true
}

func OnlyAtMe() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Being.IsAtMe
	}
}

func NoAtForOther() Rule {
	return func(ctx *Ctx) bool {
		//atme,noat
		return ctx.Being.IsAtMe || len(ctx.Being.ATList) == 0
	}
}

func OnlyTheRoom(roomid, roomid2 string) Rule {
	return func(ctx *Ctx) bool {
		return roomid == ctx.Being.GroupID && roomid2 == ctx.Being.GuildID
	}
}

func OnlyMaster() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Bot.OnlyMaster(ctx)
	}
}

// 用于getnext
func OnlyTheUser(id string) Rule {
	return func(ctx *Ctx) bool {
		return id == ctx.Being.User.ID
	}
}

func OnlyReply() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Bot.OnlyReply(ctx)
	}
}

// 大于等于群主等权限
func OnlyOverHost() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Bot.OnlyOverHost(ctx)
	}
}

// 大于等于管理员等权限
func OnlyOverAdministrator() Rule {
	return func(ctx *Ctx) bool {
		return ctx.Bot.OnlyOverAdministrator(ctx)
	}
}
func KeyWords(s ...string) Rule {
	return func(ctx *Ctx) bool {
		for _, str := range s {
			if strings.Contains(ctx.Being.KeyWord, str) {
				ctx.State["keyword"] = str
				return true
			}
		}
		return false
	}
}

// 完全匹配词
func WordRule(words ...string) Rule {
	return func(ctx *Ctx) bool {
		for _, v := range words {
			if v == ctx.Being.RawWord {
				ctx.Being.KeyWord = v
				return true
			}
		}
		return false
	}
}

// Rex匹配
func RexRule(rex string) Rule {
	r := regexp.MustCompile(rex)
	return func(ctx *Ctx) bool {
		if match := r.FindStringSubmatch(ctx.Being.RawWord); len(match) > 0 {
			ctx.Being.ResultWord = match
			return true
		}
		return false
	}
}

// Notice事件匹配
func NoticeRule(types ...string) Rule {
	return func(ctx *Ctx) bool {
		for _, v := range types {
			if ctx.State["notice_type"].(string) == v {
				return true
			}
		}
		return false
	}
}

// 模板匹配
// 一般使用UniMessage内的函数作为 message.MessageSegment
func TemplateRule(msg ...message.MessageSegment) Rule {
	rexList := []*regexp.Regexp{}
	for _, v := range msg {
		if v.Type == "rex" {
			rexList = append(rexList, regexp.MustCompile(v.Data["text"]))
		}
	}
	return func(ctx *Ctx) bool {
		ok := false
		ir := 0 //rex存储表的索引
	Next:
		for mi, i := 0, 0; mi < len(msg) && i < len(ctx.Message); {
			switch msg[mi].Type {
			case "any":
				ctx.Being.ResultWord = append(ctx.Being.ResultWord, ctx.Message[i].Text())
				mi++
				i++
				continue
			case "at":
				if ctx.Message[i].Type != "at" {
					ok = false
					break Next
				}
				if msg[mi].Data["uid"] != "" {
					for _, v := range ctx.Message[i].Data {
						if v == msg[mi].Data["uid"] {
							ctx.Being.ResultWord = append(ctx.Being.ResultWord, msg[mi].Data["uid"])
							ok = true
							i++
							mi++
							continue Next
						}
					}
					//匹配失败
					ok = false
					break Next
				}
				if id := ctx.Message[i].AtId(); id == "" {
					ok = false
					break Next
				} else {
					ctx.Being.ResultWord = append(ctx.Being.ResultWord, id)
					ok = true
					mi++
					i++
				}
			case "rex":
				if ctx.Message[i].Type != "text" {
					ok = false
					break Next
				}
				if match := rexList[ir].FindStringSubmatch(ctx.Message[i].Data["text"]); len(match) > 0 {
					ctx.Being.ResultWord = append(ctx.Being.ResultWord, match[1:]...)
					ok = true
					i++
					mi++
					ir++
				} else {
					ok = false
				}
			case "text":
				if ctx.Message[i].Type != "text" {
					ok = false
					break Next
				}
				if ctx.Message[i].TrimSpaceText() != msg[i].Data["text"] {
					ok = false
					break Next
				}
				i++
				mi++
			default:
				if ctx.Message[i].Type != msg[i].Type {
					ok = false
					break Next
				}
				ctx.Being.ResultWord = append(ctx.Being.ResultWord, ctx.Message[i].Text())
				ok = true
				i++
				mi++
			}
		}
		return ok
	}
}
