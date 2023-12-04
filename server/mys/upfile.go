package mys

import (
	"encoding/base64"
	"image"
	"os"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	log "github.com/sirupsen/logrus"
)

// 上传byte数据
func UpImgByte(ctx *rosm.CTX, img []byte) (url string, con *image.Config) {
	url, err := UploadFile(ctx, img)
	if err != nil {
		log.Warnln("[mys](upimage)上传图片失败,ERROR: ", err)
		return "", nil
	}
	c, _ := web.BytesToConfig(img)
	return url, &c
}

// 上传file
func UpImgfile(ctx *rosm.CTX, filePath string) (url string, con *image.Config) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Warnln("[mys](upimage)", err)
		return
	}
	return UpImgByte(ctx, file)
}

// 转存图片
func UpImgUrl(ctx *rosm.CTX, imgurl string) (url string) {
	url, err := TransFerImage(ctx, imgurl)
	if err != nil {
		log.Errorln("[mys](upimage)", err)
		return
	}
	return url
}

// 解析base64等data
func ImageAnalysis(ctx *rosm.CTX, data string) (url string, con *image.Config) {
	switch parts := strings.SplitN(data, "://", 2); parts[0] {
	case "base64":
		bytes, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			log.Warnln("[mys](upimage) ERROR:", err)
			return "", nil
		}
		return UpImgByte(ctx, bytes)
	case "file":
		return UpImgfile(ctx, parts[1])
	case "url":
		c, err := web.URLToConfig(url)
		if err != nil {
			log.Warnln("[mys](upimage) ERROR:", err)
			return "", nil
		}
		return UpImgUrl(ctx, parts[1]), &c
	case "consturl":
		return parts[1], nil
	default:
		return "", nil
	}
}
