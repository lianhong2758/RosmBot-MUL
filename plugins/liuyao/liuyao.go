package myplugin

import (
	"bufio"
	"bytes"
	"os"

	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/plugins/liuyao/yao"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	log "github.com/sirupsen/logrus"
)

var f *os.File

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name:       "六爻",
		Help:       "- /六爻春风",
		DataFolder: "liuyao",
	})
	{
		var err error
		f, err = os.Open(en.DataFolder + "ans.txt")
		if err != nil {
			log.Error("[liuyao]未找到解卦文件...")
		}
	}
	en.AddRex(`/六爻春风\s*(\d*)?$`).Handle(func(ctx *rosm.Ctx) {
		var buf bytes.Buffer
		seed := uint64(0)
		if ctx.Being.Rex[1] != "" {
			seed = uint64(tool.StringToInt64(ctx.Being.Rex[1]))
		} else {
			seed = uint64(tool.StringToInt64(ctx.Being.User.ID))
		}
		y := yao.NewGua(seed)
		for i := range 6 {
			y.Divination(i)
		}
		y.Changes()
		l := yao.GetItemOFGua(y.GetManifestation())
		f.Seek(0, 0)
		scanner := bufio.NewScanner(f)
		lineNumber := 0
		for scanner.Scan() {
			if lineNumber >= l && lineNumber <= l+2 {
				buf.WriteByte('\n')
				buf.WriteString(scanner.Text())
			}
			lineNumber++
		}
		if err := scanner.Err(); err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}

		ctx.Send(message.Text("你本次六爻占卜如下:\n本卦为:", y.GetInward(), "\n变卦为:", y.GetRange(), "\n", buf.String()))
	})
}
