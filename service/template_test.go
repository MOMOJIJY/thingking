package service

import (
	"testing"
	"wechat-dev/thinking/config"
)

func TestDoTemplate(t *testing.T) {
	config.InitRedis()
	DoTemplate()
}
