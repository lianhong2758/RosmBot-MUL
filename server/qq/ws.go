package qq

import (
	"encoding/json"
	"errors"
	"strconv"
	
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

// Reset 恢复到 0 值
func (wp *WebsocketPayload) Reset() {
	*wp = WebsocketPayload{}
}

// GetHeartbeatInterval OpCodeHello 获得心跳周期 单位毫秒
func (wp *WebsocketPayload) GetHeartbeatInterval() (uint32, error) {
	if wp.Op != OpCodeHello {
		return 0, errors.New("[GetHeartbeatInterval]unexpected OpCode " + strconv.Itoa(int(wp.Op)) + ", T: " + wp.T + ", D: " + tool.BytesToString(wp.D))
	}
	data := &struct {
		H uint32 `json:"heartbeat_interval"`
	}{}
	err := json.Unmarshal(wp.D, data)
	return data.H, err
}

// SendPayload 发送 ws 包
func (c *Config) SendPayload(wp *WebsocketPayload) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(wp)
}

// WrapData 将结构体序列化到 wp.D
func (wp *WebsocketPayload) WrapData(v any) (err error) {
	wp.D, err = json.Marshal(v)
	log.Debugln("[ws]包装 Identify ", string(wp.D))
	return
}
