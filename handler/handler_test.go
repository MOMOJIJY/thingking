package handler

import (
	"testing"
	"wechat-dev/thinking/config"
	"wechat-dev/thinking/models"

	"github.com/gin-gonic/gin"
)

func TestDoTemplate(t *testing.T) {
	config.InitRedis()

	h := &Handler{}
	msg := models.MsgFromUser{
		FromUserName: "123",
	}
	h.doTemplate(&gin.Context{}, 2, msg)
}
