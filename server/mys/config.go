package mys

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"

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
	BotToken  Token  `json:"token"`
	EventPath string `json:"eventpath,omitempty"`
	Port      string `json:"port,omitempty"`
}

func (c *Config) Card() *rosm.BotCard {
	return c.BotToken.BotCard
}
func NewConfig(path string) (c *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		c = new(Config)
		//初始配置
		c.BotToken.Master = []string{"123456"}
		c.BotToken.BotPubKey = "-----BEGIN PUBLIC KEY----- abcabc123 -----END PUBLIC KEY----- "
		c.EventPath = "/"
		c.Port = "0.0.0.0:80"
		c.BotToken.BotID = "bot_..."

		err = os.MkdirAll("config", os.ModePerm)
		if err != nil {
			log.Fatalln("[qq]无法创建 config 目录: ", err)
		}
		data, _ = json.MarshalIndent(c, "", "  ")
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			log.Fatalln("[mys]创建config失败: ", err)
		}
		log.Fatalln("创建初始化配置完成\n请填写config/mys.json文件后重新运行本程序\n字段解释:\ntoken:机器人基本信息:\neventpath:回调路径\nport:端口")
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
