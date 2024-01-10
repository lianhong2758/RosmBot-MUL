package driver

import (
	"strconv"

	"github.com/tidwall/gjson"
)

// Modified from https://github.com/catsworld/qq-bot-api

// Params is the params of call api
type Params map[string]interface{}

// APIResponse is the response of calling API
// https://github.com/botuniverse/onebot-11/blob/master/communication/ws.md
type APIResponse struct {
	Status  string       `json:"status"`
	Data    gjson.Result `json:"data"`
	Msg     string       `json:"msg"`
	Wording string       `json:"wording"`
	RetCode int64        `json:"retcode"`
	Echo    uint64       `json:"echo"`
}

// APIRequest is the request sending to the cqhttp
// https://github.com/botuniverse/onebot-11/blob/master/communication/ws.md
type APIRequest struct {
	Action string `json:"action"`
	Params Params `json:"params"`
	Echo   uint64 `json:"echo"` // 该项不用填写，由Driver生成
}

// User is a user on QQ.
type User struct {
	// Private sender
	// https://github.com/botuniverse/onebot-11/blob/master/event/message.md#%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
	ID       int64  `json:"user_id"`
	TinyID   string `json:"tiny_id"` // TinyID 在 guild 下为 ID 的 string
	NickName string `json:"nickname"`
	Sex      string `json:"sex"` // "male"、"female"、"unknown"
	Age      int    `json:"age"`
	Area     string `json:"area"`
	// Group member
	// https://github.com/botuniverse/onebot-11/blob/master/event/message.md#%E7%BE%A4%E6%B6%88%E6%81%AF
	Card  string `json:"card"`
	Title string `json:"title"`
	Level string `json:"level"`
	Role  string `json:"role"` // "owner"、"admin"、"member"
	// Group anonymous
	AnonymousID   int64  `json:"anonymous_id" anonymous:"id"`
	AnonymousName string `json:"anonymous_name" anonymous:"name"`
	AnonymousFlag string `json:"anonymous_flag" anonymous:"flag"`
}

// File 文件
type File struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	BusID int64  `json:"busid"`
}

// Group 群
type Group struct {
	ID             int64  `json:"group_id"`
	Name           string `json:"group_name"`
	MemberCount    int64  `json:"member_count"`
	MaxMemberCount int64  `json:"max_member_count"`
}

// Name displays a simple text version of a user.
func (u *User) Name() string {
	if u.AnonymousName != "" {
		return u.AnonymousName
	}
	if u.Card != "" {
		return u.Card
	}
	if u.NickName != "" {
		return u.NickName
	}
	return strconv.FormatInt(u.ID, 10)
}

// String displays a simple text version of a user.
// It is normally a user's card, but falls back to a nickname as available.
func (u *User) String() string {
	p := ""
	if u.Title != "" {
		p = "[" + u.Title + "]"
	}
	return p + u.Name()
}

// H 是 Params 的简称
type H = Params
