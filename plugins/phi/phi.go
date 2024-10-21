package phi

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/file"
	"github.com/lianhong2758/PhigrosAPI/phigros"
	"github.com/lianhong2758/RosmBot-MUL/message"
	"github.com/lianhong2758/RosmBot-MUL/plugins/phi/qr"
	font "github.com/lianhong2758/RosmBot-MUL/plugins/public"
	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	"github.com/skip2/go-qrcode"
)

var challengemoderank = []string{"white", "green", "blue", "red", "gold", "rainbow"}
var fontsd []byte

func init() {
	//	en := rosm.Register(rosm.NewRegist("phi", "- /bind you_Session\n- /b19", "phi"))
	en := rosm.Register(&rosm.PluginData{
		Name: "phi",
		Help: "- /bind you_Session\n" +
			"- /b19\n" +
			"- /phi扫码登录\n" +
			"- /phi帮助",
		DataFolder: "phi",
	})
	// 课题模式图标
	Challengemode = en.DataFolder + "challengemode/"
	// 字体
	Font = font.MaokenFontFile
	// 评级
	Rank = en.DataFolder + "rank/"
	// 曲绘
	Illustration = en.DataFolder + "illustration/"
	//战绩图
	output = en.DataFolder + "output/"
	//头图
	avatar = en.DataFolder + "avatar/"
	if file.IsNotExist(en.DataFolder + "sessions") {
		_ = os.MkdirAll(en.DataFolder+"sessions", 0755)
	}

	//font
	fontsd, _ = os.ReadFile(Font)
	err := phigros.LoadDifficult(en.DataFolder + "difficulty.tsv")
	if err != nil {
		panic(err)
	}
	en.AddRex(`^/b19$`).Handle(func(ctx *rosm.Ctx) {
		Session := FindSessionFID(en.DataFolder + "sessions/" + ctx.Being.User.ID + ".b19")
		if Session == "" {
			ctx.Send(message.Text("未绑定账号,请先输入`/bind you_Session`进行绑定."))
			return
		}
		j := phigros.UserRecord{}
		data, _ := phigros.GetDataFormTap(phigros.UserMeUrl, Session) //获取id
		var um phigros.UserMe
		_ = json.Unmarshal(data, &um)
		j.PlayerInfo = &phigros.PlayerInfo{
			Name:      um.Nickname,
			CreatedAt: um.CreatedAt,
			UpdatedAt: um.UpdatedAt,
			Avatar:    um.Avatar,
		}
		data, _ = phigros.GetDataFormTap(phigros.SaveUrl, Session) //获取存档链接
		var gs phigros.GameSave
		_ = json.Unmarshal(data, &gs)

		ScoreAcc, err := phigros.ParseStatsByUrl(gs.Results[0].GameFile.URL)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		j.ScoreAcc = phigros.BN(ScoreAcc, 21)
		j.Summary = phigros.ProcessSummary(gs.Results[0].Summary)
		data, err = web.GetData(j.PlayerInfo.Avatar, "")
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		f, err := os.Create(avatar + Session + ".png")
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		_, _ = f.Write(data)
		f.Close()
		err = DrawB19(0.5, j, strconv.FormatFloat(float64(j.Summary.Rks), 'f', 6, 64), challengemoderank[(j.Summary.ChallengeModeRank-(j.Summary.ChallengeModeRank%100))/100], strconv.Itoa(int(j.Summary.ChallengeModeRank%100)), Session)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		ctx.Send(message.Image("file://" + file.BOTPATH + "/" + output + Session + ".png"))
	})
	en.AddRex(`^/bind\s*(.*)$`).Handle(func(ctx *rosm.Ctx) {
		if ctx.Being.Rex[1] == "" {
			ctx.Send(message.Text("Session不能为空"))
			return
		}
		_, err := phigros.GetDataFormTap(phigros.UserMeUrl, ctx.Being.Rex[1]) //获取id
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		SaveSession(en.DataFolder+"sessions/"+ctx.Being.User.ID+".b19", ctx.Being.Rex[1])
		ctx.Send(message.Text("绑定成功,发送`/b19`查询战绩"))
	})
	en.AddRex("^/phi(扫码)?登录$").Handle(func(ctx *rosm.Ctx) {
		r, err := qr.LoginQrCode(true, "public_profile")
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		if ctx.Being.Rex[1] != "" {
			var png []byte
			png, err = qrcode.Encode(r.Data.QrcodeURL, qrcode.Medium, 256)
			if err != nil {
				ctx.Send(message.Text("ERROR: ", err))
				return
			}
			//_ = os.WriteFile("qr.png", png, 0666)
			ctx.Send(message.ImageByte(png))
		} else {
			ctx.Send(message.Text(r.Data.QrcodeURL))
		}
		result, err := qr.CheckQRCode(true, r)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		deadline := time.Now().Add(5 * time.Minute)
		for !result.Success {
			time.Sleep(2 * time.Second)
			result, err = qr.CheckQRCode(true, r)
			if err != nil {
				ctx.Send(message.Text("ERROR:", err))
				return
			}
			if time.Now().After(deadline) {
				ctx.Send(message.Text("ERROR:登录超时,超过5分钟未成功登录。"))
				return
			}
		}
		p, err := qr.GetProfile(true, result)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		k, err := qr.LoginAndGetToken(result, p, false)
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		//	fmt.Println(k.SessionToken)
		if k.SessionToken == "" {
			ctx.Send(message.Text("Session不能为空"))
			return
		}
		_, err = phigros.GetDataFormTap(phigros.UserMeUrl, k.SessionToken) //获取id
		if err != nil {
			ctx.Send(message.Text("ERROR: ", err))
			return
		}
		SaveSession(en.DataFolder+"sessions/"+ctx.Being.User.ID+".b19", k.SessionToken)
		ctx.Send(message.Text("绑定成功,发送`/b19`查询战绩"))
	})
	en.AddWord("/phi帮助").Handle(func(ctx *rosm.Ctx) {
		ctx.Send(message.Text("Session获取:\n用mt文件管理器打开`.userdata`文件\n打开后找到`sessionToken:xxx`\nxxx即为所需\n.userdata文件的相对路径: “./Android/data/com.PigeonGames.Phigros/files/.userdata"))
	})
}

func FindSessionFID(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return tool.BytesToString(data)
}

func SaveSession(path, Session string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(Session)
	return err
}
