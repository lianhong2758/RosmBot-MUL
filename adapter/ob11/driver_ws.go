package ob11

import (
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RomiChan/websocket"

	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var (
	nullResponse = APIResponse{}
)

// WSClient ...
type WSClient struct {
	seq         uint64
	conn        *websocket.Conn
	mu          sync.Mutex // 写锁
	seqMap      seqSyncMap
	Url         string // ws连接地址
	AccessToken string
	selfID      int64
}

// NewWebSocketClient 默认Driver，使用正向WS通信
func NewWebSocketClient(url, accessToken string) *WSClient {
	return &WSClient{
		Url:         url,
		AccessToken: accessToken,
	}
}

// Connect 连接ws服务端
func (ws *WSClient) Connect(config *Config) {
	log.Infof("[ws] 开始尝试连接到Websocket服务器: %v", ws.Url)
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"User-Agent":    []string{"RosmBot/1.1.0"},
	}
	if ws.AccessToken != "" {
		header["Authorization"] = []string{"Bearer " + ws.AccessToken}
	}

	network, address := resolveURI(ws.Url)
	dialer := websocket.Dialer{
		NetDial: func(_, addr string) (net.Conn, error) {
			if network == "unix" {
				host, _, err := net.SplitHostPort(addr)
				if err != nil {
					host = addr
				}
				filepath, err := base64.RawURLEncoding.DecodeString(host)
				if err == nil {
					addr = tool.BytesToString(filepath)
				}
			}
			return net.Dial(network, addr) // support unix socket transport
		},
	}

	for {
		conn, res, err := dialer.Dial(address, header)
		if err != nil {
			log.Warnf("[ws] 连接到Websocket服务器 %v 时出现错误: %v", ws.Url, err)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		ws.conn = conn
		_ = res.Body.Close()
		var rsp struct {
			SelfID int64 `json:"self_id"`
		}
		err = ws.conn.ReadJSON(&rsp)
		if err != nil {
			log.Warnf("[ws] 与Websocket服务器 %v 握手时出现错误: %v", ws.Url, err)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		ws.selfID = rsp.SelfID
		config.BotID = strconv.FormatInt(ws.selfID, 10)
		botMap[config.BotID] = config
		//zero.APICallers.Store(ws.selfID, ws) // 添加Caller到 APICaller list...
		log.Infof("[ws] 连接Websocket服务器: %s 成功, 账号: %d", ws.Url, rsp.SelfID)
		break
	}
}

// Listen 开始监听事件
func (ws *WSClient) Listen(config *Config) {
	for {
		t, payload, err := ws.conn.ReadMessage()
		if err != nil { // reconnect
			delete(botMap, strconv.FormatInt(ws.selfID, 10))
			//	zero.APICallers.Delete(ws.selfID) // 断开从apicaller中删除
			log.Warn("[ws] Websocket服务器连接断开...")
			time.Sleep(time.Millisecond * time.Duration(3))
			ws.Connect(config)
			continue
		}
		if t != websocket.TextMessage {
			continue
		}
		rsp := gjson.Parse(tool.BytesToString(payload))
		if rsp.Get("echo").Exists() { // 存在echo字段，是api调用的返回
			log.Debug("[ws] 接收到API调用返回: ", strings.TrimSpace(tool.BytesToString(payload)))
			if c, ok := ws.seqMap.LoadAndDelete(rsp.Get("echo").Uint()); ok {
				c <- APIResponse{ // 发送api调用响应
					Status:  rsp.Get("status").String(),
					Data:    rsp.Get("data"),
					Msg:     rsp.Get("msg").Str,
					Wording: rsp.Get("wording").Str,
					RetCode: rsp.Get("retcode").Int(),
					Echo:    rsp.Get("echo").Uint(),
				}
				close(c) // channel only use once
			}
			continue
		}
		if rsp.Get("meta_event_type").Str == "heartbeat" { // 忽略心跳事件
			continue
		}
		log.Debug("[ws] 接收到事件: ", tool.BytesToString(payload))
		config.processEvent(payload, ws)
	}
}

func (ws *WSClient) nextSeq() uint64 {
	return atomic.AddUint64(&ws.seq, 1)
}

// CallApi 发送ws请求
func (ws *WSClient) CallApi(req APIRequest) (APIResponse, error) {
	ch := make(chan APIResponse, 1)
	req.Echo = ws.nextSeq()
	ws.seqMap.Store(req.Echo, ch)

	// send message
	ws.mu.Lock() // websocket write is not goroutine safe
	err := ws.conn.WriteJSON(&req)
	ws.mu.Unlock()
	if err != nil {
		log.Warn("[ws] 向WebsocketServer发送API请求失败: ", err.Error())
		return nullResponse, err
	}
	log.Debug("[ws] 向服务器发送请求: ", &req)

	select { // 等待数据返回
	case rsp, ok := <-ch:
		if !ok {
			return nullResponse, io.ErrClosedPipe
		}
		return rsp, nil
	case <-time.After(time.Minute):
		return nullResponse, os.ErrDeadlineExceeded
	}
}
