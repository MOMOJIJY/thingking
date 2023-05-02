package config

import (
	"log"
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

func InitBase() {
	AppID = os.Getenv("APPID")
	AppSecret = os.Getenv("APPSECRET")
	log.Printf("getting env. appID=%s, appSecret=%s\n", AppID, AppSecret)
}
