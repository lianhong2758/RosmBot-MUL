// Package score 签到
package score

import (
	"image"
	"math/rand"
	"os"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/imgfactory"
	"github.com/lianhong2758/RosmBot-MUL/kanban"
	"github.com/lianhong2758/RosmBot-MUL/message"
	walle "github.com/lianhong2758/RosmBot-MUL/plugins/public"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
)

const (
	backgroundURL = "https://iw233.cn/api.php?sort=pc"
	referer       = "https://weibo.com/"
	signinMax     = 1
	// SCOREMAX 分数上限定为1200
	SCOREMAX = 1200
)

var (
	defKey    = "2"
	rankArray = [...]int{0, 10, 20, 50, 100, 200, 350, 550, 750, 1000, 1200}
	drawmap   = map[string]func(a *scdata) (image.Image, error){
		"1": drawScore15,
		"2": drawScore16,
		"3": drawScore17,
		"4": drawScore18,
	}
)

func init() {
	en := rosm.Register(&rosm.PluginData{
		Name:       "签到",
		Help:       "- /签到\n- /获得签到背景",
		DataFolder: "score",
	})
	cachePath := en.DataFolder + "cache/"
	go func() {
		_ = os.RemoveAll(cachePath)
		err := os.MkdirAll(cachePath, 0755)
		if err != nil {
			panic(err)
		}
		sdb = initialize(en.DataFolder + "score.db")
	}()
	en.OnRex(`^/签到\s?(\d*)$`).Handle(func(ctx *rosm.Ctx) {
		// 选择key
		var key string = defKey
		if ctx.Being.ResultWord[1] != "" {
			key = ctx.Being.ResultWord[1]
		}
		drawfunc, ok := drawmap[key]
		if !ok {
			ctx.Send(message.Text("未找到签到设定:", key)) // 避免签到配置错误造成无图发送,但是已经签到的情况
			return
		}
		uid, name := ctx.Being.User.ID, ctx.Being.User.Name
		today := time.Now().Format("20060102")
		// 签到图片
		drawedFile := cachePath + uid + today + "signin.png"
		picFile := cachePath + uid + today + ".png"
		// 获取签到时间
		si := sdb.GetSignInByUID(uid)
		siUpdateTimeStr := si.UpdatedAt.Format("20060102")
		switch {
		case si.Count >= signinMax && siUpdateTimeStr == today:
			// 如果签到时间是今天
			ctx.Send(message.Reply(), message.Text("今天你已经签到过了！"))
			if file.IsExist(drawedFile) {
				ctx.Send(message.Image("file://" + file.BOTPATH + "/" + drawedFile))
			}
			return
		case siUpdateTimeStr != today:
			// 如果是跨天签到就清数据
			err := sdb.InsertOrUpdateSignInCountByUID(uid, 0)
			if err != nil {
				ctx.Send(message.Text("ERROR1: ", err))
				return
			}
		}
		// 更新签到次数
		err := sdb.InsertOrUpdateSignInCountByUID(uid, si.Count+1)
		if err != nil {
			ctx.Send(message.Text("ERROR2: ", err))
			return
		}
		// 更新经验
		level := sdb.GetScoreByUID(uid).Score + 1
		if level > SCOREMAX {
			level = SCOREMAX
			ctx.Send(message.AT(uid, name), message.Text("你的等级已经达到上限"))
		}
		err = sdb.InsertOrUpdateScoreByUID(uid, level)
		if err != nil {
			ctx.Send(message.Text("ERROR3: ", err))
			return
		}
		// 更新钱包
		rank := getrank(level)
		add := 1 + rand.Intn(10) + rank*5 // 等级越高获得的钱越高
		err = walle.InsertWalletOf(uid, add)
		if err != nil {
			ctx.Send(message.Text("ERROR4: ", err))
			return
		}
		alldata := scdata{
			userPic:    ctx.Bot.GetPortraitURI(ctx),
			drawedfile: drawedFile,
			picfile:    picFile,
			uid:        uid,
			nickname:   ctx.Being.User.Name,
			inc:        add,
			score:      walle.GetWalletOf(uid),
			level:      level,
			rank:       rank,
		}
		drawimage, err := drawfunc(&alldata)
		if err != nil {
			ctx.Send(message.Text("ERROR5: ", err))
			return
		}
		// done.
		f, err := os.Create(drawedFile)
		if err != nil {
			data, err := imgfactory.ToBytes(drawimage)
			if err != nil {
				ctx.Send(message.Text("ERROR6: ", err))
				return
			}
			ctx.Send(message.Reply(), message.ImageByte(data))
			return
		}
		_, err = imgfactory.WriteTo(drawimage, f)
		_ = f.Close()
		if err != nil {
			ctx.Send(message.Text("ERROR7: ", err))
			return
		}
		ctx.Send(message.Reply(), message.Image("file://"+file.BOTPATH+"/"+drawedFile))
	})

	en.OnWord("/获得签到背景").Handle(func(ctx *rosm.Ctx) {
		picFile := cachePath + ctx.Being.User.ID + time.Now().Format("20060102") + ".png"
		if file.IsNotExist(picFile) {
			ctx.Send(message.Reply(), message.Text("请先签到！"))
			return
		}
		ctx.Send(message.Image("file://" + file.BOTPATH + "/" + picFile))
	})
	/*en.AddWord("查看等级排名").SetBlock(true).
		Handle(func(ctx *c.CTX) {
			today := time.Now().Format("20060102")
			drawedFile := cachePath + today + "scoreRank.png"
			if file.IsExist(drawedFile) {
				ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
				return
			}
			st, err := sdb.GetScoreRankByTopN(10)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			if len(st) == 0 {
				ctx.Send(message.Text("ERROR: 目前还没有人签到过"))
				return
			}
			_, err = file.GetLazyData(text.FontFile, control.Md5File, true)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			b, err := os.ReadFile(text.FontFile)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			font, err := freetype.ParseFont(b)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			f, err := os.Create(drawedFile)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			var bars []chart.Value
			for _, v := range st {
				if v.Score != 0 {
					bars = append(bars, chart.Value{
						Label: ctx.CardOrNickName(v.UID),
						Value: float64(v.Score),
					})
				}
			}
			err = chart.BarChart{
				Font:  font,
				Title: "等级排名(1天只刷新1次)",
				Background: chart.Style{
					Padding: chart.Box{
						Top: 40,
					},
				},
				YAxis: chart.YAxis{
					Range: &chart.ContinuousRange{
						Min: 0,
						Max: math.Ceil(bars[0].Value/10) * 10,
					},
				},
				Height:   500,
				BarWidth: 50,
				Bars:     bars,
			}.Render(chart.PNG, f)
			_ = f.Close()
			if err != nil {
				_ = os.Remove(drawedFile)
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
		})
	engine.OnRegex(`^设置(默认)?签到预设\s?(\d*)$`, zero.SuperUserPermission).Limit(ctxext.LimitByUser).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if ctx.State["regex_matched"].([]string)[2] == "" {
			ctx.Send(message.Text("设置失败,数据为空"))
		} else {
			s := ctx.State["regex_matched"].([]string)[1]
			key := ctx.State["regex_matched"].([]string)[2]
			_, ok := drawmap[key]
			if !ok {
				ctx.Send(message.Text("未找到签到设定:", key)) // 避免签到配置错误
				return
			}
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			if s != "" {
				gid = defKeyID
			}
			err := ctx.State["manager"].(*ctrl.Control[*zero.Ctx]).Manager.SetExtra(gid, key)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			ctx.Send(message.Text("设置成功,当前", s, "预设为:", key))
		}
	})*/
}

func getHourWord(t time.Time) string {
	h := t.Hour()
	switch {
	case 6 <= h && h < 12:
		return "早上好"
	case 12 <= h && h < 14:
		return "中午好"
	case 14 <= h && h < 19:
		return "下午好"
	case 19 <= h && h < 24:
		return "晚上好"
	case 0 <= h && h < 6:
		return "凌晨好"
	default:
		return ""
	}
}

func getrank(count int) int {
	for k, v := range rankArray {
		if count == v {
			return k
		} else if count < v {
			return k - 1
		}
	}
	return -1
}

func initPic(picFile, userurl string) (avatar []byte, err error) {
	if file.IsExist(picFile) {
		return nil, nil
	}
	defer process.SleepAbout1sTo2s()
	url, err := web.GetRealURL(backgroundURL)
	if err != nil {
		return nil, err
	}
	data, err := web.RequestDataWith(web.NewDefaultClient(), url, "", referer, "", nil)
	if err != nil {
		return nil, err
	}
	if userurl != "" {
		avatar, err = web.GetData(userurl, "")
		if err != nil {
			return nil, err
		}
	} else {
		avatar, _ = os.ReadFile(kanban.Path)
	}
	return avatar, os.WriteFile(picFile, data, 0644)
}
