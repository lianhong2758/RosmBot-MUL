package mys

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/RomiChan/websocket"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	log "github.com/sirupsen/logrus"
)

var botMap = map[string]*Config{}

// 默认值
type Token struct {
	*rosm.BotCard
	BotSecret      string `json:"bot_secret"`
	BotPubKey      string `json:"bot_pub_key"`
	BotSecretConst string `json:"-"`
}
type Config struct {
	Protocol  string `json:"protocol"` //协议:ws/http
	BotToken  Token  `json:"token"`
	TestVilla string `json:"test_villa,omitempty"` //测试别野id
	EventPath string `json:"eventpath,omitempty"`  //路径
	Port      string `json:"port,omitempty"`       //端口

	wr     *WebsocketInfoResp //获取的WebsocketInfoResp
	conn   *websocket.Conn    // conn 目前的 wss 连接
	hbonce sync.Once          // hbonce 保证仅执行一次 heartbeat
}

func (c *Config) Card() *rosm.BotCard {
	return c.BotToken.BotCard
}
func NewConfig(path string) (c *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		c = new(Config)
		c.BotToken.BotCard = new(rosm.BotCard)
		//初始配置
		c.BotToken.Master = []string{"123456"}
		c.BotToken.BotPubKey = "-----BEGIN PUBLIC KEY----- abcabc123 -----END PUBLIC KEY----- "
		c.BotToken.BotID = "bot_..."

		//#######################################
		//输入
		fmt.Println("请输入选择的连接方式:\n0:http连接\n1:ws正向连接")
		for {
			var t int
			fmt.Scanln(&t)
			switch t {
			case 0:
				c.Protocol = "http"
				c.EventPath = "/rosmbot"
				c.Port = "0.0.0.0:10001"
			case 1:
				c.Protocol = "ws"
				c.TestVilla = "463"
			default:
				fmt.Println("输入错误!重新输入:")
				continue
			}
			break
		}
		//#######################################
		err = os.MkdirAll("config", os.ModePerm)
		if err != nil {
			log.Fatalln("[qq]无法创建 config 目录: ", err)
		}
		data, _ = json.MarshalIndent(c, "", "  ")
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			log.Fatalln("[mys]创建config失败: ", err)
		}
		log.Infoln("创建初始化配置完成\n请填写config/mys.json文件后重新运行本程序\n字段解释:\ntoken:机器人基本信息:\neventpath:回调路径\nport:端口\ntest_villa:测试别野号")
		os.Exit(0)
	}
	c = new(Config)
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Fatalln(err)
	}
	if c.BotToken.BotID == "" || c.BotToken.BotSecret == "" {
		log.Fatalln("[mys]未设置bot信息")
	}
	//备份
	c.BotToken.BotSecretConst = c.BotToken.BotSecret
	//修正
	var pubKeynext strings.Builder
	s := strings.Fields(c.BotToken.BotPubKey)
	for k, v := range s {
		if k < 2 || k > len(s)-4 {
			pubKeynext.WriteString(v)
			pubKeynext.WriteString(" ")
		} else {
			pubKeynext.WriteString(v)
			pubKeynext.WriteString("\n")
		}
	}
	c.BotToken.BotPubKey = strings.TrimSpace(pubKeynext.String()) + "\n"
	//加密验证
	c.BotToken.BotSecret = Sha256HMac(c.BotToken.BotPubKey, c.BotToken.BotSecret)
	return
}

// HMAC/SHA256加密
func Sha256HMac(pubKey string, botSecret string) string {
	h := hmac.New(sha256.New, []byte(pubKey))
	raw := []byte(botSecret)
	h.Write(raw)
	return hex.EncodeToString(h.Sum(nil))
}
