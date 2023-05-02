package models

import "fmt"

type WeError struct {
	Msg  string
	Code int64
}

func (w *WeError) Error() string {
	return fmt.Sprintf("errcode=%d, errmsg=%s", w.Code, w.Msg)
}

var (
	ErrorGetAccessToken = WeError{Msg: "get access token failed", Code: 100001}
)
