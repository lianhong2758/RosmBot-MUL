package mys

import "go.uber.org/atomic"

var (
	counter            = atomic.NewUint64(1)
	magic       uint32 = 0xBABEFACE
	headerLenV2 uint32 = 24
	headerLenV1 uint32 = 20

	// FlagRequest ...
	FlagRequest uint32 = 1
	// FlagResponse ...
	FlagResponse uint32 = 2
)

type Message struct {
	// 定长头, 二进制
	Magic   uint32 // 用于标识报文的开始
	DataLen uint32 // 消息体长度, 编码后的数据长度，即包括消息头和body

	// 变长头,二进制
	HeaderLen uint32 // 消息头长度, 解码后的头长度
	ID        uint64 // 消息的序号，请求包的ID单调递增，回应包的ID与请求包的ID一致
	Flag      uint32 // 消息的方向: 1-请求；2-回应
	Type      uint32 // 消息的类型,相当于命令字
	AppId     int32  // 消息发出方所属的AppId

	// Body json格式
	Body []byte // 消息的内容
}

// NewRequestMsg 创建一个Message消息包
func NewRequestMsg(typ uint32, id uint64, appId int32, data []byte) *Message {
	return &Message{
		Magic: magic,
		ID:    id,
		Flag:  FlagRequest,
		Type:  typ,
		AppId: appId,
		Body:  data,
	}
}

// NewResponseMsg 创建一个Message消息包
func NewResponseMsg(typ uint32, id uint64, appId int32, data []byte) *Message {
	return &Message{
		Magic: magic,
		ID:    id,
		Flag:  FlagResponse,
		Type:  typ,
		AppId: appId,
		Body:  data,
	}
}

// GetFixHeaderLen 定长头长度
func (msg Message) GetFixHeaderLen() uint32 {
	return 8
}

// GetHeaderLen 变长头长度
func (msg Message) GetHeaderLen() uint32 {
	if msg.AppId == 0 {
		return headerLenV1
	}
	return headerLenV2
}

// GetAppId ...
func (msg Message) GetAppId() int32 {
	return msg.AppId
}

// GetType 获取消息数据段长度
func (msg *Message) GetType() uint32 {
	return msg.Type
}

// GetMagic ...
func (msg *Message) GetMagic() uint32 {
	return msg.Magic
}

// GetID ...
func (msg *Message) GetID() uint64 {
	return msg.ID
}

// GetFlag ...
func (msg *Message) GetFlag() uint32 {
	return msg.Flag
}

// CheckMagic ...
func (msg *Message) CheckMagic() bool {
	return msg.Magic == magic
}

// GetDataLen ...
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

// GetBody 获取消息内容
func (msg *Message) GetBody() []byte {
	return msg.Body
}

// SetDataLen 设置消息数据段长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

// SetHeaderLen 设置消息数据段长度
func (msg *Message) SetHeaderLen(len uint32) {
	msg.HeaderLen = len
}

// SetBody 设计消息内容
func (msg *Message) SetBody(data []byte) {
	msg.Body = data
}

// UniqMsgID ...
func UniqMsgID() uint64 {
	return counter.Add(1)
}
