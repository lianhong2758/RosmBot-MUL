package qr

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func LoginAndGetToken(r *CheckQRCodeRes, p *Profile, failOnNotExist bool) (k *AuthRes, err error) {
	var auth AuthData
	auth.AuthData.TapTap.CheckQRCodeData = r.Data.CheckQRCodeData
	auth.AuthData.TapTap.ProfileData = p.Data.ProfileData
	authData, _ := json.Marshal(auth)
	path := "users"
	if failOnNotExist {
		path = "users?failOnNotExist=true"
	}
	request, err := http.NewRequest(http.MethodPost, UserPath+path, bytes.NewReader(authData))
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-LC-Id", ClientId)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-LC-Sign", generateSign(AppKey))
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	authData, _ = io.ReadAll(res.Body)
	k = new(AuthRes)
	return k, json.Unmarshal(authData, k)
}

func generateSign(appKey string) string {
	timestamp := time.Now().Unix()
	data := strconv.FormatInt(timestamp, 10) + appKey

	hash := md5.Sum([]byte(data))
	hashStr := hex.EncodeToString(hash[:])

	sign := fmt.Sprintf("%s,%d", hashStr, timestamp)
	return sign
}
