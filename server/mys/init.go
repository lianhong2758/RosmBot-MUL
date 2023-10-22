package mys

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

var filePath = "config/mys.json"

func init() {
	data, err := os.ReadFile(filePath)
	if err != nil {
		MYSconfig.BotToken.BotPubKey = "-----BEGIN PUBLIC KEY----- abcabc123 -----END PUBLIC KEY----- "
		MYSconfig.EventPath = "/"
		MYSconfig.Port = "0.0.0.0:80"

		err = os.MkdirAll("config", os.ModePerm)
		if err != nil {
			log.Fatalln("无法创建 config 目录: ", err)
		}
		data, _ = json.MarshalIndent(MYSconfig, "", "  ")
		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			log.Fatalln("创建config失败: ", err)
		}
		log.Fatalln("创建初始化配置完成\n请填写config/mys.json文件后重新运行本程序\n字段解释:\ntoken:机器人基本信息:\neventpath:回调路径\nport:端口")
	}
	err = json.Unmarshal(data, &MYSconfig)
	if err != nil {
		log.Fatalln(err)
	}
	if MYSconfig.BotToken.BotID == "" || MYSconfig.BotToken.BotSecret == "" {
		log.Fatalln("[mys]未设置bot信息")
	}
	//备份
	MYSconfig.BotToken.BotSecretConst = MYSconfig.BotToken.BotSecret
	//修正
	var pubKeynext strings.Builder
	s := strings.Fields(MYSconfig.BotToken.BotPubKey)
	for k, v := range s {
		if k < 2 || k > len(s)-4 {
			pubKeynext.WriteString(v)
			pubKeynext.WriteString(" ")
		} else {
			pubKeynext.WriteString(v)
			pubKeynext.WriteString("\n")
		}
	}
	MYSconfig.BotToken.BotPubKey = strings.TrimSpace(pubKeynext.String()) + "\n"
	//加密验证
	MYSconfig.BotToken.BotSecret = Sha256HMac(MYSconfig.BotToken.BotPubKey, MYSconfig.BotToken.BotSecret)
}

// HMAC/SHA256加密
func Sha256HMac(pubKey string, botSecret string) string {
	h := hmac.New(sha256.New, []byte(pubKey))
	raw := []byte(botSecret)
	h.Write(raw)
	return hex.EncodeToString(h.Sum(nil))
}
