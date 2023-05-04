package models

type AccessToken struct {
	Token      string `json:"access_token"`
	ExpireTime int64  `json:"expires_in"`
}

type CommonResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
