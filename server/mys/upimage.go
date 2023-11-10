package mys

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/rosm"
	"github.com/lianhong2758/RosmBot-MUL/tool"
	"github.com/lianhong2758/RosmBot-MUL/tool/web"
	log "github.com/sirupsen/logrus"
)

const mysUpimageUrl = "https://bbs-api.miyoushe.com/vila/api/bot/platform/getUploadImageParams"

func UploadFile(ctx *rosm.CTX, path string) (imageUrl string, err error) {
	log.Info("[mys]上传图片到米游社阿里云 OSS")
	// 在这里读取本地图片文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// md5
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, file); err != nil {
		return "", err
	}
	// 在这里获取机器人开放平台下发的 oss 参数
	param, err := getParam(ctx, tool.BytesToString(md5hash.Sum(nil)), strings.ToLower(filepath.Ext(path)))
	if err != nil {
		log.Error("[mys]获取米游社阿里云 OSS 上传参数失败", err)
		return "", err
	}
	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)
	// 添加表单字段
	_ = multiPartWriter.WriteField("x:extra", param.Data.Params.CallbackVar.XExtra)
	_ = multiPartWriter.WriteField("OSSAccessKeyId", param.Data.Params.Accessid)
	_ = multiPartWriter.WriteField("signature", param.Data.Params.Signature)
	_ = multiPartWriter.WriteField("success_action_status", param.Data.Params.SuccessActionStatus)
	_ = multiPartWriter.WriteField("name", param.Data.Params.Name)
	_ = multiPartWriter.WriteField("callback", param.Data.Params.Callback)
	_ = multiPartWriter.WriteField("x-oss-content-type", param.Data.Params.XOssContentType)
	_ = multiPartWriter.WriteField("key", param.Data.Params.Key)
	_ = multiPartWriter.WriteField("policy", param.Data.Params.Policy)

	fileWriter, err := multiPartWriter.CreateFormFile("file", "test.jpg")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", err
	}

	err = multiPartWriter.Close()
	if err != nil {
		return "", err
	}
	data, err := web.Web(web.NewDefaultClient(), param.Data.Params.Host, http.MethodPost, func(r *http.Request) {
		r.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	}, &requestBody)
	if err != nil {
		return "", err
	}
	m := new(OssDownloadParam)
	err = json.Unmarshal(data, m)
	return m.URL, err
}

// mys消息的ctx,md5,扩展名
func getParam(ctx *rosm.CTX, md5, ext string) (param *OssUpParam, err error) {
	data, _ := json.Marshal(H{"md5": md5, "ext": ext})
	data, err = web.Web(web.NewDefaultClient(), mysUpimageUrl, http.MethodGet, makeHeard(ctx), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	m := new(OssUpParam)
	err = json.Unmarshal(data, m)
	return m, err
}

// 阿里云oss所需要的数据
type OssUpParam struct {
	ApiCode
	Data struct {
		Type     string `json:"type"`
		FileName string `json:"file_name"`
		Params   struct {
			Accessid    string `json:"accessid"`
			Callback    string `json:"callback"`
			CallbackVar struct {
				XExtra string `json:"x:extra"`
			} `json:"callback_var"`
			Dir                 string `json:"dir"`
			Expire              int    `json:"expire"`
			Host                string `json:"host"`
			Name                string `json:"name"`
			Policy              string `json:"policy"`
			Signature           string `json:"signature"`
			XOssContentType     string `json:"x_oss_content_type"`
			ObjectAcl           string `json:"object_acl"`
			ContentDisposition  string `json:"content_disposition"`
			Key                 string `json:"key"`
			SuccessActionStatus string `json:"success_action_status"`
		} `json:"params"`
		MaxFileSize int `json:"max_file_size"`
	} `json:"data"`
}

type OssDownloadParam struct {
	Object    int    `json:"object"`
	SecretURL string `json:"secret_url"`
	URL       string `json:"url"`
}
