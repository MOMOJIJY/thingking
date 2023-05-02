package config

import (
	"os"
	"time"
)

const (
	Token = "wechat"

	// redis key
	RedisMemberKey   = "member"
	RedisTokenKey    = "access_token"
	RedisTokenExpire = 90 * time.Minute
)

var (
	AppID     string
	AppSecret string
)

func init() {
	AppID = os.Getenv("APPID")
	AppSecret = os.Getenv("APPSECRET")
}
