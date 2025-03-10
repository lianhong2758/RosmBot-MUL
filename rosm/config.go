package rosm

import (
	"encoding/json"
	"math/rand/v2"
	"os"

	"github.com/FloatTech/floatbox/file"
	"github.com/lianhong2758/RosmBot-MUL/kanban"
	"github.com/sirupsen/logrus"
)

// 插件的config
func LoadConfig(path string, v any) error {
	if file.IsNotExist(path) {
		//new
		data, err := json.MarshalIndent(v, "", "  ") //json.Marshal(v)
		if err != nil {
			return err
		}
		return os.WriteFile(path, data, 0666)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// 覆盖型创建插件config
func CreateConfig(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

// 保存插件的config
func (p *PluginData) SaveConfig() {
	CreateConfig(p.DataFolder+"config.json", p.Config)
}

// adapter config 路径./config
func LoadBotConfig(name string, v any) error {
	if file.IsNotExist("config/" + name) {
		//new
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		if file.IsNotExist("config") {
			if err = os.MkdirAll("config", os.ModePerm); err != nil {
				return err
			}
		}
		return os.WriteFile("config/"+name, data, 0666)
	}
	data, err := os.ReadFile("config/" + name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func GetRosmConfig() *RosmConfig {
	return config
}

var config = &RosmConfig{
	LogLevel:   logrus.InfoLevel,
	ApiTimeout: 30,
	BotName:    []string{"雪儿", "梦雪"},
	CmdStar:    []string{"/", ""},
}

func init() {
	LoadBotConfig("rosm.json", config)
	if kanban.GetArgIsDebug() {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("IN DEBUG MODE")
	} else {
		logrus.SetLevel(config.LogLevel)
	}
}

type RosmConfig struct {
	LogLevel   logrus.Level
	ApiTimeout int      //s
	BotName    []string //botName
	CmdStar    []string //命令的起始标记
}

func GetRandBotName() string {
	return config.BotName[rand.IntN(len(config.BotName))]
}
