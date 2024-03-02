package fan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
)

func (c *Config) ListenHttp() {
	log.Infoln("[fan]开始尝试连接到网关:", UrlGetUpdates, ",BotID:", c.BotID)
	//准备数据结构
	d := &struct {
		Offset  int    `json:"offset,omitempty"`
		Limit   int    `json:"limit,omitempty"`
		Timeout int    `json:"timeout,omitempty"`
		MsgType string `json:"msg_type,omitempty"`
	}{} //{Offset: c.Offset, Limit: c.Limit, Timeout: c.Timeout, MsgType: c.MsgType}
	data, _ := json.Marshal(d)

	client := &http.Client{
		Timeout: 0, // 设置为0表示没有超时限制
	}
	fmt.Println(Host + fmt.Sprintf(UrlGetUpdates, c.user.Result.UserToken))
	req, err := http.NewRequest("POST", Host+fmt.Sprintf(UrlGetUpdates, c.user.Result.UserToken), bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json") // 设置请求头

	log.Infoln("[fan]连接到网关成功, 用户:", c.user.Result.ID)

	for {
		message := new(Message)
		resp, err := client.Do(req)
		if err != nil {
			log.Warn("[fan]Do Error:", err)
			resp.Body.Close()
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Warnf("[fan]Server returned non-OK status: %d", resp.StatusCode)
			resp.Body.Close()
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Warn("[fan]Read Error:", err)
			resp.Body.Close()
			continue
		}
		if len(body) == 0 {
			resp.Body.Close()
			continue
		} else {
			_ = json.Unmarshal(body, message)
			c.process(message)
		}
		// 发起下一次长轮询请求
		resp.Body.Close()
	}
}

func (c *Config) ListenWS() {
	log.Infoln("[fan]开始尝试连接到网关:", UrlGetUpdates, ",BotID:", c.BotID)

	for {
		conn, resp, err := websocket.DefaultDialer.Dial(fmt.Sprintf(UrlWs, c.user.Result.UserToken, "bot1", "eyJwbGF0Zm9ybSI6ICJib3QiLCAidmVyc2lvbiI6ICIxLjYuNjAiLCAiY2hhbm5lbCI6ICJvZmZpY2UiLCAiZGV2aWNlX2lkIjogInRlc3QiLCAiYnVpbGRfbnVtYmVyIjogIjEifQ=="), nil)
		if err != nil {
			log.Warnf("[sign]连接到网关 %v 时出现错误: %v", UrlWs, err)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		c.conn = conn
		_ = resp.Body.Close()
		break
	}
	log.Infoln("[sign]连接到网关成功, 用户:", c.user.Result.ID)
	log.Infoln("[ws]bot开始监听消息")
	for {
		req := new(WsReq)
		err := c.conn.ReadJSON(req)
		if err != nil {
			// 发送失败
			c.conn = nil
			log.Warnln("[ws]", c.BotName, "的网关连接断开, 尝试恢复:", err)
			return
		}
		log.Debugln("[ws]接收到数据:", req.Action)
		switch req.Action {
		case "pong": // Receive
			log.Debug("pong time:", req.Data)

		case "push":

		default:
		 
		}
	}

}
