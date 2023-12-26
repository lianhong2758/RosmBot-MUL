package mys

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	log "github.com/sirupsen/logrus"
)

const (
	Host           = "https://bbs-api.miyoushe.com"
	URLGetRoomList = "/vila/api/bot/platform/getVillaGroupRoomList"
	URLGetUserData = "/vila/api/bot/platform/getMember"
	URLRecall      = "/vila/api/bot/platform/recallMessage"
	URLDeleteUser  = "/vila/api/bot/platform/deleteVillaMember"
	URLGetVilla    = "/vila/api/bot/platform/getVilla"
	URLPinMessage  = "/vila/api/bot/platform/pinMessage"
	URLGetWS       = "/vila/api/bot/platform/getWebsocketInfo"
)

// 获取房间列表
func GetRoomList(ctx *rosm.Ctx) (r *RoomList, err error) {
	data, err := web.Web(&http.Client{}, Host+URLGetRoomList, http.MethodGet, makeHeard(ctx), nil)
	log.Debugln("[GetRoomList]", string(data))
	if err != nil {
		return nil, err
	}
	r = new(RoomList)
	err = json.Unmarshal(data, r)
	return
}

// 获取用户信息
func GetUserData(ctx *rosm.Ctx, uid string) (r *UserData, err error) {
	data, _ := json.Marshal(H{"uid": uid})
	data, err = web.Web(&http.Client{}, Host+URLGetUserData, http.MethodGet, makeHeard(ctx), bytes.NewReader(data))
	log.Debugln("[GetUserData]", string(data))
	if err != nil {
		return nil, err
	}
	r = new(UserData)
	err = json.Unmarshal(data, r)
	return
}

// 获取用户名
func GetUserName(ctx *rosm.Ctx, uid string) string {
	r, err := GetUserData(ctx, uid)
	if err != nil {
		return ""
	}
	return r.Data.Member.Basic.Nickname
}

// 踢人
func DeleteUser(ctx *rosm.Ctx, uid string) (err error) {
	data, _ := json.Marshal(H{"uid": uid})
	data, err = web.Web(&http.Client{}, Host+URLDeleteUser, http.MethodPost, makeHeard(ctx), bytes.NewReader(data))
	log.Debugln("[DeleteUser]", string(data))
	var r ApiCode
	_ = json.Unmarshal(data, &r)
	if r.Retcode != 0 {
		return errors.New(r.Message)
	}
	return
}

// 撤回消息,消息id,发送时间,房间id
func Recall(ctx *rosm.Ctx, msgid string, msgtime any, roomid string) (err error) {
	var t int64
	switch tt := msgtime.(type) {
	case int64:
		t = tt
	case string:
		t = tool.StringToInt64(tt)
	}
	data, _ := json.Marshal(H{"msg_uid": msgid, "room_id": tool.StringToInt64(roomid), "msg_time": t})
	data, err = web.Web(&http.Client{}, Host+URLRecall, http.MethodPost, makeHeard(ctx), bytes.NewReader(data))
	log.Debugln("[Recall]", string(data))
	var r ApiCode
	_ = json.Unmarshal(data, &r)
	if r.Retcode != 0 {
		return errors.New(r.Message)
	}
	return
}

// 获取别野信息
func GetVillaData(ctx *rosm.Ctx) (r *VillaData, err error) {
	data, err := web.Web(&http.Client{}, Host+URLGetVilla, http.MethodGet, makeHeard(ctx), nil)
	log.Debugln("[GetVillaData]", string(data))
	if err != nil {
		return nil, err
	}
	r = new(VillaData)
	err = json.Unmarshal(data, r)
	return
}

// 置顶消息,消息id,发送时间,房间id,是否取消置顶
func PinMessage(ctx *rosm.Ctx, msgid, sendat any, roomid string, iscancel bool) (err error) {
	var t int64
	switch tt := sendat.(type) {
	case int64:
		t = tt
	case string:
		t = tool.StringToInt64(tt)
	}
	data, _ := json.Marshal(H{"msg_uid": msgid, "room_id": tool.StringToInt64(roomid), "send_at": t, "is_cancel": iscancel})
	data, err = web.Web(&http.Client{}, Host+URLPinMessage, http.MethodPost, makeHeard(ctx), bytes.NewReader(data))
	log.Debugln("[PinMessage]", string(data))
	var r ApiCode
	_ = json.Unmarshal(data, &r)
	if r.Retcode != 0 {
		return errors.New(r.Message)
	}
	return
}

type RoomList struct {
	ApiCode
	Data struct {
		List []struct {
			GroupID   string `json:"group_id"`
			GroupName string `json:"group_name"`
			RoomList  []struct {
				RoomID   string `json:"room_id"`
				RoomName string `json:"room_name"`
				RoomType string `json:"room_type"`
				GroupID  string `json:"group_id"`
			} `json:"room_list"`
		} `json:"list"`
	} `json:"data"`
}

type UserData struct {
	ApiCode
	Data struct {
		Member struct {
			Basic struct {
				UID       string `json:"uid"`
				Nickname  string `json:"nickname"`
				Introduce string `json:"introduce"`
				Avatar    string `json:"avatar"`
				AvatarURL string `json:"avatar_url"`
			} `json:"basic"`
			RoleIDList []string `json:"role_id_list"`
			JoinedAt   string   `json:"joined_at"`
			RoleList   []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Color    string `json:"color"`
				RoleType string `json:"role_type"`
				VillaID  string `json:"villa_id"`
				WebColor string `json:"web_color"`
			} `json:"role_list"`
		} `json:"member"`
	} `json:"data"`
}

type VillaData struct {
	ApiCode
	Data struct {
		Villa struct {
			VillaID        string   `json:"villa_id"`
			Name           string   `json:"name"`
			VillaAvatarURL string   `json:"villa_avatar_url"`
			OwnerUID       string   `json:"owner_uid"`
			IsOfficial     bool     `json:"is_official"`
			Introduce      string   `json:"introduce"`
			CategoryID     int      `json:"category_id"`
			Tags           []string `json:"tags"`
		} `json:"villa"`
	} `json:"data"`
}
