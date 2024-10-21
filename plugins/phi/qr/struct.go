package qr

import "time"

type LoginQrCodeRes struct {
	ClientId string `json:"-"`
	Data     struct {
		Error
		LoginQrCodeData
	} `json:"data"`
	Message
}
type LoginQrCodeData struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	QrcodeURL       string `json:"qrcode_url"`
}
type CheckQRCodeRes struct {
	Data struct {
		Error
		CheckQRCodeData
	} `json:"data"`
	Message
}

type CheckQRCodeData struct {
	Kid          string `json:"kid"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	MacKey       string `json:"mac_key"`
	MacAlgorithm string `json:"mac_algorithm"`
	Scope        string `json:"scope"`
}

type Profile struct {
	Data struct {
		Error
		ProfileData
	} `json:"data"`
	Message
}

type ProfileData struct {
	Avatar  string `json:"avatar"`
	Gender  string `json:"gender"`
	Name    string `json:"name"`
	Openid  string `json:"openid"`
	Unionid string `json:"unionid"`
}

type Error struct {
	Code             int    `json:"code"`
	Msg              string `json:"msg"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type Message struct {
	Now     int  `json:"now"`
	Success bool `json:"success"`
}

type AuthData struct {
	AuthData struct {
		TapTap struct {
			CheckQRCodeData
			ProfileData
		} `json:"taptap"`
	} `json:"authData"`
}

type AuthRes struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	ACL   struct {
		NAMING_FAILED struct {
			Write bool `json:"write"`
			Read  bool `json:"read"`
		} `json:"*"`
	} `json:"ACL"`
	AuthData struct {
		Taptap struct {
			AccessToken  string `json:"access_token"`
			Avatar       string `json:"avatar"`
			Gender       string `json:"gender"`
			Kid          string `json:"kid"`
			MacAlgorithm string `json:"mac_algorithm"`
			MacKey       string `json:"mac_key"`
			Name         string `json:"name"`
			Openid       string `json:"openid"`
			Scope        string `json:"scope"`
			TokenType    string `json:"token_type"`
			Unionid      string `json:"unionid"`
		} `json:"taptap"`
	} `json:"authData"`
	Avatar              string    `json:"avatar"`
	CreatedAt           time.Time `json:"createdAt"`
	EmailVerified       bool      `json:"emailVerified"`
	MobilePhoneVerified bool      `json:"mobilePhoneVerified"`
	Nickname            string    `json:"nickname"`
	ObjectID            string    `json:"objectId"`
	SessionToken        string    `json:"sessionToken"`
	ShortID             string    `json:"shortId"`
	UpdatedAt           time.Time `json:"updatedAt"`
	Username            string    `json:"username"`
}
