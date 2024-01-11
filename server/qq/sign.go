package qq

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/RomiChan/websocket"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	log "github.com/sirupsen/logrus"
)

var clientConst = &http.Client{}

// 获取接口凭证
//
// https://docs.qq.com/doc/DRkVHT1N2a1JYSnVr
func AccessToken(appId, appSecret string) (r *AccessTokenStr, err error) {
	data, _ := json.Marshal(H{"appId": appId, "clientSecret": appSecret})
	data, err = web.Web(clientConst, urlAccessToken, http.MethodPost, func(req *http.Request) { req.Header.Add("Content-Type", "application/json") }, bytes.NewReader(data))
	log.Debugln("[AccessToken]", string(data))
	if err != nil {
		return nil, err
	}
	r = new(AccessTokenStr)
	err = json.Unmarshal(data, r)
	return
}

func (c *Config) getAccessToken() {
	/*var (
		r   = new(AccessTokenStr)
		err error
	)
	go func() {
		for {
			r, err = AccessToken(c.BotToken.AppId, c.BotToken.AppSecret)
			if err != nil {
				log.Debug("[sign]AccessToken", *r)
				log.Error("[sign]获取AccessToken失败,正在尝试重新获取...")
				time.Sleep(time.Second * 3)
				continue
			}
			c.access = r.AccessToken
			time.Sleep(time.Second * 7150) //7200过期
		}
	}()*/
	c.access = c.Authorization()
}

// getinitinfo 获得 gateway 和 shard
func (c *Config) getinitinfo() (err error) {
	var (
		shard [2]byte
		gw    string
	)
	shard[1] = 1
	if c.ShardIndex == 0 {
		gw, err = c.GetGeneralWSSGateway()
		if err != nil {
			return
		}
	} else {
		var sgw *ShardWSSGateway
		sgw, err = c.GetShardWSSGateway()
		if err != nil {
			return
		}
		if sgw.Shards <= int(c.ShardIndex) {
			err = errors.New("shard index " + strconv.Itoa(int(c.ShardIndex)) + " >= suggested size " + strconv.Itoa(sgw.Shards))
			return
		}
		gw = sgw.URL
		shard[0] = byte(c.ShardIndex)
		shard[1] = byte(sgw.Shards)
	}
	c.gateway = gw
	c.shard = shard
	return
}

// GetGeneralWSSGateway 获取通用 WSS 接入点
//
// https://bot.q.qq.com/wiki/develop/api/openapi/wss/url_get.html
func (c *Config) GetGeneralWSSGateway() (string, error) {
	resp := struct {
		CodeMessageBase
		U string `json:"url"`
	}{}
	err := c.GetOpenAPI(urlGetway, nil, &resp)
	return resp.U, err
}

// GetShardWSSGateway 获取带分片 WSS 接入点
func (c *Config) GetShardWSSGateway() (*ShardWSSGateway, error) {
	resp := struct {
		CodeMessageBase
		ShardWSSGateway
	}{}
	err := c.GetOpenAPI(urlGetwayWss, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.ShardWSSGateway, err
}

// Connect 连接到 Gateway + 鉴权连接
//
// https://bot.q.qq.com/wiki/develop/api/gateway/reference.html#_1-%E8%BF%9E%E6%8E%A5%E5%88%B0-gateway
func (c *Config) Connect() {
	network, address := resolveURI(c.gateway)
	log.Infoln("[sign]开始尝试连接到网关:", address, ", AppID:", c.BotToken.AppId)
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
		conn, resp, err := dialer.Dial(address, http.Header{})
		if err != nil {
			log.Warnf("[sign]连接到网关 %v 时出现错误: %v", c.gateway, err)
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		c.conn = conn
		_ = resp.Body.Close()
		payload, err := c.reveive()
		if err != nil {
			log.Warnln("[sign]获取心跳间隔时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		hb, err := payload.GetHeartbeatInterval()
		if err != nil {
			log.Warnln("[sign]解析心跳间隔时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		payload.Op = OpCodeIdentify
		err = payload.WrapData(&OpCodeIdentifyMessage{
			Token:      c.Authorization(),
			Intents:    c.Intents,
			Shard:      c.shard,
			Properties: nil,
		})
		if err != nil {
			log.Warnln("[sign]包装 Identify 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		err = c.SendPayload(&payload)
		if err != nil {
			log.Warnln("[sign]发送 Identify 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		payload, err = c.reveive()
		if err != nil {
			log.Warnln("[sign]获取 EventReady 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(2 * time.Second) // 等待两秒后重新连接
			continue
		}
		c.ready, err = payload.GetEventReady()
		if err != nil {
			log.Warnln("[sign]解析 EventReady 时出现错误:", err)
			_ = conn.Close()
			time.Sleep(3 * time.Second) // 等待3秒后重新连接
			continue
		}
		atomic.StoreUint32(&c.heartbeat, hb)
		break
	}
	log.Infoln("[sign]连接到网关成功, 用户名:", c.ready.User.Username)
	c.hbonce.Do(func() {
		go c.doheartbeat()
	})
}

// receive 收一个 payload
func (c *Config) reveive() (payload WebsocketPayload, err error) {
	err = c.conn.ReadJSON(&payload)
	return
}

// doheartbeat 按指定间隔进行心跳包发送
func (c *Config) doheartbeat() {
	payload := struct {
		Op OpCode  `json:"op"`
		D  *uint32 `json:"d"`
	}{Op: OpCodeHeartbeat}
	for {
		if atomic.LoadUint32(&c.heartbeat) == 0 {
			time.Sleep(time.Second)
			log.Warnln("[heart]等待服务器建立连接...")
			continue
		}
		time.Sleep(time.Duration(c.heartbeat) * time.Millisecond)
		if c.seq == 0 {
			payload.D = nil
		} else {
			payload.D = &c.seq
		}
		c.mu.Lock()
		err := c.conn.WriteJSON(&payload)
		c.mu.Unlock()
		if err != nil {
			log.Warnln("[heart]发送心跳时出现错误:", err)
		}
	}
}

// Listen 监听事件
func (c *Config) Listen() {
	log.Infoln("[qq-ws]开始监听", c.ready.User.Username, "的事件")
	payload := WebsocketPayload{}
	lastheartbeat := time.Now()
	for {
		payload.Reset()
		err := c.conn.ReadJSON(&payload)
		if err != nil {
			// 发送失败
			c.conn = nil
			atomic.StoreUint32(&c.heartbeat, 0)
			log.Warnln("[ws]", c.ready.User.Username, "的网关连接断开, 尝试恢复:", err)
			for {
				time.Sleep(time.Second)
				//开始重连
				err = c.Resume()
				if err == nil {
					break
				}
				log.Warnln("[ws]", c.ready.User.Username, "的网关连接恢复失败:", err)
			}
			continue
		}
		log.Debugln("[ws]接收到第", payload.S, "个事件:", payload.Op, ", 类型:", payload.T, ", 数据:", tool.BytesToString(payload.D))
		c.seq = payload.S
		//接收的类型进行选择
		switch payload.Op {
		case OpCodeDispatch: // Receive
			switch payload.T {
			case "RESUMED":
				log.Infoln("ws", c.ready.User.Username, "的网关连接恢复完成")
			default:
				c.process(&payload)
			}
		case OpCodeHeartbeat: // Send/Receive
			log.Debugln("[ws]收到服务端推送心跳, 间隔:", time.Since(lastheartbeat))
			lastheartbeat = time.Now()
		case OpCodeReconnect: // Receive
			log.Warnln("[ws]收到服务端通知重连")
			atomic.StoreUint32(&c.heartbeat, 0)
			c.Connect()
		case OpCodeInvalidSession: // Receive
			log.Warnln("[ws]", c.ready.User.Username, "的网关连接恢复失败: InvalidSession, 尝试重连...")
			atomic.StoreUint32(&c.heartbeat, 0)
			c.Connect()
		case OpCodeHello: // Receive
			intv, err := payload.GetHeartbeatInterval()
			if err != nil {
				log.Warnln("[ws]解析心跳间隔时出现错误:", err)
				continue
			}
			atomic.StoreUint32(&c.heartbeat, intv)
		case OpCodeHeartbeatACK: // Receive/Reply
			log.Debugln("[ws]收到心跳返回, 间隔:", time.Since(lastheartbeat))
			lastheartbeat = time.Now()
		case OpCodeHTTPCallbackACK: // Reply
		default:
			log.Warnln("[ws]未知事件, 序号:", payload.S, ", Op:", payload.Op, ", 类型:", payload.T, ", 数据:", tool.BytesToString(payload.D))
		}
	}
}

// Resume 恢复连接
//
// https://. .q.qq.com/wiki/develop/api/gateway/reference.html#_4-%E6%81%A2%E5%A4%8D%E8%BF%9E%E6%8E%A5
func (c *Config) Resume() error {
	network, address := resolveURI(c.gateway)
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
	conn, resp, err := dialer.Dial(address, http.Header{})
	if err != nil {
		return err
	}
	c.conn = conn
	_ = resp.Body.Close()
	payload := WebsocketPayload{Op: OpCodeResume}
	payload.WrapData(&struct {
		T string `json:"token"`
		S string `json:"session_id"`
		Q uint32 `json:"seq"`
	}{c.Authorization(), c.ready.SessionID, c.seq})
	return c.SendPayload(&payload)
}

// GetEventReady OpCodeDispatch READY
func (wp *WebsocketPayload) GetEventReady() (er EventReady, err error) {
	if wp.Op != OpCodeDispatch {
		err = errors.New("[opencode] unexpected OpCode " + strconv.Itoa(int(wp.Op)) + ", T: " + wp.T + ", D: " + tool.BytesToString(wp.D))
		return
	}
	if wp.T != "READY" {
		err = errors.New("[opencode] unexpected event type " + wp.T + ", T: " + wp.T + ", D: " + tool.BytesToString(wp.D))
		return
	}
	err = json.Unmarshal(wp.D, &er)
	return
}
