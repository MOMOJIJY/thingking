package models

type MsgFromUser struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string `xml:",omitempty"`
	MsgId        int64  `xml:",omitempty"`
	Event        string `xml:",omitempty"`
	EventKey     string `xml:",omitempty"`
}

type AccessToken struct {
	Token      string `json:"access_token"`
	ExpireTime int64  `json:"expires_in"`
}

type CommonResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
