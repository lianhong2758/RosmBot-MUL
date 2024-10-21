package qr

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	UserPath      = "https://rak3ffdi.cloud.tds1.tapapis.cn/1.1/"
	WebHost       = "https://accounts.tapapis.com"
	ChinaWebHost  = "https://accounts.tapapis.cn"
	ApiHost       = "https://open.tapapis.com"
	ChinaApiHost  = "https://open.tapapis.cn"
	CodeUrl       = WebHost + "/oauth2/v1/device/code"
	ChinaCodeUrl  = ChinaWebHost + "/oauth2/v1/device/code"
	TokenUrl      = WebHost + "/oauth2/v1/token"
	ChinaTokenUrl = ChinaWebHost + "/oauth2/v1/token"

	AppKey   = "Qr9AEqtuoSVS3zeD6iVbM4ZC0AtkJcQ89tywVyi0"
	ClientId = "rAK3FfdieFob2Nn8Am"
)

func LoginQrCode(useChinaEndpoint bool, permissions ...string) (r *LoginQrCodeRes, err error) {
	clientId := strings.ReplaceAll(uuid.New().String(), "-", "")
	data, err := json.Marshal(map[string]any{"device_id": clientId})
	if err != nil {
		return nil, err
	}
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	writer.WriteField("client_id", ClientId)
	writer.WriteField("response_type", "device_code")
	writer.WriteField("scope", strings.Join(permissions, ","))
	writer.WriteField("version", "2.1")
	writer.WriteField("platform", "unity")
	writer.WriteField("info", string(data))
	writer.Close()
	var endpoint string = CodeUrl
	if useChinaEndpoint {
		endpoint = ChinaCodeUrl
	}
	res, err := http.Post(endpoint, writer.FormDataContentType(), &requestBody)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, _ = io.ReadAll(res.Body)
	r = new(LoginQrCodeRes)
	r.ClientId = clientId
	return r, json.Unmarshal(data, r)
}

func CheckQRCode(useChinaEndpoint bool, r *LoginQrCodeRes) (c *CheckQRCodeRes, err error) {
	data, err := json.Marshal(map[string]any{"device_id": r.ClientId})
	if err != nil {
		return nil, err
	}
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	writer.WriteField("grant_type", "device_token")
	writer.WriteField("client_id", ClientId)
	writer.WriteField("secret_type", "hmac-sha-1")
	writer.WriteField("code", r.Data.DeviceCode)
	writer.WriteField("version", "1.0")
	writer.WriteField("platform", "unity")
	writer.WriteField("info", string(data))
	writer.Close()
	var endpoint string = TokenUrl
	if useChinaEndpoint {
		endpoint = ChinaTokenUrl
	}
	res, err := http.Post(endpoint, writer.FormDataContentType(), &requestBody)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, _ = io.ReadAll(res.Body)
	c = new(CheckQRCodeRes)
	return c, json.Unmarshal(data, c)
}

func GetProfile(useChinaEndpoint bool, c *CheckQRCodeRes) (p *Profile, err error) {
	if c.Data.Scope != "public_profile" {
		return nil, errors.New("public profile permission is required")
	}
	var urlt string
	if useChinaEndpoint {
		urlt = ChinaApiHost + "/account/profile/v1?client_id=rAK3FfdieFob2Nn8Am"
	} else {
		urlt = ApiHost + "/account/profile/v1?client_id=rAK3FfdieFob2Nn8Am"
	}
	request, err := http.NewRequest(http.MethodGet, urlt, nil)
	if err != nil {
		return nil, err
	}
	t := fmt.Sprintf("%010d", time.Now().Unix())
	uri, _ := url.Parse(urlt)
	randomStr := RandomBase64String(16)
	port := uri.Port()
	if port == "" {
		switch strings.Split(urlt, ":")[0] {
		case "https":
			port = "443"
		default:
			port = "80"
		}
	}
	request.Header.Add("Authorization", fmt.Sprintf(`MAC id="%v", ts="%v", nonce="%v", mac="%v"`,
		c.Data.Kid, t, randomStr, SignData(MergeData(t, randomStr, http.MethodGet, uri.Path+"?"+uri.RawQuery, uri.Host, port, ""), []byte(c.Data.MacKey))))
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	p = new(Profile)
	return p, json.Unmarshal(data, p)
}

func RandomBase64String(length int) string {
	randomBytes := make([]byte, length)
	_, _ = rand.Read(randomBytes)
	randomBase64String := base64.StdEncoding.EncodeToString(randomBytes)

	return randomBase64String
}

func SignData(signatureBaseString string, key []byte) string {
	h := hmac.New(sha1.New, key)
	h.Write([]byte(signatureBaseString))
	signature := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

func MergeData(time, randomCode, httpType, uri, domain, port, other string) string {
	prefix := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n", time, randomCode, httpType, uri, domain, port, other)
	return prefix
}
