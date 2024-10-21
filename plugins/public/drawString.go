package walle

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"math/rand"
	"os"
	"strings"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
)

func StringToShake(txt string, font string) (pic []byte, err error) {
	nilx, nily := 1.0, 8.0
	s := []*image.NRGBA{}
	strlist := strings.Split(txt, "\n")
	data, err := os.ReadFile(font)
	if err != nil {
		return nil, err
	}
	//获得画布预计
	testcov := gg.NewContext(1, 1)
	if err = testcov.ParseFontFace(data, 30); err != nil {
		return nil, err
	}
	//取最长段
	txt = ""
	for _, v := range strlist {
		if len([]rune(v)) > len([]rune(txt)) {
			txt = v
		}
	}
	w, h := testcov.MeasureString(txt)
	for range 10{
		cov := gg.NewContext(int(w+float64(len([]rune(txt)))*nilx)+40, int((h+nily)*float64(len(strlist)))+30)
		cov.SetRGB(1, 1, 1)
		cov.Clear()
		if err = cov.ParseFontFace(data, 30); err != nil {
			return nil, err
		}
		cov.SetColor(color.NRGBA{R: 0, G: 0, B: 0, A: 127})
		for k, v := range strlist {
			for kk, vv := range []rune(v) {
				x, y := cov.MeasureString(string([]rune(v)[:kk]))
				cov.DrawString(string(vv), x+float64(rand.Intn(5))+10+nilx, y+float64(rand.Intn(5))+15+float64(k)*(y+nily))
			}
		}
		s = append(s, imgfactory.Size(cov.Image(), 0, 0).Image())
	}
	var buf bytes.Buffer
	gif.EncodeAll(&buf, imgfactory.MergeGif(5, s))
	return buf.Bytes(), nil
}
func StringToPic(txt string, font string) (pic []byte, err error) {
	nilx, nily := 1.0, 8.0
	strlist := strings.Split(txt, "\n")
	data, err := os.ReadFile(font)
	if err != nil {
		return nil, err
	}
	//获得画布预计
	testcov := gg.NewContext(1, 1)
	if err = testcov.ParseFontFace(data, 30); err != nil {
		return nil, err
	}
	//取最长段
	txt = ""
	for _, v := range strlist {
		if len([]rune(v)) > len([]rune(txt)) {
			txt = v
		}
	}
	w, h := testcov.MeasureString(txt)
	cov := gg.NewContext(int(w+float64(len([]rune(txt)))*nilx)+40, int((h+nily)*float64(len(strlist)))+30)
	cov.SetRGB(1, 1, 1)
	cov.Clear()
	if err = cov.ParseFontFace(data, 30); err != nil {
		return nil, err
	}
	cov.SetColor(color.NRGBA{R: 0, G: 0, B: 0, A: 127})
	for k, v := range strlist {
		for kk, vv := range []rune(v) {
			x, y := cov.MeasureString(string([]rune(v)[:kk]))
			cov.DrawString(string(vv), x+10+nilx, y+15+float64(k)*(y+nily))
		}
	}
	return imgfactory.ToBytes(cov.Image())
}
