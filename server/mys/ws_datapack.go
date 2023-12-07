package mys

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/RomiChan/websocket"
)

// DataPack 封包拆包类实例
type DataPack struct {
	conn *websocket.Conn
	pkg  []byte
	err  error
}

// NewDataPack 封包拆包实例初始化方法
func NewDataPack(conn *websocket.Conn) *DataPack {
	return &DataPack{conn: conn}
}

// GetFixHeaderLen 获取包定长头长度方法
func (dp *DataPack) GetFixHeaderLen() uint32 {
	return Message{}.GetFixHeaderLen()
}

// Pack 封包方法(压缩数据)
func (dp *DataPack) Pack(msg *Message) *DataPack {
	dp.err = nil
	data, err := dp.packData(msg)
	if err != nil {
		dp.err = err
		return dp
	}
	fixHead, err := dp.packFixHeader(msg, len(data))
	if err != nil {
		dp.err = err
		return dp
	}
	dp.pkg = append(fixHead, data...)
	return dp
}

// 发送
func (dp *DataPack) Send() error {
	if dp.err != nil {
		return dp.err
	}
	return dp.conn.WriteMessage(websocket.BinaryMessage, dp.pkg)
}

// 写定长头数据
func (dp *DataPack) packFixHeader(msg *Message, dataLen int) ([]byte, error) {
	headerBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(headerBuff, binary.LittleEndian, msg.GetMagic()); err != nil {
		return nil, err
	}
	if err := binary.Write(headerBuff, binary.LittleEndian, uint32(dataLen)); err != nil {
		return nil, err
	}
	return headerBuff.Bytes(), nil
}

func (dp *DataPack) packData(msg *Message) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetHeaderLen()); err != nil {
		return nil, err
	}

	// 写Header数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetID()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetFlag()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetType()); err != nil {
		return nil, err
	}
	if msg.GetHeaderLen() == headerLenV2 { // 兼容v1版本json编码的协议，只有headerLen超过20时才填充Appid
		if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetAppId()); err != nil {
			return nil, err
		}
	}

	// 写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetBody()); err != nil {
		return nil, err
	}

	out := dataBuff.Bytes()

	return out, nil
}

// 读取下一条信息并解析
func (dp *DataPack) UnpackWs() (*Message, error) {
	_, payload, err := dp.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	msg := &Message{}
	// 解析定长头
	if err := dp.unpackFixHeader(payload[:dp.GetFixHeaderLen()], msg); err != nil {
		return nil, err
	}
	// 解析不定长正文
	payLoadLen := dp.GetFixHeaderLen() + msg.GetDataLen()
	if int(payLoadLen) != len(payload) {
		return nil, fmt.Errorf("websocket invlaid payload len expectLen: %d realLen:%d", payLoadLen, len(payload))
	}
	data := payload[dp.GetFixHeaderLen():payLoadLen]
	return dp.unpackBody(data, msg)
}

func (dp *DataPack) unpackFixHeader(fixHeadData []byte, msg *Message) error {
	dataBuff := bytes.NewReader(fixHeadData)

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Magic); err != nil {
		return err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return err
	}

	if !msg.CheckMagic() {
		return fmt.Errorf("wrong magic in message: %X", msg.Magic)
	}

	return nil
}

func (dp *DataPack) unpackBody(bodyData []byte, msg *Message) (*Message, error) {
	decodedData := bodyData

	dataBuff := bytes.NewReader(decodedData)

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.HeaderLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Flag); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Type); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.AppId); err != nil {
		return nil, err
	}

	if msg.HeaderLen > uint32(len(decodedData)) {
		return nil, fmt.Errorf("too large header len: %d(%X) > %d", msg.HeaderLen, msg.HeaderLen, len(decodedData))
	}

	body := decodedData[msg.HeaderLen:]
	msg.SetBody(body)
	return msg, nil
}
