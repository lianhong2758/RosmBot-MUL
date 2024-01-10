package ob11

import (
	"time"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
)

// Run 主函数初始化
func (c *Config) Run() {
	runinit(c)
	for _, driver := range c.Driver {
		driver.Connect()
		driver.Listen(process)
		c.mul(tool.Int64ToString(driver.GetSelfID()))
	}
}

func runinit(c *Config) {
	if c.MaxProcessTime == 0 {
		c.MaxProcessTime = time.Minute * 4
	}
}
func (c *Config) mul(id string) {
	rosm.MULChan <- rosm.MUL{Types: "ob11", Name: c.BotName, BotID: id}
}
