package mys

import (
	vila_bot "github.com/lianhong2758/RosmBot-MUL/server/mys/proto"
)

type H = map[string]any

// 回调的请求结构

type InfoSTR struct {
	Event struct {
		Robot      *vila_bot.Robot               `protobuf:"bytes,1,opt,name=robot,proto3" json:"robot,omitempty"`                                   // 事件相关机器人
		Type       vila_bot.RobotEvent_EventType `protobuf:"varint,2,opt,name=type,proto3,enum=vila_bot.RobotEvent_EventType" json:"type,omitempty"` // 事件类型
		CreatedAt  int64                         `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`         // 事件发生时间
		Id         string                        `protobuf:"bytes,5,opt,name=id,proto3" json:"id,omitempty"`                                         // 事件 id
		SendAt     int64                         `protobuf:"varint,6,opt,name=send_at,json=sendAt,proto3" json:"send_at,omitempty"`                  // 事件消息投递时间
		ExtendData struct {                      // 包含事件的具体数据
			EventData EventData `json:"EventData"`
		} `protobuf:"bytes,3,opt,name=extend_data,json=extendData,proto3" json:"extend_data,omitempty"` // 事件拓展信息
	} `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
}

// 所有事件
type EventData struct {
	vila_bot.RobotEvent_ExtendData_JoinVilla
	vila_bot.RobotEvent_ExtendData_SendMessage
	vila_bot.RobotEvent_ExtendData_CreateRobot
	vila_bot.RobotEvent_ExtendData_DeleteRobot
	vila_bot.RobotEvent_ExtendData_AddQuickEmoticon
	vila_bot.RobotEvent_ExtendData_AuditCallback
	vila_bot.RobotEvent_ExtendData_ClickMsgComponent
}

// 接收的原始消息,解析
type MessageContent struct {
	Trace struct {
		VisualRoomVersion string `json:"visual_room_version"`
		AppVersion        string `json:"app_version"`
		ActionType        int    `json:"action_type"`
		BotMsgID          string `json:"bot_msg_id"`
		Client            string `json:"client"`
		Env               string `json:"env"`
		RongSdkVersion    string `json:"rong_sdk_version"`
	} `json:"trace"`
	Quote struct {
		QuotedMessageSendTime   int64  `json:"quoted_message_send_time"`
		OriginalMessageID       string `json:"original_message_id"`
		OriginalMessageSendTime int64  `json:"original_message_send_time"`
		QuotedMessageID         string `json:"quoted_message_id"`
	} `json:"quote"`
	MentionedInfo struct {
		MentionedContent string   `json:"mentionedContent"`
		UserIDList       []string `json:"userIdList"`
		Type             int      `json:"type"`
	} `json:"mentionedInfo"`
	User    user    `json:"user"`
	Content content `json:"content"`
}

type content struct {
	Images   []any `json:"images"`
	Entities []struct {
		Offset int `json:"offset"`
		Length int `json:"length"`
		Entity struct {
			Type  string `json:"type"`
			BotID string `json:"bot_id"`
		} `json:"entity"`
	} `json:"entities"`
	Text string `json:"text"`
}
type user struct {
	PortraitURI string `json:"portraitUri"`
	Extra       string `json:"extra"`
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	ID          string `json:"id"`
	Portrait    string `json:"portrait"`
}

// 消息发送回调
type SendState struct {
	ApiCode
	Data struct {
		BotMsgID string `json:"bot_msg_id"`
	} `json:"data"`
}

// api返回
type ApiCode struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
}
