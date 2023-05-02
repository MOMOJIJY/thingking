package config

import "time"

const (
	Token     = "wechat"
	AppID     = "wx5046b7123c5bb976"
	AppSecret = "6702f6773874195b079c1a161861946a"

	// redis key
	RedisMemberKey   = "member"
	RedisTokenKey    = "access_token"
	RedisTokenExpire = 90 * time.Minute
)
