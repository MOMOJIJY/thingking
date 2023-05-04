package config

import (
	"os"
	"time"

	"github.com/kpango/glg"
)

const (
	Token = "wechat"

	// redis key
	RedisMemberKey   = "member"
	RedisTokenKey    = "access_token"
	RedisTokenExpire = 60 * time.Minute
)

var (
	AppID          string
	AppSecret      string
	ChangeBgApiKey string
)

func InitBase() {
	AppID = os.Getenv("APPID")
	AppSecret = os.Getenv("APPSECRET")
	ChangeBgApiKey = os.Getenv("BGKEY")
	glg.Infof("getting env. appID=%s, appSecret=%s, bg_key=%s\n", AppID, AppSecret, ChangeBgApiKey)
}
